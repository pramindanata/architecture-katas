# ADR05 - Scaling Notification Service

Date: 2026-01-03

## Status

Accepted

## Context

The notification service is responsible to send email notification to users and webhook notification to reseller system for certain events. During high traffic, this service need to be capable to scale efficiently, especially generating numbers of emails and send them to the users.

## Decision

### Generating Email Content

The email content will be generated using server-side templating engine such as Handlebars or EJS. This approach allows dynamic content generation based on user data and event type. This approach requires more server resources during high traffic, but it provides flexibility in email design and content.

### Email Provider

To handle email sending, I will use the **Amazon SES** as the email provider. Amazon SES is a scalable and cost-effective email service that can handle high volumes of email sending. It also provides features such as email tracking, bounce handling, and complaint management.

Amazon SES mainly focuses on basic sending email functionality and more complex to set up.

Considered alternatives are SendGrid & Mailgun. Both are popular email service with a robust API, good scalability, & have more features (such as marketing automation) but they are more expensive & not offer flexible pricing.

### Triggering Email Notification & Distributing Load

Other services that need to send email notification will publish an event to the existing **Apache Kafka** instance. The notification service will have a consumer that listens to those events and generate & send the email notification accordingly.

With how Kafka works, the load can be distributed to multiple notification service instances.

### Protocol for Webhook Notification

For webhook notification to reseller system, I will use **HTTP** protocol with JSON payload. This approach is widely adopted and easy to implement. The payload will contain relevant information about the event, such as ticket sold out or fraud detected.

The gRPC is considered as an alternative. It offers better performance and efficiency compared to HTTP/JSON, especially for high-throughput scenarios. However, it requires more complex setup, complex debugging, and may not be supported by all reseller systems.

### Result

![diagram](../asset/notification-service-architecture.svg)

## Consequences

### Positive

- The notification service can scale efficiently during high traffic periods.
- Using Amazon SES reduces the complexity of email sending and allows focusing on core functionalities.
- The use of Kafka as a queue helps to distribute the load and ensures reliable delivery of notifications.
- HTTP/JSON webhook notification is easy to implement and widely supported.

### Negative

- Generating email content on the server side requires more resources during high traffic.
- Relying on third-party email provider may introduce dependency and potential latency.
- HTTP/JSON webhook notification may have higher overhead compared to gRPC.
