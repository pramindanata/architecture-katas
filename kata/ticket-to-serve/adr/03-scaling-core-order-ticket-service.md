# ADR03 - Scaling Core Order Ticket Service

Date: 2026-01-02

## Status

Accepted

## Context

The ticket purchase flow is mainly handled by the core order ticket service. This service needs to be highly available and elastic to handle massive traffic spikes during high-demand events. Also, ensure data consistency during the process is crucial to prevent overselling.

How I handle the flow is already described in [this document](../spike/ticketing-platform-auth.md). This document will focus on the architecture of the service itself.

## Decision

### Queue

Distributing the massive traffic will be handled by a queue mechanism. I need a queue that support the following requirements:

- Can handle massive enqueue operation during traffic spike (hundreds of thousands per second at peak time).
- The consumer can consume 1 message per topic per partition at a time to lower the server stress.
- Guarantee message ordering.
- Support at-least-once delivery guarantee.

I choose **Apache Kafka** as the queue mechanism because it meets all of those requirements. Also, I can use Kafka as a message broker for inter-service communication. There are some cons with this choice:

- More complex infrastructure to maintain.
- Higher learning curve for the development team.
- Higher resource consumption.
- It only guarantees global ordering if a topic only has 1 partition. Increasing the partition will make ordering only per partition. At least the whole message ordering is not random.

Considered alternatives:

- NATS JetStream: simpler infrastructure & lower resource consumption, but it doesn't support ordering guarantee.
- RabbitMQ: doesn't support massive enqueue operation & ordering guarantee.
- Apache Pulsar: more complex infrastructure & higher learning curve.

### Database

For the database, I can still use **PostgreSQL** as the database to ensure the stock consistency, but this service will have its own dedicated DB instance to isolate the load. Horizontal scale also will be applied via Citus to distribute the lock contention when decrementing the ticket stock.

A single PostgreSQL node can handle around ~1,600 update operations on the same record (according my benchmark, [see this code](../rnd/load-test/load-test-pg-decrement-ticket.js)). The number is very low comparing the incoming traffic, but with help of Kafka to distribute the load and horizontal scaling of PostgreSQL, I think this architecture can handle the load well.

With 10 nodes of PostgreSQL with Citus, I can handle around ~16,000 update operations on the same record. To improve the throughput further, I can split the ticket further so 1 shard can contain `n` ticket units with its own stock.

Considered alternatives:

- CockroachDB: natively support sharding & distributed SQL, but the ecosystem is not as mature as PostgreSQL.
- MySQL: support sharding natively with NDB Cluster, but the feature is not mature enough and the community support is not as strong as PostgreSQL.

### In-Memory Store

Before publishing the "order ticket request", system need to select which shard that has enough stock for the requested ticket and decrement it. To support fast read & write operation, I can use **Redis** as the in-memory store to store the available stock per shard. Redis support various data support and its Lua scripting can be used to ensure atomic read & decrement operation.

To ensure the Redis availability, I can use Redis Cluster with multiple nodes and replicas.

There are some cons with this approach:

- Data in Redis can be lost if the instance crashed. To mitigate this problem, system need to acquire a global lock to warm up the data from DB if Redis crashed.

Considered alternatives:

- Memcached: simpler infrastructure, but it doesn't support atomic read & decrement operation.

### Payment Provider

I will use a third party payment provider such as Stripe or Braintree to handle the payment process. It will reduce the complexity of the system and allow focusing on the core domain. The provider must support the following requirements:

- Support various payment methods (credit card, e-wallet, bank transfer, etc.).
- Support high traffic to create & handle a payment (about a thousand of transactions per second at peak time because of the queue).
- Can detect & prevent fraudulent transactions.
- Have webhook support to notify the system such as payment success, failure, fraud detected, etc.

I choose **Payment X**.

> To be honest, I didn't find any payment provider that can handle this kind of traffic, either it has lower rate limit or the number is not published. Contacting the provider sales team to discuss this requirement is a must.

Considered alternatives:

- Stripe
- Ayden
- Braintree

### Syncing Stock Data

The reseller service's DB will the source of truth of the ticket stock data. To sync the stock data to the core order ticket service, I can use an asynchronous approach (event-driven) where the reseller service will publish an event to Kafka whenever there is a change in the ticket stock data. The core order ticket service will have a consumer that listens to those events and update its own DB accordingly. Also, the core order ticket service need to inform the reseller service about the stock change after processing the order.

The sync communication can also be an alternative (such as HTTP request), but system need to handle retry the request in case of failure. Also, this approach will add more load to the reseller service.

### Sending Notification

Any notification for the user will be handled by the notification service via asynchronous communication. The core order ticket service will publish an event to Kafka whenever there is a need to send a notification such as ticket purchased and refund processed. 

### Final Result

![diagram](../asset/core-order-ticket-architecture.svg)

## Consequences

### Positive

- Allow high availability and elasticity during traffic spike.
- Data consistency is ensured during the ticket purchase process.
- The use of dedicated services and databases allows for better isolation of load and easier scaling.
- The use of third-party payment providers reduces the complexity of the system and allows focusing on the core domain.

### Negative

- The architecture is becoming more complex and requires more infrastructure to maintain (e.g., Kafka, Redis Cluster, PostgreSQL with Citus).
- Higher learning curve for the development team due to the use of multiple technologies.
- Increased infrastructure cost.
