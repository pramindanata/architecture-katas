# ADR09 - Add API Gateway

Date: 2025-10-16

## Status

Accepted

## Context

The main web application needs to communicate to some backend services. To make service discovery easier, limit APIs that can be accessed, & handle request authentication, a dedicated API gateway is needed.  

## Decision

We will add **Traefik** as the API Gateway. It covers all requirements. Also, it has better performance efficiency, resource usage, & learning curve.

Other alternatives:

- Kong: offer richer features & plugins but has higher resource usage & learning curve.
- NGINX: simplest API gateway with higher performance & resource efficiency than Kong & Traefik. Unfortunately, it can't handle request authentication.

## Consequences

### Positive

- Traefik offer better performance efficiency, resource usage, & learning curve.

### Negative

- Traefik has fewer features & plugins than Kong. We can move to Kong in the future if needed.
