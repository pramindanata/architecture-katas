# Purchase Ticket Flow

## Problem

We need to design a ticket purchase flow that can fulfill the following constraints:

- Handle massive traffic spike when the ticket purchasing period opens for highly anticipated shows.
- Ensure strong consistency and correctness of ticket inventory. No double purchase & oversell.
- Ensure fairness. Users who order first should have higher chance to get the ticket.

## Solution

Here is the general idea of the flow. Some aspect may be missing such as check available stock from memory store before publishing the order into queue.

![diagram](../asset/purchase-ticket-flow.svg)

### Handling the Traffic

To handle massive traffic, we can use queue to distribute the load. Each order request will be placed in a queue and users will wait until their requests are processed. When the request out from the queue, system starts to reserve the ticket by decrement the ticket stock, create order records, and notify client side, so user can complete the payment process.

Each order have its own expiration time. If the expiration time reached, system will revert the stock and cancel the order.

Users also need to know how long to wait and number of users that are ahead of them in the queue. We can provide an estimated waiting time based on the average processing time of the queue handler and number of messages in the queue.

Using stream-based or log-based queue (such as Kafka or Redis Stream) can slow down the traffic well. Each consumer can only process 1 message per partition per topic at a time. Applying custom topic & partition can also maximize the throughput, for example:

```js
 // dedicated for highly anticipated shows that have million of tickets.
- `orderTicket:topA`
- `orderTicket:topB`
 // dedicated for other shows
- `orderTicket:normal` 
```

Admin or system can configure which topic to use for a show when publishing the "order ticket" message to the queue. If we have 10 consumers and set each topic with 10 partitions, then all messages can be distributed to all consumers equally and each consumer can handle 3 topics at the same time. This will prevent 1 highly anticipated show blocks all consumers.

### Ensuring the Stock consistency

There are some approaches that can applied here. Each have its own trade off between complexity, performance, and consistency. Personally I prefer the "Simple DB Implementation" if performance is acceptable. If not, then "DB - Split the Record" can be the alternative. More thorough test need to be done to measure the performance of each approach.

The terms of `Memory Store` in this section refers to technology like Redis & Memcached that allow fast read & write via memory, while `DB` refers to databases that guarantee ACID such as PostgreSQL & MySQL.

#### Simple DB Implementation

Inside the queue handler, we apply row level lock and decrement the stock. This is the simplest solution, but throughput will suffer because of the lock contention unless there are some DBs that have better performance even when lock contention occurs.

Row level lock usually a pessimistic lock where each write on a same record must wait until the lock is released. What about optimistic lock? This kind of lock provide faster throughput, but we need to return error or retry the process if the lock was acquired by another process. Retrying the write operation can make more active queue handlers storm the DB while returning an error will make user to start the process over again.

#### Using Memory Store

Using a memory store that allow atomic decrement/increment provide better throughput, while combined with DB to ensure data consistency. We can apply implementation like the following.

```
batch = []

// This is the queue job handler
handleJob(job):
  redis.decr("ticket:${job.data.ticketId}", job.data.ticketCount)
  batch.push(job.data)

handleFlush():
  map = {}
  createOrderOps = []

  for item of batch:
    map[item.ticketId] += item.ticketCount
    createOrderOps(...)

  transaction.run((tx) => {
    for id, count of map:
      db.exec("update set stock = stock - ? where id = ? and stock - ? >= 0", count, id, count)

    for op of ops:
      db.exec("insert into order ...")
  }) 

setInterval(handleFlush, 1) // process batch every 1 second.
```

What happens if the `handleJob()` crashed after decrementing the stock in the memory store? We may need to retry it and apply idempotency to ensure no double decrement. Now, what happens if the `handleFlush` crashed? This problem can make the implementation becoming more complex.

Moving the source of truth to memory store also can be dangerous because if the memory store is crashed, then all data will be lost. Enabling persistence can be a solution but this only minimize the amount of data loss. There is a stricter persistence configuration, but each write operation need to write data into the disk so it will impact to the performance.

Now we end up with complex implementation without significant improvement (need to be tested further).

#### DB - Split the Record

Instead of storing ticket data in a DB with the following structure.

```js
// table tickets
id, name, stock
```

We can split each `tickets` record into multiple "units".

```js
// table ticket_units
id, ticket_id, shard_id name, stock
```

Inside the queue handler, system will select which ticket units need to be decremented. With this approach, we can distribute the lock contention of a ticket into multiple unit records. If we don't use DB sharding/partitioning, "selecting" the ticket unit need advanced approach, but it is doable. Because of the possibility of high traffic case, sharding the DB can be applied to distribute the DB load.


```
handleJob(job):
  shardId = hash(job.userId) 

  transaction.run((tx) => {
    for ticketId of job.ticketIds:
      // The DB will handle the shard routing based on the given shard key. The `shard_id` will be the shard key.
      db.exec(
        "update set stock = stock - ? where id = ? and shard_id = ? and stock - ? >= 0", 
        count, 
        ticketId, 
        shardId,
        count, 
      )
  }) 
```

What happens if the ordered ticket count is larger than the selected unit stock? Trying the process to make it finds matched unit is unreliable especially if all units have low stock. To solve this problem, we can select the unit early before publishing the message. We create a data in memory store where it contains the available stock of the ticket in each shard.

```js
// In memory
key: `ticket:shard:{id}`
value: {
  <shard_a>: 200
  <shard_b>: 200
  ...
  <shard_z>: 200
}
```

When the "request order" HTTP handler is executed, system will select the shard from those data and decrement the shard key. It allows faster read & write operation to ensure the "request order" throughput. If the memory store crashed, system need to acquire a global lock to warm up the data from DB. Incoming request will be rejected when the lock is acquired.

The read & decrement process must be atomic, so the memory store need to support atomic process, for example the Lua script in Redis.

If it is possible, the read, decrement, & publish operations must be atomic. If the publish operation use different technology, at least the delayed "revert stock" job (explained in the another section), will ensure the stock consistency.

Also, if the stock value in memory store is larger than in DB, at least the queue handler will not decrement the stock if it was already 0.

#### DB - Split the Record into the Smallest Value

We can split the ticket record further. If a ticket stock is 1,000,000, then we store 1,000,000 of ticket records instead of 1 record. This will avoid lock contention but data size grow drastically. If there are 100 tickets that each has a million stocks, then the DB will store 100,000,000 records.

This approach also have the same problem as the "DB - Split the Record" approach.

### Notifying the Client Side

To notify client side after an order record is created, system can use technology like SSE or Web Socket. The details will be explored further in a dedicated ADR document.

### Reverting Ticket Stock

What if user don't complete the payment? The ticket stock should be reverted so other users can order it. We can use technique like scheduler or delayed job to revert the stock. Scheduler is the simpler approach but not real time. Delayed job offer is real time, but it will use more resources to store and process the job.

Scheduler can be a solution to handle this problem because there is no restriction about how fast user could see the reverted ticket stock. We can make the scheduler to run frequently (such as each 30s) but we need to ensure each execution don't process same reserved ticket twice.

### What If User Accidentally Closed the Browser?

We may need to store the following states in a memory store to allow the user see the queue position when the user open the browser again.

- Total of received order.
- Total of processed order.
- Number of the user order in the queue.
- Estimated order process time.

If the order already created, user can see it immediately in web page, so he/she can complete the payment.
