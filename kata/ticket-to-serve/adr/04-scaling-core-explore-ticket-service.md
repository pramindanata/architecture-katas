# ADR04 - Scaling Core Explore Ticket Service

Date: 2026-01-03

## Status

Accepted

## Context

All show & ticket read operation will be handled by the core explore ticket service. This service needs to be highly available and elastic to handle massive traffic spikes during peak traffic.

## Decision

### Database

There is no need for specialized search operation such as full text search or complex filtering. So I can still use **PostgreSQL** as the main database for this service. This service will have its own dedicated DB instance to isolate the load. Horizontal scale will be applied via Citus to distribute the read load.

The scaling level of this DB isn't as complex as the core order ticket service because it only handles read operation.

If specialized data structure is needed in the future, PostgreSQL extension can be used unless the achieved performance is not enough.

### In-Memory Store

To support fast read operation especially during traffic spike, I can use its own **Redis** instance as the in-memory store to cache the show & ticket data. Redis Cluster with multiple nodes and replicas will be used to ensure the availability.

Using the same Redis instance used by the core order ticket service is not a good idea because it will increase the blast radius if the Redis instance goes down.

### Show & Ticket Data Sync

The ticket management service is the source of truth for show & ticket data. To sync the data to the core explore ticket service, I can use the same **Kafka** instance proposed for the core order ticket service architecture to send the data update event.

### Result

![diagram](../asset/core-explore-ticket-service-architecture.svg)

## Consequences

### Positive

- High availability and elasticity during traffic spike.
- Fast read operation for show & ticket data.
- Easy to maintain & operate because it uses familiar technology stack.

### Negative

- More complex infrastructure to maintain.
