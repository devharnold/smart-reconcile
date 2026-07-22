# Smart Reconciliation Platform (Go)

> **Architecture Goal:** Build a production-grade, event-driven payment reconciliation platform using idiomatic Go, minimal dependencies, and cloud-native principles.

---

# Guiding Principles

* The Go Standard Library whenever possible.
* Use small, focused libraries instead of full frameworks.
* Keep services loosely coupled through events.
* Design for observability from day one.
* Build services that can be independently deployed and scaled.
* Follow Clean Architecture principles where appropriate, but avoid over-engineering.
* Favor explicit code over "magic."

---

# Technology Stack

## API Layer

| Component   | Choice                    | Reason                              |
| ----------- | ------------------------- | ----------------------------------- |
| HTTP Server | `net/http`                | Built into Go, stable, fast         |
| Router      | `chi`                     | Thin wrapper around `net/http`      |
| Middleware  | `chi/middleware` + custom | Request ID, recovery, logging, CORS |

---

## Database

| Component        | Choice         |
| ---------------- | -------------- |
| Database         | PostgreSQL     |
| Driver           | pgx            |
| Query Generation | sqlc           |
| Migrations       | golang-migrate |

**Why?**

* Excellent PostgreSQL support
* Compile-time SQL validation
* No ORM overhead
* Type-safe queries
* Better performance

---

## Messaging

| Component | Choice       |
| --------- | ------------ |
| Broker    | Apache Kafka |
| Client    | franz-go     |

**Topics**

```
jobs.created

jobs.started

transactions.fetched

transactions.normalized

reconciliation.completed

reconciliation.failed

notifications.send

notifications.sent
```

---

## Caching

Redis

Use for:

* Session storage
* Idempotency keys
* Frequently accessed jobs
* Distributed locks
* Rate limiting

---

## Scheduling

```
robfig/cron
```

Responsibilities:

* Poll external providers
* Schedule reconciliation jobs
* Retry failed jobs
* Cleanup tasks

---

## Logging

```
zerolog
```

Structured JSON logging.

Every log should include:

```
timestamp
level
service
job_id
request_id
trace_id
message
```

---

## Metrics

Prometheus

Collect:

* API latency
* Kafka consumer lag
* Kafka publish failures
* Active jobs
* Failed jobs
* Worker execution time
* PostgreSQL latency
* Redis latency

Visualization:

Grafana

---

## Distributed Tracing

OpenTelemetry

Trace requests across:

```
API

↓

Kafka

↓

Workers

↓

Database
```

---

## WebSockets

```
nhooyr.io/websocket
```

Used for:

* Live reconciliation updates
* Dashboard notifications
* Job progress

---

# Project Structure

```
smart-reconcile/

├── cmd/
│   ├── api/
│   ├── scheduler/
│   ├── fetch-worker/
│   ├── normalize-worker/
│   ├── reconcile-worker/
│   ├── notification-worker/
│   ├── outbox-worker/
│   └── dlq-worker/
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   │
│   ├── reconciliation/
│   ├── normalization/
│   ├── fetching/
│   ├── notifications/
│   ├── scheduler/
│   ├── websocket/
│   │
│   ├── kafka/
│   ├── postgres/
│   ├── redis/
│   │
│   ├── config/
│   ├── logger/
│   ├── metrics/
│   └── tracing/
│
├── pkg/
│   ├── events/
│   ├── models/
│   ├── dto/
│   ├── errors/
│   └── utils/
│
├── sql/
│   ├── migrations/
│   └── queries/
│
├── configs/
│
├── scripts/
│
├── deployments/
│
├── docs/
│
└── Makefile
```

---

# Core Services

## 1. API Service

Responsibilities

* Authentication
* CRUD APIs
* Start reconciliation
* View jobs
* View reports
* WebSocket gateway

Consumes

```
HTTP
```

Produces

```
jobs.created
```

---

## 2. Scheduler

Responsibilities

* Schedule polling
* Retry jobs
* Cleanup
* Trigger reconciliation windows

Produces

```
jobs.started
```

---

## 3. Fetch Worker

Responsibilities

* Connect to payment providers
* Download transactions
* Store raw transactions

Consumes

```
jobs.started
```

Produces

```
transactions.fetched
```

---

## 4. Normalize Worker

Responsibilities

* Normalize provider-specific data
* Currency conversion
* Standardize timestamps
* Validate records

Consumes

```
transactions.fetched
```

Produces

```
transactions.normalized
```

---

## 5. Reconciliation Worker

Responsibilities

* Compare provider transactions
* Compare internal ledger
* Detect mismatches
* Calculate differences
* Produce reconciliation report

Consumes

```
transactions.normalized
```

Produces

```
reconciliation.completed

or

reconciliation.failed
```

---

## 6. Notification Worker

Responsibilities

* Send emails
* Push WebSocket updates
* Store notifications

Consumes

```
notifications.send
```

---

## 7. Outbox Worker

Responsibilities

Publish database events to Kafka.

Guarantees:

* No lost events
* Atomic database writes
* Reliable messaging

---

## 8. Dead Letter Queue Worker

Responsibilities

Handle messages that permanently fail.

Strategies:

* Retry
* Alert
* Archive
* Manual replay

---

# Event Flow

```
Scheduler
    │
    ▼
jobs.started
    │
    ▼
Fetch Worker
    │
    ▼
transactions.fetched
    │
    ▼
Normalize Worker
    │
    ▼
transactions.normalized
    │
    ▼
Reconciliation Worker
    │
    ├──────────────► reconciliation.completed
    │
    └──────────────► reconciliation.failed
                           │
                           ▼
                Notification Worker
```

---

# Job Lifecycle

```
QUEUED

↓

FETCHING

↓

FETCHED

↓

NORMALIZING

↓

NORMALIZED

↓

RECONCILING

↓

COMPLETED
```

Failure can occur at any stage.

```
FAILED
```

Each transition should be timestamped.

---

# Outbox Pattern

Instead of:

```
Insert reconciliation

Publish Kafka
```

Do:

```
BEGIN TRANSACTION

Insert reconciliation

Insert outbox_event

COMMIT
```

Outbox Worker

```
Read unpublished events

↓

Publish to Kafka

↓

Mark as published
```

Benefits:

* No lost messages
* Strong consistency
* Easier retries

---

# Database Tables

Core tables

```
providers

accounts

transactions

normalized_transactions

reconciliations

reconciliation_results

jobs

notifications

outbox_events

failed_messages
```

---

# Configuration

Environment variables

```
DATABASE_URL=

REDIS_URL=

KAFKA_BROKERS=

JWT_SECRET=

SMTP_HOST=

SMTP_PORT=

SMTP_USER=

SMTP_PASSWORD=

WS_PORT=

PROMETHEUS_PORT=
```

---

# Observability

Every request should include

```
Request ID

Trace ID

Job ID

Correlation ID
```

Logs

↓

Metrics

↓

Tracing

Should all point to the same identifiers.

---

# Security

* JWT Authentication
* Role-Based Access Control (RBAC)
* HTTPS everywhere
* Secrets via environment variables
* Input validation
* SQL injection protection (pgx/sqlc)
* Idempotent APIs for payment operations
* Audit logging for sensitive actions

---

# Testing Strategy

## Unit Tests

```
testing
```

Target:

> 80%+ coverage for business logic.

---

## Integration Tests

Test:

* PostgreSQL
* Kafka
* Redis

using Docker Compose or Testcontainers.

---

## End-to-End Tests

Validate:

* Job creation
* Fetch
* Normalize
* Reconcile
* Notification
* Dashboard updates

---

# Deployment

Containers

```
Docker
```

Orchestration

```
Kubernetes
```

Ingress

```
NGINX
```

CI/CD

```
GitHub Actions
```

Cloud

* AWS (EKS, RDS, MSK, ElastiCache)
* GCP (GKE, Cloud SQL, Memorystore)

---

# Future Enhancements

* Multi-tenancy
* Event sourcing
* CQRS for reporting
* Workflow retries with exponential backoff
* Provider plugin system
* gRPC internal communication
* OpenAPI documentation generation
* Horizontal autoscaling
* Data warehouse integration
* ML-powered anomaly detection
* Rule engine for reconciliation policies

---

# Why This Stack?

This architecture embraces Go's philosophy: **small, composable tools over heavyweight frameworks**.

By building on the standard library (`net/http`) with `chi`, the platform avoids unnecessary framework lock-in while retaining clean routing and middleware. Using `pgx` and `sqlc` keeps database access explicit and type-safe without the complexity of an ORM. Kafka provides reliable event-driven communication between workers, while the Outbox Pattern ensures database changes and published events remain consistent.

The result is a system that is easy to understand, test, scale, and maintain. Each service has a single responsibility, communicates through well-defined events, and can evolve independently. Observability, reliability, and operational simplicity are treated as first-class concerns from the beginning rather than retrofitted later.

This is the kind of architecture that scales well from an MVP to a production-grade financial platform and reflects the design patterns commonly found in modern Go backend systems.
