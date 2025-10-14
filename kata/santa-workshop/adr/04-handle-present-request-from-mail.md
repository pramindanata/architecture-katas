# ADR04 - Handle Present Request from Mail

Date: 2025-10-14

## Status

Accepted

## Context

Based on [ADR03](./03-standardize-present-request-analyzer.md), handling present requests require the system to retrieve the mail from a certain channel and analyze it using a dedicated parser service. We need an architecture that capable of handling both processes and it support the following criteria.

For retrieving the mail.

- Elves can upload the the scanned files into the system.
- The service that responsible for uploading the scanned files should capable of handling peak traffic without performance degradation.

For the service.

- It capable of retrieveing a text from an image containg a terrible handwritten mail.
- It capable of analyzing the retrieved text into structured format.

## Decision

### Retrieving the Scanned Files

The system must support a bulk upload feature that only accept an image files. This uploading process will be handled by a dedicated upload service because it need to scale well with the peak traffic. The upload service will upload the file into a certain cloud storage. After a file is uploaded, this service will inform the mail request parser about this new uploaded file via async communication.

Scanning the files need to be handled manually by the Elves. The elves can use the provided upload UI in the system to upload the files from the local storage.

There is another alternative, by make the scanner upload the file automatically to the cloud storage via upload service. Unfortunately, this approach requires extensive operational because it need to change the target URL frequently. The URL must be signed so no external entity can't upload files without proper authentication.

### Analyzing the Mail Image

Retrieveing text from the scanned file and analyzing the text requires AI solutions. We choose to develop our in-house AI models using existing models and fine-tune them to meet our needs. For retrieving the text from scanned file, we choose **TrOCR** model, while analyzing the text we choose **Llama 3.1 8B**.

Developing in-house solutions may requires huge investement upfront (including hiring ML engineers) but it is cheaper than using available Generative AI API. Based on our research, here are the cost estimation for both in-house solutions & using Generative AI API.

| Approach | Infrastructure Cost | Team Cost | Total Year 1 | Total Year 2+ | Notes |
|----------|-------------------|-----------|--------------|---------------|-------|
| **In-House** | **$1M-1.5M** | **$700K-1M** | **$1.7M-2.5M** | **$1.7M-2.5M** | Best ROI for seasonal traffic |
| Pure API (Vision LLM) | $50M-100M | $400K-600K | $50.4M-100.6M | $50.4M-100.6M | Prohibitively expensive |
| Pure API (OCR API + LLM) | $20M-67M | $400K-600K | $20.4M-67.6M | $20.4M-67.6M | Still too expensive |
| Hybrid (On-prem + Cloud) | $2M-3M | $900K-1.2M | $3.9M-5.2M* | $2.9M-4.2M | *Includes $1M capex |
| Year-Round Large Infrastructure | $18M-72M | $1.2M-1.8M | $19.2M-73.8M | $19.2M-73.8M | Massive waste |
| Prototype: Single Vision LLM | $100K-500K | $400K-600K | $500K-1.1M | $500K-1.1M | Good for MVP only |

For more details about choosing the AI solutions, refer to [this TRD](./../trd/handwritten-letter-ocr-at-scale.md).

### The Design

![diagram](./../asset/adr-04-diagram.svg)

## Consequences

### Positive

- Lower cost for maintaining the AI solution in long term.
- Allow fine tune the models based on our needs.
- In-house solution provides lower latency and better scaling.

### Negative

- Memory usage will be high because the upload & parser services need to process multiple image files.
- A new service need to be maintained.
- Need more time to develop the in-house AI solutions.
