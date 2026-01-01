# Read Me

This directory contains load test scripts and related files for performance testing of various components for the "Ticket to Serve" kata.

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [k6](https://k6.io/docs/getting-started/installation/) (for load testing)
- [Go](https://go.dev/doc/install) (for setting up some initial data)

## Installation

1. Run `docker-compose up` to start the required services.
2. Run the load test script using k6, e.g., `k6 run load-test-simple-redis-decrement.js`. See the "Scenarios" section for more details.

## Machine

- CPU: Intel(R) Core(TM) Ultra 5 125H (14 cores, 18 threads)
- RAM: 32 GB
- OS: Ubuntu 24.04 LTS

## Scenarios

### Redis Multi Decrement Using Lua Script

I'm curious to see if the Lua Script is really atomic and want to know if there are significant differences in performance if one check whether the value is below 0 or not. The result was that the Lua Script is indeed atomic and there were no significant performance differences either. The results were varied but both of them can reach ~20k RPS with 10k VUs & 1 million iterations.

Decrement without checking the value:

```bash
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./load-test-simple-redis-decrement.js
        output: -

     scenarios: (100.00%) 1 scenario, 10000 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 1000000 iterations shared among 10000 VUs (maxDuration: 10m0s, gracefulStop: 30s)

  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=501.3ms  min=11.29ms med=513.74ms max=1.27s p(90)=679.93ms p(95)=736.67ms
      { expected_response:true }...: avg=501.3ms  min=11.29ms med=513.74ms max=1.27s p(90)=679.93ms p(95)=736.67ms
    http_req_failed................: 0.00%   0 out of 1000000
    http_reqs......................: 1000000 18637.399289/s

    EXECUTION
    iteration_duration.............: avg=513.07ms min=33.66ms med=515ms    max=1.33s p(90)=719.76ms p(95)=792.44ms
    iterations.....................: 1000000 18637.399289/s
    vus............................: 108     min=108          max=10000
    vus_max........................: 10000   min=10000        max=10000

    NETWORK
    data_received..................: 123 MB  2.3 MB/s
    data_sent......................: 200 MB  3.7 MB/s
```

Decrement with checking the value:

```bash

         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./load-test-complex-redis-decrement.js
        output: -

     scenarios: (100.00%) 1 scenario, 10000 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 1000000 iterations shared among 10000 VUs (maxDuration: 10m0s, gracefulStop: 30s)

  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=499.06ms min=414.79µs med=544.73ms max=1.12s p(90)=677.31ms p(95)=756.21ms
      { expected_response:true }...: avg=499.06ms min=414.79µs med=544.73ms max=1.12s p(90)=677.31ms p(95)=756.21ms
    http_req_failed................: 0.00%   0 out of 1000000
    http_reqs......................: 1000000 19253.100114/s

    EXECUTION
    iteration_duration.............: avg=506.43ms min=36.02ms  med=545.02ms max=1.16s p(90)=695.91ms p(95)=799.4ms 
    iterations.....................: 1000000 19253.100114/s
    vus............................: 5537    min=5537         max=10000
    vus_max........................: 10000   min=10000        max=10000

    NETWORK
    data_received..................: 123 MB  2.4 MB/s
    data_sent......................: 201 MB  3.9 MB/s
```

To run this tests, follow these steps:

1. Run `go run ./cmd/setup-redis` to setup Redis with initial stock value.
2. Run `k6 run load-test-simple-redis-decrement.js` or `k6 run load-test-complex-redis-decrement.js` to run the load test.

### PostgreSQL `SELECT ... FOR UPDATE SKIP LOCKED`

With the sharding ticket solution that explained in [this doc](../../spike/purchase-ticket-flow.md), I want to see how well the query perform with a million of records.

With 1000 VUs & 100,000 iterations, it got ~3400 RPS with p95 latency ~380ms.

```bash

         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load-test-pg-reserve-ticket.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 100000 iterations shared among 1000 VUs (maxDuration: 10m0s, gracefulStop: 30s)

  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=294.92ms min=10.86ms med=315.95ms max=470.79ms p(90)=367.78ms p(95)=379.07ms
      { expected_response:true }...: avg=294.92ms min=10.86ms med=315.95ms max=470.79ms p(90)=367.78ms p(95)=379.07ms
    http_req_failed................: 0.00%  0 out of 100000
    http_reqs......................: 100000 3360.405034/s

    EXECUTION
    iteration_duration.............: avg=295.56ms min=11.87ms med=316.38ms max=471.24ms p(90)=368.29ms p(95)=379.6ms 
    iterations.....................: 100000 3360.405034/s
    vus............................: 1000   min=1000        max=1000
    vus_max........................: 1000   min=1000        max=1000

    NETWORK
    data_received..................: 15 MB  491 kB/s
    data_sent......................: 20 MB  655 kB/s
```

To run this test, follow these steps:

1. Run `go run ./cmd/setup-pg` to setup PostgreSQL with initial ticket units.
2. Run `k6 run load-test-pg-reserve-ticket.js` to run the load test.

To achieve this performance, some tuning in PostgreSQL are required such as.

#### Indexing

Applying a composite index that used in the `WHERE` clause of the query can significantly improve the performance to avoid extra sequential scan after index scan. I applied the index in `ticket_units(ticket_id, reserved)` as shown below:

```sql
CREATE INDEX idx_ticket_units_on_ticket_id_and_reserved ON ticket_units (ticket_id, reserved);
```

#### Sequential Scan Issue

When the load test was run for the first time, the p95 latency was around ~15ms. But after some tests, the latency spike to ~2s.

This spike was caused by PostgreSQL that decided to do sequential scan instead of index scan when the available data reach certain number. To solve this issue, I forced the PostgreSQL to always use index scan by using `SET LOCAL enable_seqscan = off;` before running the query in the transaction.

Some said using `enable_seqscan = off` is a bad practice because it can lead to worse performance in some cases. There are some alternatives such as configuring `random_page_cost`, but I haven't tried it yet.

#### Connection Pooling

High traffic requests require proper connection pooling to avoid overwhelming the database. In this load test, I apply the following configuration.

- PG's `max_connections` is set to 100 (default value).
- Client pool's max connection is set to 50.

Previously I set the pool's max connection to 1,000 without knowing that it should not exceed PG's `max_connections`, which caused why there are so many connection errors and slow test.

Remember that increasing the number of connection is not always the best solution because it can lead to connection thrashing. It's better to find the optimal number of connections based on the workload and database capacity. Also, more connections require more resources on the database server.

