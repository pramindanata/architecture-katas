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

- Ticketing Platform: provide ticketing APIs.
- Admin: resolving reseller support tickets.
- Reseller: have a system where it use our APIs to manage its show & sell tickets.
- Reseller system: where users can see & purchase tickets. It uses the Ticketing Platform APIs for seeing & purchasing tickets.
- User: purchase ticket from reseller system.
- Payment Provider: provide APIs for billing & payment.

### Assumptions

- A1: Ticket inventory is finite and must be pre-created by reseller.
- A2: Ticketing Platform is the source of truth for ticket availability and ownership

### Functional Requirements

- R1: Reseller system can fetch list of owned shows via APIs.
- R2: Reseller system can fetch list of available tickets of an owned show via APIs.
- R3: Reseller system can create payment of an available ticket via APIs.
- R4: Reseller can send a support ticket.
- R5: Admin can process support tickets such as force refund, complain, etc.
- R6: Reseller system can be notified certain events such as when a payment flagged as a fraud, payment paid, and looked ticket just get purchased.
- R7: Reseller can see detailed ticket sales report.
- R8: Reseller can manage shows and their tickets (CRUD).
- R9: Reseller can register an account to start managing shows & using the Ticketing Platform APIs. Including registering a billing information.
- R10: User can complete a ticket payment through a Payment Provider. Each complete purchase will be deducted by x% for the platform and payment provider.
- R11: Reseller system can request a payment refund.
- R12: When payment is created, the ticket will be reserved for certain duration (e.g. 15 minutes) until the payment is paid.
- R13: Reseller can create a custom web directly instead of consuming the APIs. Extra fees will apply.
- R14: Reseller can customize the custom web such as theme and text.

### Non-Functional Requirements

- NFR1: User must be notified in real time when ticket they are looking at jus got bought.
- NFR2: System must handle traffic spike at certain period especially for high profiled shows (max to a million requests at a same time).
- NFR3: System must ensure the concistency of the ticket inventory.
- NFR4: System must ensure reseller can only see it's own data.
- NFR5: System must provide clear APIs contracts, backward compatibility, & versioning.
- NFR6: System must handle thousand of resellers where each can have thousands to millions users.
- NFR7: System must ensure a ticket can't be bought twice.

## Design

- Ticket purchase idempotency.
- Real time notification whe looked ticket just get bought.
- Create custom web per reseller with their own theme & config (how?)
- System must support per-region isolation of reseller and user data (??? bingung apa yg dipisah dan apakah ini berarti multitenancy. Skip aja untuk dipelajari?)

## Architecture Drives

## Domain Understanding

- Ticket lifecycle
- Payment lifecycle
- Reservation expiration
- Refund path

## High Level Design
