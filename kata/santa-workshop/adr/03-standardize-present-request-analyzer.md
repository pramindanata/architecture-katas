# ADR03 - Standardize Present Request Analyzer

Date: 2025-10-14

## Status

Accepted

## Context

This system will receive present requests from different channels (mail, email, & voice call) and represented in a different format (e.g. scanned image, text, and audio). Then, the received requests will be analyzed to extract necessary details that needed by the system.

With different channels to retrieve the message and same mechanism to analyze the request, we need to set up an architecture that meet the following criteria:

- Accommodate multiple input channel with minimal coupling.
- Allow adding new channels in the future without major refactoring.

## Decision

We will introduce new services that share a same purpose, this service called **Request Parser Service**. This service responsible for.

1. Receiving/pulling data from a specific input channel.
2. Analyze the data into structured format by leveraging AI tools.
3. Publish the analysis result into other services via asynchronous communication.

![diagram](../asset/adr-03-diagram.svg)

Each channel (mail, email, and voice call) will have its own dedicated parser service. Each service will received different format of data, for example:

- The mail channel parser will receive/pull scanned mails from a physical machine.
- The email channel parser will receive/pull email texts from an email server.
- The voice call channel parser will receive/pull audio or its transcript from a Voice Gateway service.

The details of each parser will be explained in other ADRs.

The analysis result will be used by the main service in the system, so Elves can start produce toys and Santa can view the delivery details or review present requests manually (if the analysis process resulting a low confident level).

To scale better in high load, the received data from the channel will be published into Kafka and process it later instead of processing it immediately.

## Consequences

### Positive

- Adding a new channel by deploying a new parser service with minimal coupling.
- Scaling each channel parser service independently based on its traffic volume.

### Negative

- Increased number of services to manage and monitor.
- Requires consistent message schema and validation across all parser services.
