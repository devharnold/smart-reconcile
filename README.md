# smart-reconcile

> **Financial transaction reconciliation infrastructure for teams that can't afford to be wrong.**

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![GCP](https://img.shields.io/badge/infra-GCP-4285F4?style=flat-square&logo=google-cloud)](https://cloud.google.com)
[![Status](https://img.shields.io/badge/status-active%20development-orange?style=flat-square)]()

---

Most fintech systems break in the same places: a Stripe webhook fires twice, M-Pesa returns a timezone-shifted timestamp, a bank CSV has a malformed row, and suddenly your ledger is off by KES 3,000 with no trail explaining why.

**smart-reconcile** is a backend platform built to prevent exactly that. It ingests transactions from multiple payment providers, stores raw payloads before touching them, normalizes everything into a unified format, and runs a configurable reconciliation engine that detects discrepancies, flags mismatches, and produces an auditable paper trail — automatically, on a schedule, at scale.

This is not a CRUD API. It is financial infrastructure.

---

## Table of Contents

- [Why This Exists](#why-this-exists)
- [How It Works](#how-it-works)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Core Concepts](#core-concepts)
  - [The Provider Interface](#the-provider-interface)
  - [Raw Storage](#raw-storage)
  - [Normalization](#normalization)
  - [The Reconciliation Engine](#the-reconciliation-engine)
  - [Reconciliation Statuses](#reconciliation-statuses)
- [Tech Stack](#tech-stack)
- [Database Schema](#database-schema)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Adding a New Provider](#adding-a-new-provider)
- [What's Still Missing](#whats-still-missing)
- [Contributing](#contributing)

---

## Why This Exists

If you've built payment systems, you've encountered this problem. Providers don't agree with each other. They don't always agree with themselves. Consider:

- **M-Pesa** returns `TransID` and `MSISDN` in a proprietary JSON envelope
- **Stripe** returns amounts in the lowest currency denomination (cents, not dollars)
- **Bank CSV exports** have inconsistent date formats, missing references, and occasional blank rows
- **PayPal** settlement reports lag by 24–48 hours

You can write one-off scripts to handle each provider. You can build a rats' nest of if-statements. Or you can build a proper abstraction layer that handles provider differences at the boundary, so your reconciliation logic never has to care what the source was.

That's what this is.

The secondary concern is auditability. In financial systems, you must be able to answer:

- Why was this transaction marked as reconciled?
- What was the raw payload from the provider at the time?
- Who flagged this for manual review and when?
- What was the variance?

Without these answers, reconciliation reports are just numbers. With them, they're a legal record.

---

## How It Works

At a high level, the flow is:

```
Merchant / Internal System
          ↓
  Provider Integration Layer
  (M-Pesa · Stripe · PayPal · Bank)
          ↓
   Polling & Sync Engine
   (Cloud Scheduler → Pub/Sub → sync-worker)
          ↓
  Raw Transaction Storage
  (original payload, untouched, JSONB)
          ↓
   Normalization Pipeline
   (normalize-worker maps to shared model)
          ↓
   Reconciliation Engine
   (match by reference · tolerance · fuzzy)
          ↓
   Mismatch Detection & Result Storage
          ↓
   Reports · Dashboard · Alerts
```

Each stage is handled by a dedicated worker, decoupled via Pub/Sub. This means:

- stages can fail independently and retry without cascading
- each worker can scale horizontally without affecting others
- you get a natural audit boundary between raw ingestion and processed output

---

## Architecture

The system is composed of five binaries, each with a single responsibility:

| Binary | Responsibility |
|---|---|
| `cmd/api` | REST API for reports, manual review, dashboard data |
| `cmd/scheduler` | Triggers sync jobs on a cron schedule per provider |
| `cmd/sync-worker` | Calls provider APIs, stores raw payloads, emits normalization events |
| `cmd/normalize-worker` | Converts raw payloads to `NormalizedTransaction`, emits reconciliation events |
| `cmd/reconcile-worker` | Matches internal vs external records, produces reconciliation results |

Workers communicate via GCP Pub/Sub. Each job is a JSON message with a provider name and time range. Workers are stateless and idempotent — safe to retry.

---

## Project Structure

```
smart-reconcile/
│
├── cmd/
│   ├── api/                    # HTTP server (Gin/Chi)
│   ├── scheduler/              # Cron-triggered job emitter
│   ├── sync-worker/            # Fetches from providers, stores raw
│   ├── normalize-worker/       # Maps raw → NormalizedTransaction
│   └── reconcile-worker/       # Runs the reconciliation engine
│
├── internal/
│   │
│   ├── providers/              # Payment provider integrations
│   │   ├── provider.go         # Provider interface definition
│   │   ├── mpesa/              # M-Pesa Daraja API integration
│   │   ├── stripe/             # Stripe API integration
│   │   ├── paypal/             # PayPal Reports API integration
│   │   └── bank/               # Bank CSV/SFTP integration
│   │
│   ├── sync/
│   │   ├── service.go          # Orchestrates fetch → store → publish
│   │   ├── fetcher.go          # Calls provider.FetchTransactions()
│   │   └── scheduler.go        # Interprets time range (daily/weekly/etc)
│   │
│   ├── normalization/
│   │   ├── service.go          # Dispatches to correct mapper
│   │   ├── mapper.go           # Mapper interface
│   │   └── currency.go         # Currency/denomination normalisation
│   │
│   ├── reconciler/
│   │   ├── engine.go           # Core reconciliation logic
│   │   ├── matcher.go          # Reference matching (exact + fuzzy)
│   │   ├── tolerance.go        # Variance threshold evaluation
│   │   └── discrepancy.go      # Discrepancy classification
│   │
│   ├── reporting/
│   │   ├── summary.go          # Reconciliation summary generation
│   │   ├── exports.go          # CSV/PDF export handlers
│   │   └── dashboard.go        # Aggregated metrics for UI
│   │
│   ├── alerts/
│   │   ├── email.go            # Email alerts for failed reconciliation
│   │   └── webhook.go          # Outbound webhook notifications
│   │
│   ├── queue/
│   │   └── pubsub.go           # Pub/Sub publisher and consumer wrappers
│   │
│   ├── storage/
│   │   ├── postgres/           # SQL queries and repository implementations
│   │   ├── redis/              # Caching and idempotency key storage
│   │   └── gcs/                # Raw payload archival to Cloud Storage
│   │
│   └── users/                  # Authentication, tenants, RBAC (stub)
│
├── pkg/
│   ├── models/                 # Shared domain types
│   └── dto/                    # API request/response structs
│
├── migrations/                 # SQL migration files (ordered)
├── terraform/                  # GCP infrastructure as code
└── scripts/                    # Local dev and CI helpers
```

---

## Core Concepts

### The Provider Interface

Every payment provider in this system implements one interface:

```go
package providers

import "context"

type Transaction struct {
    ExternalID string
    Amount     decimal.Decimal
    Currency   string
    Reference  string
    Timestamp  time.Time
    Provider   string
}

type Provider interface {
    Name() string
    FetchTransactions(ctx context.Context, from, to time.Time) ([]Transaction, error)
}
```

That's it. The sync engine doesn't know or care whether it's talking to M-Pesa or a bank SFTP server. It calls `FetchTransactions`, gets back a slice of `Transaction`, and hands it to the storage layer.

**Why this matters:** Without this abstraction, adding a new provider means touching the sync engine, the normalization layer, the reconciler, and probably a config file or two. With it, adding a provider means creating one folder under `internal/providers/` and implementing two methods. The rest of the system doesn't change.

The alternative — a proliferation of SDKs, type switches, and provider-specific conditionals scattered throughout the codebase — is how fintech projects become unmaintainable.

---

### Raw Storage

Before any normalization or processing occurs, every transaction payload is stored exactly as received from the provider:

```sql
CREATE TABLE raw_transactions (
    id         UUID PRIMARY KEY,
    provider   TEXT NOT NULL,
    payload    JSONB NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

The `payload` column is `JSONB` — unmodified, unparsed, exactly what the API returned. Bank CSVs are stored as JSON-encoded rows. Binary formats are stored as base64.

**This is non-negotiable in financial systems.** Here's why:

- **Audits.** Regulators and auditors want to see the original record, not what you did with it.
- **Disputes.** When a merchant disputes a transaction, you need the original provider payload to prove what was received and when.
- **Debugging.** Normalization bugs are subtle. When your mapper produces wrong output, you need the original input to reproduce and fix the issue.
- **Provider inconsistencies.** Providers quietly change their API response shapes. When that happens and your normalization breaks, raw storage lets you reprocess without re-fetching.
- **Reprocessing.** If the normalization or reconciliation logic changes, you can replay historical raw records through the new pipeline without going back to the provider.

Financial systems without raw audit storage don't stay financial systems for long.

---

### Normalization

Different providers return wildly different data shapes:

**M-Pesa:**
```json
{
  "TransID": "QWE123",
  "TransAmount": "1500",
  "MSISDN": "2547XXXXXXXX",
  "TransTime": "20260512143022"
}
```

**Stripe:**
```json
{
  "id": "txn_1abc",
  "amount": 150000,
  "currency": "usd",
  "created": 1715521822
}
```

**Bank CSV:**
```
reference,amount,date
INV-001,5000.00,2026-05-12
```

The normalization pipeline converts all of these into a single unified model:

```go
type NormalizedTransaction struct {
    ID         string
    Provider   string
    Amount     decimal.Decimal
    Currency   string
    Reference  string
    OccurredAt time.Time
}
```

Note the use of `decimal.Decimal` (via `github.com/shopspring/decimal`) rather than `float64`. This is intentional and important — see the section on precision below.

Each provider has a dedicated mapper implementing:

```go
type Mapper interface {
    Normalize(payload []byte) (*NormalizedTransaction, error)
}
```

Mappers handle:
- field name differences
- timestamp format parsing (Unix epoch, RFC3339, M-Pesa's `YYYYMMDDHHmmss`, etc.)
- amount denomination conversion (Stripe's cents → decimal currency units)
- currency code normalization (uppercase ISO 4217)
- missing or optional field defaults

Once a transaction is normalized, the reconciliation engine operates on a consistent model regardless of source.

---

### The Reconciliation Engine

The engine takes two `NormalizedTransaction` values — one internal record, one external provider record — and produces a `Result`:

```go
type Result struct {
    ID           uuid.UUID
    InternalTxID uuid.UUID
    ExternalTxID uuid.UUID
    Status       string
    Variance     decimal.Decimal
    Reason       string
    ReconciledAt time.Time
}
```

The matching process runs in order:

1. **Currency check** — mismatched currencies are an immediate failure. No amount comparison is meaningful across currency boundaries without FX rates.

2. **Reference matching** — references are normalized (trimmed, uppercased, special characters removed) and compared exactly first. If that fails, Levenshtein distance is used as a fallback to catch typos and minor formatting differences.

   ```
   INV-001   →  INV001
   INV_001   →  INV001   ← same after normalization
   PAYMNT123 →  fuzzy match against PAYMENT123
   ```

3. **Tolerance check** — if references match, the variance between amounts is calculated:

   ```go
   variance := internal.Amount.Sub(external.Amount).Abs()
   ```

   If the variance is within the configured tolerance (e.g. ±1.00 KES), the transaction is marked `MATCHED`. If it exceeds tolerance, it goes to `MANUAL_REVIEW`.

The engine is designed around a few important properties:

- **Decimal arithmetic throughout.** `0.1 + 0.2 = 0.30000000000000004` in float64. In a reconciliation system that processes millions of transactions, that drift compounds into real money. Decimal arithmetic is exact.
- **Absolute variance.** The engine computes `Abs()` before tolerance comparison so negative variances (internal < external) and positive variances (internal > external) are treated consistently.
- **Immutable results.** Reconciliation results are insert-only. If a result needs to change, a new result is created. The history is preserved.
- **Idempotency.** Transactions are uniquely indexed by `(provider, external_id)`. Duplicate webhooks and retried sync jobs produce no duplicate records.

---

### Reconciliation Statuses

| Status | Meaning |
|---|---|
| `PENDING` | Transaction ingested, not yet reconciled |
| `MATCHED` | Internal and external records agree within tolerance |
| `PARTIAL_MATCH` | Reference matched, amount variance exceeds tolerance but is close |
| `MANUAL_REVIEW` | Significant discrepancy or reference mismatch — human review required |
| `FAILED` | Reconciliation process error (logged, dead-lettered, retriable) |
| `SETTLED` | Manually confirmed and closed |

Every status transition is recorded with a timestamp and reason. The full transition history is queryable via the API.

---

## Tech Stack

| Layer | Choice | Why |
|---|---|---|
| Language | Go 1.22 | Strong concurrency model, low overhead for worker processes, excellent standard library for financial workloads |
| Decimal arithmetic | `shopspring/decimal` | Exact monetary arithmetic, no float drift |
| Fuzzy matching | `agnivade/levenshtein` | Lightweight, no ML dependency, tunable threshold |
| Queue | GCP Pub/Sub | At-least-once delivery, dead-letter queues, replay support |
| Scheduler | GCP Cloud Scheduler | Managed cron with Pub/Sub integration |
| Database | Cloud SQL (PostgreSQL) | ACID transactions, JSONB for raw payloads, mature ecosystem |
| Workers | Cloud Run Jobs | Stateless, horizontally scalable, pay-per-execution |
| API | Cloud Run | Auto-scaling, containerised, zero idle cost |
| Object storage | Cloud Storage | Raw payload archival, CSV exports |
| Secrets | Secret Manager | No credentials in environment variables or config files |
| Observability | Cloud Monitoring | Metrics, alerting, uptime checks |
| IaC | Terraform | Reproducible GCP infrastructure |

---

## Database Schema

### `transactions`
```sql
CREATE TABLE transactions (
    id          UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    provider    TEXT NOT NULL,
    external_id TEXT NOT NULL,
    reference   TEXT NOT NULL,
    amount      NUMERIC(18, 2) NOT NULL,
    currency    TEXT NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    status      TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Prevents duplicate ingestion of the same provider transaction
CREATE UNIQUE INDEX idx_provider_external_id
    ON transactions (provider, external_id);
```

### `raw_transactions`
```sql
CREATE TABLE raw_transactions (
    id         UUID PRIMARY KEY,
    provider   TEXT NOT NULL,
    payload    JSONB NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### `reconciliation_results`
```sql
CREATE TABLE reconciliation_results (
    id            UUID PRIMARY KEY,
    internal_tx_id UUID NOT NULL REFERENCES transactions(id),
    external_tx_id UUID NOT NULL REFERENCES transactions(id),
    status        TEXT NOT NULL,
    variance      NUMERIC(18, 2),
    reason        TEXT,
    reconciled_at TIMESTAMP NOT NULL
);
```

Migrations live in `migrations/` and are numbered sequentially. Run them in order.

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (for local Postgres and Pub/Sub emulator)
- GCP project with billing enabled (for cloud deployment)
- Terraform (for infrastructure provisioning)

### Local development

```bash
# Clone the repo
git clone https://github.com/your-org/smart-reconcile.git
cd smart-reconcile

# Start local dependencies
docker compose up -d   # starts Postgres + Pub/Sub emulator

# Run migrations
psql $DATABASE_URL -f migrations/001_create_transactions.sql
psql $DATABASE_URL -f migrations/002_create_reconciliation_results.sql
psql $DATABASE_URL -f migrations/003_create_raw_transactions.sql

# Install dependencies
go mod download

# Run the API
go run ./cmd/api

# In separate terminals, run the workers
go run ./cmd/scheduler-worker
go run ./cmd/normalize-worker
go run ./cmd/reconcile-worker
```

### Cloud deployment

```bash
# Provision GCP infrastructure
cd terraform/
terraform init
terraform plan
terraform apply

# Build and push containers
gcloud builds submit --config=cloudbuild.yaml
```

---

## Configuration

Configuration is read from environment variables. Secrets are pulled from GCP Secret Manager at startup.

| Variable | Description | Example |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host/db` |
| `PUBSUB_PROJECT_ID` | GCP project ID | `my-project-123` |
| `PUBSUB_SYNC_TOPIC` | Topic for sync jobs | `sync-jobs` |
| `PUBSUB_NORMALIZE_TOPIC` | Topic for normalization jobs | `normalize-jobs` |
| `PUBSUB_RECONCILE_TOPIC` | Topic for reconciliation jobs | `reconcile-jobs` |
| `RECONCILE_TOLERANCE` | Max allowed monetary variance | `1.00` |
| `FUZZY_MATCH_THRESHOLD` | Max Levenshtein distance for reference fuzzy match | `2` |
| `MPESA_CONSUMER_KEY` | M-Pesa Daraja API consumer key | (from Secret Manager) |
| `MPESA_CONSUMER_SECRET` | M-Pesa Daraja API consumer secret | (from Secret Manager) |
| `STRIPE_SECRET_KEY` | Stripe secret key | (from Secret Manager) |
| `PAYPAL_CLIENT_ID` | PayPal API client ID | (from Secret Manager) |
| `PAYPAL_CLIENT_SECRET` | PayPal API client secret | (from Secret Manager) |
| `ALERT_EMAIL_FROM` | Sender address for alert emails | `alerts@yourcompany.com` |
| `ALERT_WEBHOOK_URL` | Webhook URL for mismatch notifications | `https://hooks.slack.com/...` |

---

## Adding a New Provider

1. Create a folder under `internal/providers/yourprovider/`

2. Implement the `Provider` interface:

```go
package yourprovider

type Client struct {
    apiKey string
}

func New(apiKey string) *Client {
    return &Client{apiKey: apiKey}
}

func (c *Client) Name() string {
    return "yourprovider"
}

func (c *Client) FetchTransactions(
    ctx context.Context,
    from, to time.Time,
) ([]providers.Transaction, error) {
    // call your provider API
    // map response fields to providers.Transaction
    // return the slice
}
```

3. Implement a normalization mapper under `internal/normalization/`:

```go
type YourProviderMapper struct{}

func (m YourProviderMapper) Normalize(payload []byte) (*NormalizedTransaction, error) {
    // parse the raw payload
    // return NormalizedTransaction
}
```

4. Register both in the provider registry and normalization dispatcher.

5. Add the provider's credentials to Secret Manager and the relevant environment variable documentation.

That's it. No changes to the sync engine, reconciliation engine, or reporting layer.

---

## What's Still Missing

This is a foundation, not a finished product. The following are known gaps that a production deployment would need to address:

**Observability**
- OpenTelemetry distributed tracing across workers
- Prometheus metrics (reconciliation rate, match rate, error rate per provider)
- Structured logging with correlation IDs

**Reliability**
- Dead-letter queues for failed Pub/Sub messages
- Exponential backoff with jitter on provider API calls
- Database-level row locking for concurrent reconciliation jobs
- Explicit database transaction wrapping for multi-table writes

**Security**
- RBAC for the API (multi-tenant access control)
- Encryption at rest for `raw_transactions.payload`
- Secret rotation automation via Secret Manager

**Provider completeness**
- Full M-Pesa Daraja B2C/B2B reconciliation
- PayPal settlement report ingestion
- Generic bank CSV/SFTP adapter with configurable column mapping

**Testing**
- Integration tests against provider sandbox APIs
- Contract tests for the Provider interface
- Chaos testing (simulated provider outages, duplicate callbacks, malformed payloads)

**Scale**
- Multi-region failover for Cloud SQL
- Horizontal sharding strategy for high transaction volumes
- Multi-currency FX normalisation for cross-currency reconciliation

Contributions in any of these areas are welcome.

---

## Contributing

Pull requests are welcome. For significant changes, open an issue first to discuss what you'd like to change.

When contributing a new provider integration, please include:
- The provider client implementation
- The normalization mapper
- A test against the provider's sandbox API (or fixture-based unit tests if no sandbox exists)
- Documentation in this README

For reconciliation engine changes, be especially careful with:
- Decimal arithmetic (never introduce `float64` for monetary values)
- Status transition logic (reconciliation results are append-only)
- Tolerance configuration (changes affect how historical results would have been classified)

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

*Built because someone had to. Financial infrastructure is not glamorous work. It is, however, necessary work...and the difference between doing it well and doing it badly is the difference between a reconciliation report and a crisis.*