# ADR05 - Handle Present Request from Email

Date: 2025-10-15

## Status

Accepted

## Context

With present requests sent into Santa email address, the system need a mechanism to retrieve inbound emails and analyze it into a format that can be understood by the rest of system. The system must also capable of retrieving and handling many hundred thousand emails per second at peak period.

## Decision

We choose **Amazon SES** to handle the retrieve inbound emails. It offers flexible pricing and can scale better than other alternatives. We also add an **Email Request Parser Service** to analyze the email content.

![diagram](../asset/adr-05-diagram.svg)

Using SES to retrieve inbound emails also mean we will use **S3 & SNS**. S3 will be used to store the raw email MIME (including the email body) while SNS will be used to send the email metadata into our system. With retrieved email metadata, our system can fetch the raw email MIME later from the S3.

For analyzing the email body, system will use the Parser AI model that explained in [ADR5](./05-handle-present-request-from-email.md).

With SES, here are the cost estimation assuming in a month our system process 3 million of emails and the size of each email is 75 KB.

| Component                        | Unit Price                            | Usage                               | Monthly Cost |
| -------------------------------- | ------------------------------------- | ----------------------------------- | ------------ |
| **Inbound email processing**     | $0.10 / 1,000 emails                  | 3,000,000                           | **$300.00**  |
| **Data transfer (75 KB/email)**  | $0.09 / GB                            | 3,000,000 × 75 KB ≈ 214 GB          | **$19.26**   |
| **SNS delivery (HTTP)**          | $0.50 / 1M deliveries (first 1M free) | 3M                                  | **$1.00**    |
| **S3 storage**                   | $0.023 / GB                           | 214 GB                              | **$4.92**    |

Total: $320.26/month

Here are alternatives for retrieving the emails.

- CloudMailin: it has simpler flow on sending the email information, but the pricing is more expensive ($800 per month for 2 million emails).
- MailerSend: it can only process max 500,000 emails daily. Also, the pricing package include other features that not needed like sending emails, tracking, etc.
- Set up our own infrastructure: there are many needed to be developed such as SMTP server, load balancing, anti-spam, email retry, etc.

## Consequences

### Positive

- High scalability and availability with minimal operational effort.  
- Pay-as-you-go pricing model — cost scales with actual load.  
- Tight integration with S3/SNS simplifies event-driven workflows.  

### Negative

- Reliance on AWS ecosystem (vendor lock-in).
- Slightly higher integration complexity (multiple AWS services involved).
- A new service need to be maintained.
