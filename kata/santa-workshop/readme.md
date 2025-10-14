# Santa's Workshop Needs Your Help!

Santa gets millions of letters asking for presents each year. He needs help sorting through it all to make sure he fulfils each child’s wishes!

Many children have terrible handwriting, so the solution will need to be able to decipher the mail, as well as process emails, and voice phone calls, to generate toy orders for his elves to produce. These toy orders need to be consolidated into reports, as the elves need to know how many of each kind of toy needs to be produced.

The toy orders then have to be put together and stocked into Santa’s magic storeroom sack, with its specific location recorded along with child’s details.

While riding his sleigh, Santa needs locations of each home, what presents to deliver to which child, where in his sack he can find the presents, and any hazards he needs to be aware of when delivering these presents. Santa’s sleigh is as fast as lightning. He values speed and accuracy because the smallest mistakes can bring a child disappointment.

## Discovery

### Assumptions

- A1: A present request details consist of what toys the children want, his/her location, and his/her details. Each child must put this information into the mail, email, or voice call.
- A2: A present can contain multiple toys.
- A3: Duplicate requests may occurs from a same child. System will automatically consolidate them into 1 present.
- A4: In Christmas season period (2 weeks or 1.209.600 seconds), the present requests estimated to peak at 100K request per second.

### Actors

- Santa: deliver presents to children.
- Elves: view present order reports & produce presents.
- Children: send present requests in form of mail, email, or a phone call.

### Functional Requirements

- R1: System can retrieve present request email texts from an email server.
- R2: System can retrieve present request mail images from a physical scanner machine.
- R3: System can extract the mail text from present request mail images.
- R4: System can analyze present request texts (from mail, email, & vocie call) to extract the structured details (see A1 for the details specification).
- R5: Elves can view the report of toys that need to be produced.
- R6: Elves can update the present request production status after the toys production is done so it can be delivered.
- R7: Santa can view the delivery routes. Each delivery route dedicated for a single present location. It also contain information about the hazard that need to be aware of.
- R8: Santa can update the present request delivery status after the present is delivered.
- R9: System can retrieve recorder present request audio (or transcript) from a Voice Gateway service.
- R10: System can generate report of what toys need to be produced and their count from the available present requests.
- R11: Santa can view the GPS navigation of the current active delivery route in real time. It will be changed to the next route after the previous delivery is done.
- R12: System can mark a present request details to be reviewed manually when confident score is low (below 90%).
- R13: Santa can manual review present request details with low confident score.
- R14: Santa can view list of present that ready to be delivered.
- R15: Santa can start generating delivery routes for presents that ready to be delivered.
- R16: System can consolidate duplicate present requests from a same child into 1.

### Non-Functional Requirements

- NFR1: System should capable of analyze the requested present texts accurately.
- NFR2: System should capable of analyze the requested present texts in large number for shorter time (1.000 texts per minute).
- NFR3: System should ensure GPS accuracy when delivering the present.
- NFR4: System should maximize the overall delivery speed from the available delivery routes.
- NFR5: Analyzing text & delivery navigation features should have 99.9% availability during the Christmas Eve period.
- NFR6: System should capable of analyze tens-to-hundreds thousand of requested present during Christmas Eve period without performance degradation.
- NFR7: Children data must be handled securely
- NFR8: Communication between external services to retrieve the requested present information must be handled securely.

### Driven Architecture Characteristics

- (TOP) Accuracy: NFR1, NFR3, & NFR4
- Performance: NFR2
- (TOP) Availability: NFR5 & NFR6
- (TOP) Elasticity: NFR6
- Security: NFR7 & NFR8

## Implementation Concern

### Overview of Retrieving Present Request Text

> TODO

Gambar diagram supaya kebayang gimana peroleh request dari berbagai sumber (voice, mail, & email). Detailnya akan dijelaskan di section lain.

### Receiving Request from Phone Calls

> TODO

Use Voice Gateway to record phone call and send the text into our system. Example: Twillio Programmable Voice.

### Receiving Request from Mails

> TODO

Use OCR. Need to explore:

- Should we use existing model or train our own model to make system capable decipher bad handwriting from the images?

### Receiving Request from Emails

> TODO

- Tarik dari email server.

### Retrieving Present Details from Mail, Email, & Phone Call Texts

> TODO

Use AI. Hmm GPT should be enough.

### Routing Algorithm

> TODO

Santa deliver the present by flying. Need fast & accurate algorithm to generate the delivery routes.

### Viewing Delivery GPS Navigation

> TODO

Pakai short poll?

## Architecture

