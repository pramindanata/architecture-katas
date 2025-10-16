# ADR07 - Retrieve Coordinate from an Address

Date: 2025-10-15

## Status

Accepted

## Context

The given mails, emails, & voice calls only include the address information in a text format. System need to know the actual coordinate of those addresses, so system can generate the delivery routes properly. The solution must capable of processing large number of addresses in peak period and ensure the accuracy of the location.

## Decision

We choose **Google Maps API** to handle the geocoding process (retrieving coordinate from an address). It covers more complete addresses than other alternatives, so it ensures the accuracy of the address location and reduce number of addresses not found. Unfortunately, it has some major issues:

- Higher cost ($5 per 1,000 calls)
- Lower rate limit (50 RPS per GCP project)
- Restriction about caching the query result, for example, geocoding result can be cached only for 30 days.

To handle the rate limit issue, we can try to contact Google Maps sales to establish an enterprise contract or use multiple projects. With this low RPS, our system need to fetch the address coordinate as many as it can while waiting until toys production is complete.

A **distributed cache** will be added to reduce the API calls in case some present came from same addresses. The benefit of this cache maybe low because the majority of present requests are unique.

Other alternatives:

- [Pelias](https://github.com/pelias/pelias) (self-hosted): less address data. We need to use open source data (for example OpenStreetMap) or create it manually.
- [Geocode Earth](https://geocode.earth/) (cloud-managed): less address data because it uses open source data and lower RPS (20 RPS).

## Consequences

### Positive

- Google Maps API has more address dataset, so it maximizes the accuracy of address location & reduce not found address.

### Negative

- Higher cost to call the API.
- Low RPS to fetch the address coordinate.
- Some restriction in storing the API result.

## Reference

- [Google Maps API Cache Policy](https://developers.google.com/maps/documentation/geocoding/policies#cache-policy)
