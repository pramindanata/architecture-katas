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
- R5: Guest books a room directly.
- R6: Guest receives book confirmation after payment is done. It can contain temporary credential for A4.
- R7: Guest receives "complete payment" link after the book is created manually by a reservation staff.
- R8: Based on A4, guest can open/close room smartlock through the system web/mobile app.
- R9: Reservation staff proceed guest books (accept/deny).
- R10: Reservation staff book a room for guests.
- R11: Reservation staff assign cleaning staff to do housekeeping for a certain room.
- R12: Reservation staff view housekeeping tasks.
- R13: Reservation staff can revoke the guest temporary credential (see A4).
- R14: Reservation staff can re-genereate the guest temporary credential (see A4).
- R15: Cleaning staff views assigned housekeeping tasks.
- R16: Cleaning staff receives assigned housekeeping task notification via in-app notification.
- R17: Cleaning staff update housekeeping status.
- R18: Guest request housekeeping (see A5).
- R19: Reservation staff re-generate "complete payment" link for a room with an expired payment.
- R20: System assign housekeeping request to cleaning staff using round-robin.
- R21: Reservation staff can lock/unlock a room smart lock.
- R22: Reservation/cleaning staff login into the system (see A10).

### Non Functional

- N1: New system need handle massive traffic properly on peak season (see A2 for the assumed traffic).
- N2: New system must communicate properly with the existing reservation system & smart lock system.
- N3: Communication between new service, existing reservation system, and smart lock system must be secure (API auth & encryption).
- N4: Ensure smart lock room can only be accessed by the guest who book the room.
- N5: Apply expiration to the guest temporary credentials.
- N6: New system must maintain 99.9% availability in peak season.
- N7: Room state (booked, cleaned, locked/unlock, etc) must be sync in each systems within 1 minute.
- N8: Architecture should favor rapid development because the new system is needed in the next year.
- N9: New system must provide monitoring for certain issues such as fail lock/unlock smart lock.

### Assumption

- A1: The existing reservation system also handle the following functionalities:
  - Manage room data.
  - Update room status manually (available, ready to clean, etc).
- A2: The peak seasons can happen once each month and it can reach 1 million guests and 10K concurrent requests.
- A3: Payment for the guest book will be done in the existing reservation system. New system will provide API to fetch the needed information for the payment activities.
- A4: The room smart lock is connected to the new system backend so opening it can be done using the new system web or mobile app. For a guest that book manually (walk-in/phone-call), they will receive a temporary credential to login into the web.
- ~~A5: When guest request housekeeping, the request will be automatically assigned to the available cleaning/maintenance staff using round-robin approach.~~ (See R20)
- A6: In R5, R10, R11, R17, & R18, the new system will sync room state into the existing reservation system. The existing reservation system will be the source of truth of the room data.
- A7: In R8 & R22, new system will communicate to smart lock system to lock/unlock the door and sync the state into the existing reservation system.
- A8: In R3 & R6, the credential only active based on the guest check in/out period.
- A9: In R7, the complete payment link will contain a long-unique token that represent the payment resource. Also the link will be expired within 24 hours.
- A10: Reservation & cleaning staff use SSO to login to the new service. Their user credentials (source of truth) is located in the existing reservation system.
- A11: The smart lock system is hosted by the reservation company instead of using service provided by external entity.

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

The selected architecture style is **Service Based**. It favor cheap and fast development but also provide room for the system elasticsity. There are 2 services; reservation & booking services.

Booking service a dedicated service to handle user booking requests especially at peak season. Separating it will maximalize the service performance at the peak time.

Th reservation service will handle the rest of the features. This service will act like a monolith service. No need for seperate it into multiple services because the need of rapid development and no performance urgency in the other features except the booking feature.

Both services will have own DB without shared access to ensure elasticsity & availability especially for the booking service.

Each services will communicate using both sync and async communication based on certain activities. For activities that require immediately result like create payment after booking request is created, it will use sync communication. For activities that need to fan out the communication to multiple services like sync room lock state to smart lock and existing system, it will use async communication.
