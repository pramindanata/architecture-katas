# Notifying Reseller Web When Ticket Sold Out

## Problem

There is the following statement in the kata.

> User wants to be notified when a ticket they are looking at just got bought.

With custom web, the web can establish an SSE connection with the ticketing platform to receive the real time notification. But how I handle this kind notification if the reseller use our APIs to build their own web?

## Solution

Using a webhook is enough. When the ticket stock reach 0, system will send a POST request to the reseller webhook URL with certain payload. It will be the responsibility of the reseller to handle the notification and notify the user in their web.

Also, there is another approach where the reseller backend can establish an SSE connection directly into the ticketing platform. It sounds same as using webhook but I think it more complicated to maintain because it uses the same event as the custom web. If the API is changed, then reseller need to update their backend too. Also, reseller need to scale the SSE connection in their side if they have massive traffic.
