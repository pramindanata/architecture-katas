# ADR07 - Scaling Admin and Reseller Services

Date: 2026-01-04

## Status

Accepted

## Context

Only the admin & reseller services left from the planned initial architecture. Both of these services have the lowest traffic compared to the other, even the reseller service only need to handle a few hundred of resellers.

With this lower requirement, I need to decide how to scale these services properly without over engineering the architecture.

## Decision

### Database

Each service will have its own dedicated database. I choose **PostgreSQL** for this purpose because both services don't have massive read & write operation that need to be handled. A SQL database that guarantee ACID, relationship, and offer various extension for flexibility is enough for this use case.

Sharding & read replication may need to be applied to optimize some read operations such as reading ticket sales report & list of created orders in the reseller service.

### In-Memory Store

No in-memory store is needed for both services because the low traffic that they have.

### Syncing Ticket Order Data from Core Order Ticket Service

The reseller need to read details of ticket sales report from the core order ticket service. Because the massive traffic that the core order ticket service have, the core order ticket service need to communicate with the reseller service in an optimized way to avoid overloading the reseller service & its DB, for example, using batching.

### Result

![diagram](../asset/admin-and-reseller-services-architecture.svg)

## Consequences

### Positive

- Both services have simple architecture that is easy to maintain.

### Negative

- Both services may need to be optimized further if the traffic increases drastically in the future.
- Syncing ticket order data from the core order ticket service need to be designed properly to avoid overloading the reseller service.
