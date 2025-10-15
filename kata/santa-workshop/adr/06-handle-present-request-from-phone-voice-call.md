# ADR06 - Handle Present Request from Phone Voice Call

Date: 2025-10-15

## Status

Accepted

## Context

With present requests sent into Santa via phone voice-call, the system need a mechanism to retrieve the voice-call information and analyze it into a format that can be understand by the rest of system. The system must also capable of retrieving and handling hundred thousand of voice-call per second at peak period.

## Decision

Use **human-assisted call handling (elves)** to receive and record present requests, instead of full automated call transcription. The elves later can input the request details into the system manually without additional analysis process like in mail & email requests.

Reasoning:

- Significantly lower operational cost at peak season.
- More accurate interpretation of children’s speech (accents, emotion, and context).
- Simpler to implement and maintain.
- Allows future extension to partial automation (e.g., AI transcription for voicemail only).

Here are samples of the estimated cost if an automation is implemented (assuming there are 100.000 calls and each duration is 5 minutes).

- Using Amazon Transcribe to only transcribe the audio call
  - Transcribe pricing = $0.024 per minute.
  - **Total cost = $9.750**
- Using Twillio Voice API to receive the call, record call, store audio, & transcribe audio (record, store, & transcribe are inseparable).
  - Receive call = $0.0085 per minute = $4.250
  - Record call = $0.0025 per min = $1.250
  - Storing call audio = $0.0005 per minute = $250
  - Transcribe call audio = $0.05 per min = $25.000
  - **Total cost = $31.750**

## Consequences

### Positive

- Reduce development & operation costs.
- Low impact for the system architecture.

### Negative

- Elves must handle the call, summarize the call, and input the summarize result into system manually.
