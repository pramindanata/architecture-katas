# Ensuring Ticket Purchase Idempotency

## Problem

A user can accidentally or deliberately purchase a same ticket twice.

## Solution

Applying idempotency for a process requires cooperation for both client & server side. Client side must create a unique idempotency key (such as UUID) and server side must store that key in the database with a unique constraint. When a same request come to the server side with same key, the server side can throw an error because that key already exist in the DB.

The client side must guarantee that a same purchase have a same key. Sometimes a bug can occur where a client side call `n` purchase endpoints for same payload with different key. To avoid this kind of bug, client can create the key once when the page is loaded, before the purchase endpoint called.

In the server side, storing the key into DB must be atomic with other write operations. If any operation failed, all writes will be aborted.

A consumer (in this case is the server side) can't handle idempotency alone. Consumer don't know whether the received request is a different or same request. This problem requires both publisher (client side) & consumer (server side) to work together.
