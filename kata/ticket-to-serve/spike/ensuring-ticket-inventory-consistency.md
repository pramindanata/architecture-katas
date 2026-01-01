# Ensuring Ticket Inventory Consistency

## Problem

When a ticket is purchased, system must ensure the ticket final stock number is correct and consistent.

## Solution

Well, the solution is tech specific. The database that store the ticket stock number must support the row level lock and prevent dirty read when reading data inside a transaction.

When a transaction occurs (where it will decrement the stock number), we can apply the following queries (for example in PostgreSQL).

```sql
begin transaction;

-- Apply pesimistic lock 
select * from tickets where id in ($idA, $idB) for update;

-- Do update with minimum "finalStock >= 0" constraint
update tickets
  set stock = stock - $qty
  where id = $idA and stock - $qty >= 0

update tickets
  set stock = stock - $qty
  where id = $idB and stock - $qty >= 0

commit;
```

In PostgreSQL, update & delete statement will acquire a row level lock for the modified records.  It ensures others write queries that will write the same records need to wait until the lock is released. Technically the `select ... for update` is optional, it allows us to lock the row for longer period in case we need to do certain operations before writing the records.

Also, all PostgreSQL isolation levels prevent dirty read, so it ensures any read inside the transaction can only read the latest committed values.

Make sure to update the stock value with `stock = stock - $qty` instead of `stock = $newStock`. The later approach is more vulnerable to a racing condition while the first approach ensure the update query to read the last committed `stock` value. The `stock = $newStock` may work well if the write query to that record do `select ... for update` **before** calculating the `newStock` value.

This approach will ensure consistent decrement because each write will be treated like serial operation for each ticket record. Throughput will suffer, but it is the trade off that we need to make to solve this consistency problem.

If the throughput becoming a concern, optimistic lock can be the alternative. The consequences of this approach is if there are `n` purchase happening at same time for ticket x, only 1 can succeed while the other will receive an error. Applying retry can be a solution but it will add more pressure to the DB.

Adding a constraint like `stock int not null check (stock >= 0)` can be good practice to avoid other application codes accidentally decrement the stock value to less than 0.
