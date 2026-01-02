# Creating Custom Web for Reseller

## Problem

To lower reseller technical burden, Ticketing Platform can provide hosting capability for reseller to create their own custom web to sell tickets instead of consuming the Ticketing Platform APIs. How I can implement this custom web feature?

## Solution

I can treat the custom web specification as a data instead of creating raw HTML, CSS, and JS files per reseller. The details of the web customization can be stored in the database such as color scheme, section position, text content, etc. Then, when a user accesses the reseller custom web, the server will render the web based on the stored configuration.

What about the backend? Instead of communicating with the Ticketing Platform APIs, I may need to create a dedicated backend service that will handle the business logic for the custom web. The Ticketing Platform APIs require certain authentication key (like how reseller system communicate with the platform.), so I don't want to expose that code key into the frontend.

This dedicated backend service will communicate with the Ticketing Platform APIs on behalf of the custom web. It will handle user registration, login, ticket browsing, ordering, payment processing, and real-time notifications. The custom web frontend will communicate with this backend service using its own set of APIs.
