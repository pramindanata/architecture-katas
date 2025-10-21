# ADR03 - Handle Real Time Chat

Date: 2025-10-21

## Status

Accepted

## Context

Inside the video chat session, participants can chat or send messages in the video chat's chat column in real time. The architecture must ensure participants can chat with latency below 300ms.

## Decision

We choose the **Daily built-in APIs** to handle real time chat. We may need to persist the messages manually so participant who connect to a session can view the history of the chat.

Here are the options.

### Daily Built-In Message

Daily has APIs to send and receive messages in a video session.

Pros:

- No extra setup.

Cons:

- Message isn't persisted.

### WebSockets

A persistent TCP connection allowing both server → client and client → server messages instantly.

Pros:

- Truly real-time (tens of ms latency).
- Bidirectional, ideal for chat, presence, typing indicators.
- Integrates cleanly with SFU events (e.g., “user muted”, “moderator warning”).

Cons:

- Requires stateful connections (each client holds an open socket).
- Needs a load balancer with sticky sessions or session-aware gateway.
- Slightly higher infra cost (connections consume RAM + file descriptors).

### Polling

Client sends a request; server holds it open until a message arrives, then replies; client immediately re-requests.

Pros:

- Works everywhere, no special protocol.
- Can reuse existing HTTP infra.

Cons:

- Latency spikes (depends on polling cycle).
- Heavy server load (many pending requests).
- Inefficient with large concurrent users.

### SSE

One-way channel: server → client messages over a persistent HTTP connection.

Pros:

- Simple to implement (native `EventSource` in browser).
- Works through proxies/firewalls that block raw WebSockets.
- Lightweight, just text over HTTP.

Cons:

- One-way only: clients still need separate HTTP POSTs for sending messages.
- Not ideal for high fan-out (500+ users) since each connection is server-pushed.
- Less robust reconnection behavior than mature WebSocket libs.

## Consequences

### Positive

- No extra setup in infrastructure.

### Negative

- Vendor lock-in.

## References

- [Three ways to add chat to your video calls with the Daily API](https://www.daily.co/blog/three-ways-to-add-chat-to-your-video-calls-with-the-daily-api/)
