# ADR00 - Offload Delivery Routes

Date: 2025-10-15

## Status

Accepted

## Context

With thousand-to-million presents need to be delivered, system need to generate the delivery routes as optimal as it can with less time. Because a lot of data need to be processed, the main service may consume more resources and can disturb its other functionalities.

## Decision

We will **offload** the "generate routes" process into a dedicated microservice with its own resources. This service will communicate with main service with sync communication because the route generation must be triggered manually by Santa. This service will use the main service DB for fetching & storing data.

The details of how the routes will be generated can be viewed in this [TRD](../trd/generating-delivery-routes.md).

## Consequences

### Positive

- Reduce disturbance in the main service.

### Negative

- A new service need to be maintained.
- Increase resource cost.
- Add more complexity to the service communication.
