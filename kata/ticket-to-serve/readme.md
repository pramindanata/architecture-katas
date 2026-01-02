# Ticket to Serve

Your goal is to provide a centralized ticketing API for every ticket seller to use under the hood of their system.

You need to provide a various set of feature to help your client manage their show and sell their ticket, but there is also a bunch of feature that may be over scope

Reseller: about a few thousand, thousands to million user per reseller

Requirements:

- You need to provide a system for your client so their user can get a list of available ticket for a show and buy them
- User wants to be notified when a ticket they are looking at just got bought
- You have to absolutely make sure no ticket get bought twice
- You want to get yourself an admin interface to help reseller support with specific operation
- You should think about implementing a fee system over each purchase
- You want to provide an event interface to notify reseller (Like some payment being flagged as fraudulent several day after the payment)
- You need to provide a clear accounting system providing enough information to your reseller so they know how much they are going to make

Bonus:

- You may want to provide some hosting capability to your reseller to lower their technical burden
- As some country have constraining legislation, you want to keep in mind the possibility of completely separating your client or user base for specific legislation

> I skip the last bonus part because I need to explore more information about the constraint.

## Discovery

### Actors

- Reseller: have a system where it uses our APIs to sell the tickets or host a custom web.
- User: purchase ticket from reseller's system or reseller's custom web.
- Admin: resolving reseller support tickets.

### Assumptions

- A1: Ticket inventory is finite and must be pre-created by reseller.
- A2: Ticketing Platform is the source of truth for ticket availability and ownership.
- A3: When the ticket purchasing period of a highly anticipated show opens, the traffic can peak at a hundred of thousand purchases at a same time.
- A4: A show can have multiple type of tickets and each type of ticket will have its own stock.
- A5: Fraud detection will be handled by the selected payment provider.
- A6: Posted shows are like concerts, sport events, etc. that only held at specific location, date, & time.

### Functional Requirements

#### Ticketing Platform APIs

- Platform can return list of reseller's shows.
- Platform can return list of available show's tickets.
- Platform can create order & payment of the selected tickets.
- Platform can create request refund.
- Platform can receive notification from payment provider such as payment success, failure, refund processed, fraud detected, etc.
- Platform can send email notification to user such as ticket purchased, refund processed, fraud detected, etc.
- Platform can return queue status for the user order.
- Platform can return created orders.
- Platform can send real time notification to reseller system such as ticket sold out, fraud detected, etc. via webhook.

#### Reseller Dashboard

- Reseller can register an account to start managing shows & using the Ticketing Platform APIs. Including billing information.
- Reseller can manage shows and their tickets.
- Reseller can see detailed ticket sales report.
- Reseller can send a support ticket.
- Reseller can create & customize a custom web instead of consuming the APIs. Extra fees will apply.

#### Custom Web

- User can register & login into the reseller custom web.
- User can explore available shows.
- User can order tickets.
- User can complete ticket payment.
- User can receive real time notification when the ticket they are looking at just got purchased.
- User can see his/her queue status.
- User can see his/her created orders.

#### Admin Dashboard

- Admin can process support tickets such as force refund, complain, etc.
- Admin can mark shows that potentially will have high traffic when the ticket purchasing period opens.

### Non-Functional Requirements

- NFR1: System must handle traffic spike at certain period especially for highly anticipated shows (see A3).
- NFR2: System must ensure strong consistency and correctness of ticket inventory (no double purchase).
- NFR3: System must ensure reseller can only see its own data.
- NFR4: System must provide clear APIs contracts, backward compatibility, & versioning.
- NFR5: System must handle a thousand of resellers where each can have thousands to millions users.
- NFR6: System must have 99.9% availability when the ticket purchasing period opens.

## Architecture Drives

- (Top) Consistency: NFR2
- (Top) Availability: NFR6
- (Top) Elasticity: NFR1
- Scalability: NFR5
- Interoperability: NFR4
- Security: NFR3

Note:

- Ensuring ticket stock correctness while handling massive traffic spike when the ticket purchasing period opens is the main challenge for this system.
- Ticket oversell need to be avoided because it will lead to potential user dissatisfaction and bad reputation to the company. Unlike airline ticketing where oversell is common, airlines can accommodate more passengers due to no-shows or cancellations.
- Interoperability seems a priority because it affects reseller adoption, but it doesn't force fundamental architecture trade-off (scaling, locking, etc.) and only requires higher discipline in API design.
- Scalability is important, but it can be achieved by applying proper scaling technique such as sharding, partitioning, etc. which is a common solution.

## Spikes

- X Ticket purchase idempotency.
- X Ticket inventory consistency.
- X Purchase ticket flow
- X Send ticket sold out notification in real time from ticketing to reseller system.
- X Auth mechanism
- X Create custom web per reseller with their own theme & config (how?)

## Design

- Choosing architecture (split domain but single DB)
- Handle ticket purchase (queue & notify)
- Handling payment
- Handling custom web
- Handling SSE from ticketing to reseller system (reseller system need to send SSE request manually to the ???)

## Reference

- [Scaling Ticket Booking Systems: A Deep Dive into High-Performance Design](https://blog.vineet.pro/blog/ticket-booking-system-design)
