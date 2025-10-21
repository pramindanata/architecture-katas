# Mock Internet UN (United Nation)

Organization running "Mock UN" events wants to take its events online, permitting students to participate online

Requirements: student-diplomats must be able to video-chat with one another; student-diplomats must be able to "give speeches" to the "assembly" (video-chat to the entire group); (mocked) world events (created by moderators) distributed via (mock) "news sites"; moderators must be able to monitor any video chat for appropriateness

Users: 500 or so "diplomats" per "mock UN" gathering; dozens of moderators per "mock UN"; many "mock UN"s simultaneously; no new hardware requirements on students

## Discovery

### Actor

- Student: video chat to another student or group of students.
- Moderator: create event & monitor video-chat.

### Assumption

- A1: There less than 100 gatherings happened at the same time.
- A2: A Mock UN event contains many video-chat.
- A3: Only moderator can create a video chat session.
- A4: The event will be held every week.
- A5: The event will be held on a global scale.
- A6: In the video chat session, only max 25 participants can open mic & video.

### Functional Requirement

- R1: Moderator can create a Mock UN event.
- R2: Moderator can create a video chat session inside a Mock UN event. Also, he/she can invite other students.
- R3: Student and moderator can join a video chat.
- R4: Student and moderator can leave a video chat.
- R5: Student and moderator can speak and view video of each student camera video in a video chat.
- R6: Moderator can warn student in a video chat via messages in case the student breaks the chat rule.
- R7: Moderator can configure which students can open mic or muted.
- R8: Moderator can manage and publish "mock world events" within the Mock UN event.
- R9: Student can view published mock world event.
- R9: Student and moderator can send messages in video chat's chat column.
- R10: Student and moderator can see sent messages in the video chat's chat column.

### Non-Functional Requirement

- NFR1: Student must not be disconnected during the debate.
- NFR2: Moderator must be able to monitor live.
- NFR3: Dozens of simultaneous event, each up to 500 participants.
- NFR4: Live video and audio streaming must be low latency (under 300ms).
- NFR5: Students must not intrude into another events or chats.
- NFR6: Moderators need privileged access.
- NFR7: Multiple events run in parallel but only temporarily.
- NFR8: Student should not need to install extra software.
- NFR9: Video chat session can contain max 500 students and dozens of moderators. Total 600 users per session.

### Architecture Characteristic

- (TOP) Availability: support NFR1 & NFR2
- (TOP) Scalability: support NFR3 & NFR7
- (TOP) Performance: support NFR4 & NFR9
- Cost efficiency: support NFR7
- Security: support NFR5 & NFR6
- Usability: support NFR8
