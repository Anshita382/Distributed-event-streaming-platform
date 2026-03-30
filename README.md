# Planet-Scale Event Streaming Platform

A GitHub-ready distributed systems project that mimics the kinds of event pipelines used in Amazon- and Netflix-style backend platforms.

It demonstrates how to ingest, publish, process, retry, deduplicate, persist, observe, and stress-test high-volume event flows using **Kafka + Go + Java + Redis + PostgreSQL**.

## Why this repo matters

This project is intentionally designed to signal strong backend and infrastructure fundamentals:

- distributed event-driven architecture
- producer/consumer workflow design
- idempotent processing
- retry queues and DLQ handling
- backpressure-aware consumers
- service observability with Prometheus and Grafana
- failure testing with container restarts
- clean multi-service repo organization for real engineering teams

## Architecture

```text
                +----------------------+
                |  Java Core Service   |
                |  Spring Boot API     |
                +----------+-----------+
                           |
                           v
+-------------+      +-----------+      +------------------+
| Go Producer | ---> |   Kafka   | ---> |   Go Consumer    |
| Load Source |      |  topics   |      | processing layer |
+-------------+      +-----------+      +---------+--------+
                                                 / \
                                                /   \
                                               v     v
                                           +------+ +-----------+
                                           |Redis | |PostgreSQL |
                                           |dedupe| |durability |
                                           +------+ +-----------+

Failures -> events.retry -> re-consume -> events.dlq after max retries
Metrics  -> Prometheus -> Grafana dashboards
```

## Tech stack

- **Java 21 + Spring Boot** for core ingest APIs
- **Go 1.22** for high-throughput producer and consumer services
- **Apache Kafka** for streaming backbone
- **Redis** for fast idempotency checks
- **PostgreSQL** for durable processed-event records
- **Prometheus + Grafana** for throughput and latency dashboards
- **Docker Compose** for local multi-service orchestration
- **GitHub Actions** for CI

## Core features

### 1. Event ingestion
The Java service exposes REST APIs that publish events into Kafka.

### 2. High-throughput producers
A Go load producer continuously emits events to simulate burst traffic.

### 3. Consumer-side idempotency
The Go consumer uses Redis `SETNX` plus PostgreSQL uniqueness guarantees to avoid double-processing.

### 4. Retry queue and DLQ
Failures are re-routed to `events.retry`. After the configured retry ceiling, they are pushed to `events.dlq`.

### 5. Backpressure-aware settings
Consumer configuration uses bounded buffering and controlled processing behavior to better mimic real streaming services.

### 6. Observability
Each service exposes metrics for Prometheus scraping, with Grafana dashboards provisioned in the repo.

### 7. Fault injection
A helper script can stop, kill, or restart the consumer so you can observe failure and recovery behavior.

## Repository structure

```text
planet-scale-event-streaming/
├── .github/workflows/ci.yml
├── docker-compose.yml
├── Makefile
├── .env.example
├── java-core-service/
│   ├── Dockerfile
│   ├── pom.xml
│   └── src/
├── go-producer/
│   ├── Dockerfile
│   ├── go.mod
│   └── ...
├── go-consumer/
│   ├── Dockerfile
│   ├── go.mod
│   └── ...
├── ops/
│   ├── prometheus.yml
│   ├── grafana/
│   └── postgres/
└── scripts/
    ├── create-topics.sh
    └── fault-injection.sh
```

## Local run

### Option A: one-command startup

```bash
make up
```

This starts:
- Kafka + Zookeeper
- Redis
- PostgreSQL
- Java ingest service
- Go producer
- Go consumer
- Prometheus
- Grafana

### Option B: start manually

```bash
docker compose up -d zookeeper kafka redis postgres prometheus grafana
bash scripts/create-topics.sh
```

Then run each app locally:

```bash
cd java-core-service && mvn spring-boot:run
cd go-consumer && go run ./cmd/consumer
cd go-producer && go run ./cmd/producer
```

## Create Kafka topics

```bash
bash scripts/create-topics.sh
```

Topics created:
- `events.main`
- `events.retry`
- `events.dlq`

## APIs

### Publish one event

```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "eventId":"evt-1001",
    "eventType":"order.created",
    "partitionKey":"customer-42",
    "payload":{"amount":199,"source":"manual-test"},
    "retryCount":0,
    "createdAt":"2026-03-29T12:00:00Z"
  }'
```

### Generate load

```bash
curl -X POST "http://localhost:8080/api/events/load-test?count=10000"
```

### Health check

```bash
curl http://localhost:8080/api/events/health
```

## Observability

- Grafana: `http://localhost:3000`
- Prometheus: `http://localhost:9090`
- Java metrics: `http://localhost:8080/actuator/prometheus`
- Go consumer metrics: `http://localhost:2112/metrics`
- Go producer metrics: `http://localhost:2113/metrics`

Grafana default credentials:
- user: `admin`
- password: `admin`

## Fault injection

Restart the consumer:

```bash
bash scripts/fault-injection.sh go-consumer restart
```

Kill the consumer:

```bash
bash scripts/fault-injection.sh go-consumer kill
```

Stop the consumer:

```bash
bash scripts/fault-injection.sh go-consumer stop
```

