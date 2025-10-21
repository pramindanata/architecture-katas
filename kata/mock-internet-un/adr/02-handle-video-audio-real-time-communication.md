# ADR02 - Handle Video & Audio Real Time Communication

Date: 2025-10-20

## Status

Accepted

## Context

The system need to handle video & audio real time communication with the following criteria:

- Live video and audio streaming must be low latency (under 300ms).
- There will be multiple video chat sessions that take place at the same time.
- Each video chat session contain max 600 participants, but only some can have active audio and video.
- Each video chat session will be attended by participants globally.
- These video chat sessions take place in a single event and this event only take place once per week.

A WebRTC framework can be used to handle this kind of problem but the chosen WebRTC topology and our infrastructure need to be explored.

## Decision

### WebRTC Topology

For the WebRTC topology, we will choose **SFU**.

#### SFU (Selective Forwarding Unit)

Each participant sends one media stream to a server (SFU). The SFU forwards it to others, optionally downscaling/resizing streams.

```txt
     ┌──────────────┐
     │     SFU      │
     └──────────────┘
       ↙    ↓    ↘
 StudentA  B  Moderator
```

- Pros: Scales to hundreds of users per event, allows moderator to “monitor” rooms easily, & keeps latency low (100–300 ms).
- Cons: Server must handle high network I/O.

#### Mesh

Each participant sends video directly to every other participant.

```txt
Student A  ↔  Student B
     ↖      ↗
       Moderator
```

- Pros: Simple (no server load).
- Cons: Bandwidth heavy, not scalable beyond ~4 people.

#### MCU (Multipoint Control Unit)

The server mixes all streams into one composite video (like one grid view). Usually used for recording or broadcast.

- Pros: Reduces client CPU usage
- Cons: High server cost & latency

### Infrastructure

We choose **Daily** for our infra to handle the video chat.

> For estimating the infra cost, let's assume we have the following scenario in 1 month.
>
> - There are 100 active video chat sessions.
> - Each room has 500 participants and only 25 open mic & video for the whole of the session (active participants).
> - Viewer can only see 4 active participant streams at the same time.
> - Participants that can open mic and video can be changed throughout the session.
> - Each room take 120 minutes to complete the session.
> - Users around the globe can connect to a room with audio & video latency below 300ms.

#### Daily (Cloud)

<https://www.daily.co/>

Daily provides a cloud solution to handle video call with pricing of $0.004 per participant minute.

```txt
Estimated cost = roomCount * participantCount * durationInMinute * 0.004
    = 100 * 500 * 120 * 0.004
    = 6,000,000 * 0.004
    = $24,000
```

Pros:

- No setup extra infrastructure.
- Flexible pricing.
- Up to 1,000 active participants per room. Can support up to 100,000 participants but only 25 are active. The number can be increased by contacting the Daily team.
- Support video recording, storage, and AI transcription if needed in the future.
- Guaranteed SLAs.

Cons:

- More expensive in the long run.

#### Twillio Video (Cloud)

<https://www.twilio.com/en-us/video>

Twillio provides a cloud solution to handle video call with pricing of $0.004 per participant minute.

```txt
Estimated cost = roomCount * participantCount * durationInMinute * 0.004
    = 100 * 500 * 120 * 0.004
    = 6,000,000 * 0.004
    = $24,000
```

Pros:

- No setup extra infrastructure.
- Flexible pricing.
- Support video recording & storage if needed in the future.
- Guaranteed SLAs.

Cons:

- Only max 50 active participants per session.
- More expensive in the long run.

#### LiveKit (Cloud)

<https://livekit.io>

Estimated cost (assume we use the $500 tier):

```txt
Assume:

- Each of active participant has average bit-rate of ~1.5 Mbps.
- Each of viewer only receive 4 streams of active participants. Average of stream received by viewer is 6 Mbps.

WebRTC participant minutes = (roomCount * participantCount * durationInMinute - includedMinutes) * 0.0004
  = (100 * 500 * 120 - 1,500,000) * 0.0004
  = $1,800

Upstream cost is free.

Downstream cost = (roomCount * participantCount * (sessionDurationInSecond * receivedStream * MbToGBMultiplier) - includedSizeInGB) * 0.10
  = (100 * 500 * (7200 * 6 * 0.000125) - 3000) * 0.10
  = (50,000 * 5.4 - 3000) * 0.10
  = $26.700

Total = tier + participantMinutes + upstream + downstream
  = 500 + 1,800 + 0 + 26,700
  = $29,000
```

Pros:

- Tier include features like recording, streaming, & AI (LLM, STT, & TTS models).
- Guaranteed SLAs.

Cons:

- More expensive than other cloud providers for our use case.
- Tiered pricing with some features that don't needed by this system.
- Connecting to WebRTC, down streaming data, and others have different pricing.

#### LiveKit (Self-Host)

Cost estimation. It can vary based on each scenario. Also, the following calculation doesn't include the cost of hiring extra engineers to maintain the infrastructure.

| Solution | Monthly Cost |
|----------|--------------|
| **All 100 Sessions Concurrent (Once Per Month)** | **~$23,472** |
| **Sequential (~3 sessions/day, 6 hours peak)** | **~$28,200** |
| **Batched (10 sessions at a time, 10 batches/month)** | **~$23,943** |

For more details, see [this document](../other/livekit-self-host-cost-calculation.md).

Pros:

- It has lower infra cost and modern API than Jitsi.
- We have more control in infrastructure.
- Infra cost can be cheaper than cloud in the long run.

Cons:

- Need extra effort to develop and maintain the infrastructure.
- Non-infra cost like extra manpower and time to develop/maintain the infra can be higher than using the cloud solution.

#### Jitsi (Self-Host)

Cost estimation. It can vary based on each scenario. Also, the following calculation doesn't include the cost of hiring extra engineers to maintain the infrastructure.

| Solution | Monthly Cost |
|----------|--------------|
| **All 100 Sessions Concurrent (Once Per Month)** | **$24,280** |
| **Sequential (~3 sessions/day, 6 hours peak)** | **$29,008** |
| **Batched (10 sessions at a time, 10 batches/month)** | **$24,751** |

For more details, see [this document](../other/jitsi-self-host-cost-calculation.md).

Pros:

- More mature than LiveKit.
- We have more control in infrastructure.
- Infra cost can be cheaper than cloud in the long run.

Cons:

- Need extra effort to develop and maintain the infrastructure.
- Non-infra cost like extra manpower and time to develop/maintain the infra can be higher than using the cloud solution.

## Consequences

### Positive

- SFU has better performance to handle video chat with hundreds of participants.
- Daily offers flexible pricing and much higher participants count in each session than other cloud solution.
- We don't need to maintain the infrastructure.

### Negative

- Infrastructure cost can be higher in the long run.
- Vendor lock.
