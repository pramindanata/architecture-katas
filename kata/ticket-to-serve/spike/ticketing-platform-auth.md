# Ticketing Platform Auth

## Problem

How the ticketing platform authenticate the reseller backend request?

## Solution

There are some approaches that can be applied here. Personally I prefer API Key for its simplicity. OAuth 2.0 can be the alternative if I need more secure & complex auth mechanism.

### API Key

Each reseller will have its own API key that must be included in each request header. The ticketing platform will validate the API key and authorize the request based on the reseller associated with the API key.

### OAuth 2.0

The ticketing platform can implement OAuth 2.0 protocol where reseller need to register their application and obtain client ID and client secret. Reseller will use these credentials to obtain access token from the ticketing platform's authorization server. Each request to the ticketing platform APIs must include the access token in the Authorization header.
