# ADR02 - Initial Architecture

Date: 2026-01-02

## Status

Accepted

## Context

With the discovered requirements and constraints from the discovery phase, I need to design the initial architecture for this system.

## Decision

### Client Application

There will be 3 client applications:

- Reseller dashboard (web): for reseller admin to manage their shows, tickets, orders, support tickets, and see reports.
- Reseller custom web (web): alternative for reseller instead of using the Ticketing Platform APIs.
- Admin dashboard (web): for platform admin manage the system.

There is also the reseller system that will consume the Ticketing Platform APIs.

### Backend Service

Based on the available functional requirements, I separate the backend service into several bounded contexts:

#### Core Order Ticket Service

This is the service that requires the highest availability & elasticity because it will handle the ticket purchase flow with massive traffic spike. Responsibilities:

- Platform can create order & payment of the selected tickets.
- Platform can return queue status for the user order.
- Platform can return created orders.
- Platform can create request refund.
- Platform can receive notification from payment provider such as payment success, failure, refund processed, fraud detected, etc.

The last 2 points can be separated into another service if needed, but I think this service can handle them because the traffic should be smaller.

#### Core Explore Ticket Service

This is the second service that also need high availability & elasticity because it will handle read operation from many users. Responsibilities:

- Platform can return reseller's shows details.
- Platform can return available show's tickets details.

#### Ticket Management Service

This service will handle reseller's shows & tickets management. Responsibilities:

- Reseller can register an account to start managing shows & using the Ticketing Platform APIs. Including billing information.
- Reseller can manage shows and their tickets.
- Reseller can see detailed ticket sales report.
- Reseller can send a support ticket.
- Reseller can create & customize a custom web instead of consuming the APIs. Extra fees will apply.

#### Admin Service

This service will handle admin operations. Responsibilities:

- Admin can process support tickets such as force refund, complain, etc.
- Admin can mark shows that potentially will have high traffic when the ticket purchasing period opens.

I think this service can be merged into Ticket Management Service if the operational load is small.

#### Notification Service

This service will handle all notification sending such as email, webhook, etc. Because it is responsible for sending emails, it needs higher memory & CPU for processing email templates. Responsibilities:

- Platform can send email notification to user such as ticket purchased, refund processed, fraud detected, etc.
- Platform can send notification to reseller system such as ticket sold out, fraud detected, etc. via webhook.

#### Custom Web Service

This service will handle some functionalities for reseller custom web. Responsibilities:

- User can register & login into the reseller custom web.
- User can explore available shows & tickets details.
- User can receive real time notification when the ticket they are looking at just got purchased.

Some functionalities such as explore shows will require have higher load, but because the majority of the logics will be handled by the core explore ticket service, this service should have smaller load.

The real time notification can be handled by WebSocket or Server-Sent Events (SSE). It will be explored further in the next ADR.

#### Custom Web Order Service

This service will handle ticket order from reseller custom web. Responsibilities:

- User can order tickets.
- User can see his/her queue status.
- User can see his/her created orders.

I need to separate this service from the custom web service because the high load during ticket purchase period. The majority of the logics will be handled by the core order ticket service.

### Database

For the initial architecture, I will use a single relational database for all services. I choose PostgreSQL for its reliability, data correctness, features, and strong community support. Next ADR will discuss more about using dedicated or even specialized database per service.

To achieve strong availability, each DB especially for the core order & explore ticket services will be deployed horizontally. Unfortunately I need to use extension like Citus to achieve this because PostgreSQL does not support sharding natively.

The cons of choosing PostgreSQL is slower read & write performance compared to other distributed DB such as CockroachDB or NewSQL DBs. But with proper optimization and caching layer, I think PostgreSQL can handle the load well. Also, other features don't require special data structure that need specialized DB to achieve better performance.

Alternatives considered:

- MySQL: support sharding natively with NDB Cluster, but the feature is not mature enough and the community support is not as strong as PostgreSQL.
- CockroachDB: natively support sharding & distributed SQL, but the ecosystem is not as mature as PostgreSQL.

### API Gateway

Because there are many services that need to be exposed to the client applications, I will use an API Gateway to simplify the service discovery from the client side. Also, the API Gateway must capable of handling the following requirements.

- Handling hundreds of thousands of requests at peak time.
- Rate limiting to prevent abuse from clients especially when the order ticket period is open.

Logging & authentication can be handled by each service because of certain needs such as ensuring no service can be accessed without proper authentication and complete log data.

I choose **Kong Gateway** as the API Gateway because of its performance, ease of use, and global rate limiting feature. The cons of this choice are:

- Redis will be the hot path.
- Higher memory usage than NGINX.
- Limited Open Source features such as no built-in dashboard.

### Backend Service Communication

For service-to-service communication, I will use synchronous HTTP communication for most of the cases because it is simpler to implement and debug. Adding message queue for asynchronous communication will be considered in the next ADR because the ticket order flow requires queue to distribute the load.

The gRPC can be considered in the future if I need better performance and strong typing between services, but tooling and debugging is more complex than HTTP.

### Result

![diagram](../asset/initial-architecture.svg)

## Consequences

### Positive

- Each service have clear bounded context that can be developed, deployed, and scaled independently.
- The level of granularity of each service is balanced between complexity and development effort.
- Using relational DB provide strong consistency and correctness of data.
- Using API Gateway simplify the client application communication to the backend services.

### Negative

- Using single DB can lead to potential bottleneck when the system scale. Also, schema change need to be coordinated between services.
- Synchronous HTTP communication can lead to cascading failure when one service is down. Each service need to implement proper circuit breaker and fallback mechanism.
- Using API Gateway can lead to single point of failure if not implemented properly.
- API Gateway can add extra latency for each request.
- More operational overhead because of many services that need to be monitored, deployed, and managed.
