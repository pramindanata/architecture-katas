# Room With a View

## Kata

A large hotel reservation company wants to build the next generation hotel reservation and management system specifically tailored to high-end resorts and spas where guests can view and reserve specific rooms.

Users: Guests (hundreds), hotel staff (less than 20)

Requirements:

- Registration can be made via web, mobile, phone call, or walk-in.
- Guests have the ability to either book a type of room (standard, deluxe, or suite) or choose a specific room to stay in by viewing pictures of each room and its location in the hotel.
- The system must be able to maintain room status (booked, available, ready to clean, etc.) as well as when the room will be needed next.
- It must also have state-of-the-art housekeeping management functionality so that cleaning and maintenance staff can be directed to various rooms based on priority and reservation need using proprietary devices supplied by the reservation company attached to the cleaning carts.
- Standard reservation functionality (e.g., payments, registration info, etc.) will be done by leveraging the existing reservations system.
- The system will be web-based and hosted by the reservation company.

Additional Context: 'Peak season is quickly approaching, so the system must be ready quickly or we have to wait until next year!'; Company is also investing heavily in cutting edge technology like smart room locks that open via a cell phone; Only interested in the high-end market; Sales people have tremendous clout in the organization; people often scramble to make their promises true.

## Requirements

### Functional

- R1: Guest registers into the system.
- R2: Guest logins into the system using its own credential (email).
- R3: Guest logins into the system using the temporary credential based on A4.
- R4: Guest views room list.
- R5: Guest create reservation a room by type or choose specific room based on its location in the hotel.
- R6: Guest receives reservation confirmation after payment is done. It can contain temporary credential for A4.
- R7: Guest receives "complete payment" link after the reservation is created manually by a reservation staff.
- R8: Based on A4, guest can open/close room smartlock through the system web/mobile app.
- ~~R9: Reservation staff proceed guest reservation (accept/deny)~~ (cancelled because user already paid the reservation).
- R10: Reservation staff create reservation for a guest (walk-in/phone-call reservation).
- ~~R11: Reservation staff assign cleaning staff to do housekeeping for a certain room.~~
- R12: Reservation staff view housekeeping tasks.
- R13: Reservation staff can revoke the guest temporary credential (see A4).
- R14: Reservation staff can re-genereate the guest temporary credential (see A4).
- R15: Cleaning staff views assigned housekeeping tasks.
- R16: Cleaning staff receives assigned housekeeping task notification via in-app notification.
- R17: Cleaning staff update housekeeping status.
- R18: Guest request housekeeping (see A5).
- R19: Reservation staff re-generate "complete payment" link for a room with an expired payment.
- R20: System assign housekeeping request to cleaning staff based on the room type priority.
- ~~R21: Reservation staff can lock/unlock a room smart lock.~~ (see A12).
- R22: Reservation/cleaning staff login into the system (see A10).
- R23: Reservation staff see failed lock/unlock room logs.

### Non Functional

- N1: New system need handle massive traffic properly on peak season (see A2 for the assumed traffic).
- N2: New system must communicate properly with the existing reservation system & smart lock system.
- N3: Communication between new service, existing reservation system, and smart lock system must be secure (API auth & encryption).
- N4: Ensure smart lock room can only be accessed by the guest who own the room.
- N5: Apply expiration to the guest temporary credentials.
- N6: New system must maintain 99.9% availability in peak season.
- N7: Room state (reserved, cleaned, locked/unlock, etc) must be sync in each systems within 1 minute.
- N8: Architecture should favor rapid development because the new system is needed in the next year.
- N9: New system must provide monitoring for certain issues such as fail lock/unlock smart lock.

### Assumption

- A1: The existing reservation system also handle the following functionalities:
  - Manage room data.
  - Update room status manually (available, ready to clean, etc).
- A2: The peak seasons can happen once each month and it can reach 1 million guests and 10K concurrent requests.
- A3: Payment for the guest reservation will created and completed in the existing reservation system.
- A4: The room smart lock is connected to the new system backend so opening it can be done using the new system web or mobile app. For a guest that book manually (walk-in/phone-call), they will receive a temporary credential to login into the web.
- ~~A5: When guest request housekeeping, the request will be automatically assigned to the available cleaning/maintenance staff using round-robin approach.~~ (see R20)
- A6: In R5, R10, R11, R17, & R18, the new system will sync room state into the existing reservation system. The existing reservation system will be the source of truth of the room data.
- A7: In R8 & R21, new system will communicate to smart lock system to lock/unlock the door and sync the state into the existing reservation system.
- A8: In R3 & R6, the credential only active based on the guest check in/out period.
- A9: In R7, the complete payment link will contain a long-unique token that represent the payment resource. Also the link will be expired within 24 hours.
- A10: Reservation & cleaning staff use SSO to login to the new service. Their user credentials (source of truth) is located in the existing reservation system.
- A11: The smart lock system is hosted by the reservation company instead of using service provided by external entity.
- A12: Reservation staff can lock/unlock room using existing reservation system.
- A13: 1 reservation can contain multiple rooms.
- A14: Each room will have its own category or type. A room category own `n` number of rooms.
- A15: When a guest create a reservation by the room type, system will automatically select specific room(s) for the guest immediately before payment is created.

## Architecture Characteristics

### Driving Characteristic

- (TOP) **Elasticsity**: N1
- (TOP) **Security**: N3, N4, & N5
- **Interoperability**: N2 & N7
- (TOP) **Feasibility**: N8
- **Availability**: N6
- **Observability**: N9

### Others Considered

None

## Architecture

### Details

The selected architecture style is **Service Based**. It favor cheap and fast development but also provide room for the system elasticsity. There are 2 services; booking & reservation services.

Booking service a dedicated service to handle user booking requests especially at peak season. Separating it will maximalize the service performance at the peak time.

Th reservation service will handle the rest of the features. This service will act like a monolith service. No need for seperate it into multiple services because the need of rapid development and no performance urgency in the other features except the booking feature.

Both services will have own DB without shared access to ensure elasticsity & availability especially for the booking service.

Each services will communicate using both sync and async communication based on certain activities. Details:

- Sync communication (HTTP request) will be used for activities that require immediate result lock/unlock door in smart lock system. Because these activities has low impact when peak season occurs, this communication is okay to use.
- Async communication (Pub/Sub) will be used to handle high traffic request in peak season such as create payment & update room state when a booking request is created. This communication is more expensive to build & doesn't give immediate result but it allow system to handle massive traffic by queueing the them.

### Diagram

#### C4 (Context)

![diagram](./c1-diagram.svg)

#### C4 (Container)

![diagram](./c2-diagram.svg)

## Implementation Concern

### Payment Creation in Peak Season

Payment creation is done in the old system. When peak season occurs, the legacy system will be flooded with massive traffic. To reduce the traffic, we employed pub/sub as communication mechanism between the new and old system. With pub/sub, the traffic can be queued in old system but a new issue arises, guest will not receive the payment information immediately after a booking is created.

To fix this issue, there are 2 solutions.

1. Send an email containing payment information to the guest.
2. Old system dispatch topic to new system when payment created and use short-polling in client side to fetch the payment info.

The first solution is the simplest but it has bad UX. The guest need to wait for the email and use the link inside the email to complete the payment.

The second solution is more complicated but it offer better UX. Client application need to show a loading state after booking is created. When the loading occurs, client application will use short polling to fetch the payment data in new system. We store payment data in the new system to reduce traffic in the old system.

Second solution is preferred because it still feasible to develop and reduce confusion in the guest side.

To improve short-polling performance, we can introduce distributed cache like Redis or Memcached to store the payment info temporarily instead fetching in from the DB.

Here is the flow diagram.

![diagram](./payment-creation-in-peak-season.svg)

### Reservation concurrency in Peak Season

> Please read A13, A14, & A15 regarding the assumption for the room requirements.

Handling reservation creation in massive traffic can cause a same room being reserved by multiple guests. To avoid this, a special concurrency process must be done.

The system will store all room IDs in the cache and when a guest create a reservation, system will apply optimistic lock for selected rooms. When a guest reserve booked rooms, system will abort the reservation process. The lock will be released when the guest checked-out or not completing the payment.

The reservation payment duration must be not long to reduce number of rooms getting locked by guests who not complete the payment. The duration can be 15 minutes.

When a guest create a reservation by choosing the room type, system will select available room IDs randomly from cache and apply locks.

Because the available room IDs are taken from cache, warm-up and sync mechanism is needed to ensure the cache is empty. The warm-up and sync process can be done when the booking service is started and periodically using a scheduler. When a sync process is triggered, the reservation process must wait the sync process to be done. The waiting process can be done by checking whether a "sync lock" is exist in cache and wait until it disappeared.

### Housekeeping Priority

When a guest request a housekeeping service, system need to prioritize rooms based on the room type. With limited number of staff, it must ensure the prioritization and housekeeping assignment are done properly. This housekeeping feature will be done using event-driven & queue.

When certain events are triggered (guest checked-out, housekeeping requested, etc), the event handler will create a `housekeeping task`object. This object contains a prioritization score that were calculated from the room type. When the task is created, system will dispatch a job to queue to assign the cleaning staff to that task.

The queue must have concurrency of 1 (it means only 1 job can be processed at the same time) to avoid system accidentally assign same staff to different rooms at the same time when there are multiple dispached "task created" jobs.
