# ADR01 - Choosing Architecture Style

Date: 2025-01-02

## Status

Accepted

## Context

Selecting architecture style will have huge impact on how the system will structure, built and performs. From the discovery process, the key architecture characteristic that drive our design are.

- (Top) Consistency
- (Top) Availability
- (Top) Elasticity
- Scalability
- Interoperability
- Security

Not all features require the same level of consistency, availability, and elasticity. Only the ticket purchasing flow needs these top characteristics.

## Decision

We choose **Service Based Architecture**.

![diagram](../asset/architecture-styles-worksheet.jpg)

It offers flexibility in how the system will structure, built and performs. Service based allow us selectively choose which features that require high level of availability & elasticity while the rest of features can be built using simpler architecture.

The microservice architecture can also be considered. This architecture can maximize the performance, scalability, and elasticity but not all features requires all of those characteristics.

## Consequences

### Positive

- Allow selective optimization.
- Reduce development cost & time for less critical components.
- Provide a clear migration path toward microservices if future scalability needs increase.

### Negative

- Require careful boundary definition for each services to avoid tight coupling.
- More complex inter-service communication need to be developed.
