# ADR02 - Choose Inter-Service Communication Style

Date: 2025-10-14

## Status

accepted

## Context

Selected architecture style introduce inter-service communication to the system. The communications can be sync or async depending on the scenarios. Standardizing both communication style and technology are essential to ensure the system meets the its driven architecture characteristic especially: elasticity, availability, & performance.

## Decision

We choose to use both **sync & async communication** for this system because both of them can accomodate certain scenarios that exist in the system.

- Sync communication is used for a communication where immediate response is required.
- Async communication is used for a communication that decoupled from the caller, allowing the system to remain resilient & responsive during high load.

### Sync Communication

For the sync communication, **HTTPS is enough**. Currently there is no sync communication that require low latency and high throughput. Also, HTTPS is simpler to implemented.

### Aysnc Communication

For the async communication, the system must handle coordination between services to process a large number of present requests within a short period of time. The messaging technology must support the following criteria.

- Capable handling hundreds thousand of requests at a short period
- Support at-least once delivery semantic. It ensure no lost messages with cons of possibility of duplicated messages. Duplications are accepted because our system capable of consolidating duplicate data. Exactly once delivery semantic can be considered but it resource intensive.
- Support distributing messages equally to available subscribers that deployed in different pod.
- If possible, it has lower resource cost.

We choose **Kafka** to handle async communication. It is battle tested and support majority of our criteria. Unfortunately, it has higher infrastructure and operational cost compared to some alternatives.

Considered alternative:

- Apache Pulsar: high throughput (less than Kafka), cheaper in resource cost, and not mature as Kafka.
- RabbitMQ: moderate throughput (10K-to-hundred-thousand per second) and cheaper in resource cost. Not good for elasticity & scaling for future.
- Cloud-managed (such as Google Pub/Sub): high throughput (more than Pulsar), less setup in infra, but the cost calculation is more complex & we have less control. The total cost can be higher than Kafka.

## Consequences

### Positive

- Flexibility to choose communication style depending on the scenario.
- Kafka ensures high throughput and reliability for asynchronous workflows.  
- Synchronous operations remain simple, low-cost, and easy to integrate via HTTPS.

### Negative

- Increased system complexity due to managing both sync and async patterns.  
- Kafka requires additional operational expertise and infrastructure resources.  
- Development teams must handle potential message duplication and ensure idempotency in consumers.

## Reference

- [Confluent: Kafka vs Pulsar](https://www.confluent.io/kafka-vs-pulsar/)
- [AutoMQ: Kafka vs Google Pub/Sub](https://www.automq.com/blog/apache-kafka-vs-google-pub-sub-differences-and-comparison)
- [Confluent Blog: Kafka — The Fastest Messaging System](https://www.confluent.io/blog/kafka-fastest-messaging-system/)
