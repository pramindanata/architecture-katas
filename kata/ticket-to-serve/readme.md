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

## Discovery

### Actors

- Admin: resolving reseller support tickets.
- Reseller: have a system where it uses our APIs to manage its show & sell tickets.
- User: purchase ticket from reseller system or reseller custom web.
- Ticketing Platform: provide ticketing APIs.
- Reseller system: where users can see & purchase tickets. It uses the Ticketing Platform APIs for seeing & purchasing tickets.
- Reseller custom web: same as reseller system but it is hosted by the Ticketing Platform.
- Payment Provider: provide APIs for billing & payment.

### Assumptions

- A1: Ticket inventory is finite and must be pre-created by reseller.
- A2: Ticketing Platform is the source of truth for ticket availability and ownership.
- A3: When the ticket purchasing period of a highly anticipated show opens, the traffic can peak at a hundred of thousand purchases at a same time.
- A4: A show can have multiple type of tickets and each type of ticket will have its own stock.
- A5: Fraud detection will be handled by the selected payment provider.

### Functional Requirements

- R1: Reseller system can fetch list of owned shows via APIs.
- R2: Reseller system can fetch list of available tickets of an owned show via APIs.
- R3: Reseller system can create payment of an available ticket via APIs.
- R4: Reseller can send a support ticket.
- R5: Admin can process support tickets such as force refund, complain, etc.
- ~~R6: Reseller system can be notified certain events such as when a payment flagged as a fraud, payment paid, and looked ticket just get purchased~~.
- R7: Reseller can see detailed ticket sales report.
- R8: Reseller can manage shows and their tickets (CRUD).
- R9: Reseller can register an account to start managing shows & using the Ticketing Platform APIs. Including registering a billing information.
- R10: User can complete a ticket payment through a Payment Provider. Each complete purchase will be deducted by x% for the platform and payment provider.
- R11: Reseller system can request a payment refund.
- R12: When payment is created, the ticket will be reserved for certain duration (e.g. 15 minutes) until the payment is paid.
- R13: Reseller can create a custom web directly instead of consuming the APIs. Extra fees will apply.
- R14: Reseller can customize the custom web such as theme and text.
- R15: System can detect fraud payment.
- R16: Admin can mark shows that potentially will have high traffic when the ticket purchasing period opens.
- R17: User can register & login into the reseller custom web.
- R18: Reseller system can receive real time such as notification when looked ticket just get purchased.
- R19: Reseller system can receive notification such as payment paid & payment flagged as a fraud.

### Non-Functional Requirements

- NFR1: System must handle traffic spike at certain period especially for highly anticipted shows (see A3).
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
- Interoperability seems a priority because it affects reseller adoption, but it doesn't force fundamental architecture trade-off (scaling, locking, etc.) and only requires higher discipline in API design.
- Scalability is important, but it can be achieved by applying proper scaling technique such as sharding, partitioning, etc. which is a common solution.

## Spikes

- X Ticket purchase idempotency.
- X Ticket inventory consistency.
- X Purchase ticket flow
- auth mechanism
- Create custom web per reseller with their own theme & config (how?)
- Send ticket sold out notification in real time from ticketing to reseller system.
- Choosing payment provider (need support high traffic & fraud detection) (langsung ADR)
  - ugh kyknya beberapa provider engga publish angka. Alhasil perlu contact provider tsb atau bahkan perlu buat perjanjian.

## Design

- Choosing architecture (split domain but single DB)
- Handle ticket purchase (queue & notify)
- Handling payment
- Handling custom web
- Handling SSE from ticketing to reseller system (reseller system need to send SSE request manually to the ???)

## Reference

- [Scaling Ticket Booking Systems: A Deep Dive into High-Performance Design](https://blog.vineet.pro/blog/ticket-booking-system-design)
