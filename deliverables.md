# Smart Reconcile SaaS - Deliverables & System Architecture

## Product Overview

Smart Reconcile is a multi-tenant B2B SaaS platform that automates transaction reconciliation across multiple payment providers and financial systems.

The platform enables merchants to:

* Connect payment providers
* Automatically fetch transaction data
* Normalize provider-specific transaction formats
* Store and manage transaction records
* Reconcile transactions
* Detect discrepancies
* Generate reports
* Export statements
* Receive alerts and notifications

---

# High-Level Business Flow

```text
Connect Provider
       │
       ▼
Fetch Transactions
       │
       ▼
Normalize Data
       │
       ▼
Store Transactions
       │
       ▼
Run Reconciliation
       │
       ▼
Generate Reports
       │
       ▼
Export Statements
       │
       ▼
Send Alerts
```

---

# Target Architecture

```text
smart-reconcile/

├── cmd/
│   └── api/
│
├── internal/
│
│   ├── auth/
│   ├── users/
│   ├── merchants/
│   ├── providers/
│   ├── sync/
│   ├── transactions/
│   ├── normalizer/
│   ├── reconciler/
│   ├── reporting/
│   ├── exports/
│   ├── alerts/
│   ├── audit/
│   ├── observability/
│   ├── config/
│   └── storage/
│
├── api/
├── deployments/
├── scripts/
└── docs/
```

---

# Domain Deliverables

## 1. Auth Domain

### Purpose

Authentication and authorization.

### Deliverables

* User login
* JWT generation
* Password hashing
* Password reset
* Session validation
* Role enforcement

### Roles

* Platform Admin
* Merchant Admin
* Auditor
* Read-Only User

### Future Enhancements

* SSO
* OAuth
* MFA

---

## 2. Users Domain

### Purpose

User lifecycle management.

### Deliverables

* User registration
* User updates
* User deactivation
* User invitations
* Role assignment

### Audit Events

Track:

* Login activity
* Password changes
* Permission changes

---

## 3. Merchants Domain

### Purpose

Tenant management.

Every customer of the SaaS platform is a merchant.

### Deliverables

### Merchant Profiles

Store:

* Business Name
* Contact Details
* Status

### Merchant Settings

Configure:

* Sync schedules
* Reconciliation rules
* Alert preferences
* Export preferences

### Provider Connections

Store:

* API credentials
* OAuth tokens
* Webhook secrets

---

## 4. Providers Domain

### Purpose

External system integrations.

### Deliverables

Provider adapters for:

* M-Pesa
* Stripe
* Flutterwave
* PayPal
* Bank APIs

### Responsibilities

* Authentication
* API communication
* Pagination
* Rate limiting
* Retry handling

### Provider Interface

Each provider should support:

```text
Validate Connection
Fetch Transactions
Fetch Settlements
Fetch Balances
```

---

## 5. Sync Domain

### Purpose

Transaction acquisition.

### Deliverables

### Manual Sync

User-triggered synchronization.

```text
Sync Now
```

### Scheduled Sync

Merchant-configurable frequencies:

```text
Hourly
Daily
Weekly
Monthly
```

### Incremental Sync

Store:

```text
last_synced_at
```

Fetch only newly available transactions.

### Sync Tracking

Track:

* Started
* Completed
* Failed
* Duration

---

## Scheduler Architecture

### MVP

Use:

```text
Google Cloud Scheduler
```

Flow:

```text
Cloud Scheduler
      │
      ▼
Cloud Run Endpoint
      │
      ▼
Sync Service
```

### Future

Add:

```text
Pub/Sub
Cloud Tasks
```

for larger workloads.

---

## 6. Transactions Domain

### Purpose

Central transaction management.

### Deliverables

### Transaction Storage

Store:

* Merchant
* Provider
* Reference
* Amount
* Currency
* Timestamp
* Status

### Validation

Check:

* Required fields
* Invalid values
* Duplicate references

### Deduplication

Prevent duplicate transaction ingestion.

### Idempotency

Multiple sync runs should never create duplicate records.

---

## 7. Normalizer Domain

### Purpose

Transform provider payloads into a common format.

### Deliverables

### Mapping Engine

Convert:

```text
Provider Payload
        ↓
Internal Transaction Model
```

### Validation Rules

Verify:

* Reference format
* Amount validity
* Required attributes

### Extensibility

New providers should require only:

```text
Provider Adapter
+
Provider Mapper
```

---

## 8. Reconciler Domain

### Purpose

Core reconciliation engine.

### Deliverables

### Matching Engine

Support:

* Exact Match
* Reference Match
* Amount Match
* Date Match

### Tolerance Rules

Examples:

```text
± Amount Variance

± Time Variance
```

### Merchant Rules

Merchant-specific reconciliation behavior.

### Reconciliation Runs

Track:

* Pending
* Running
* Completed
* Failed

### Reconciliation Results

Statuses:

```text
Matched
Unmatched
Partial Match
Exception
```

### Discrepancy Detection

Detect:

* Missing transactions
* Duplicate transactions
* Amount mismatches
* Status mismatches

---

## 9. Reporting Domain

### Purpose

Business analytics.

### Deliverables

### Dashboards

Metrics:

* Reconciliation Rate
* Transaction Volume
* Failed Transactions
* Exceptions

### Historical Reporting

Views:

* Daily
* Weekly
* Monthly
* Yearly

### Merchant Analytics

Track:

* Settlement Delays
* Provider Performance
* Reconciliation Trends

---

## 10. Exports Domain

### Purpose

Generate downloadable files.

### Deliverables

### CSV Exports

Generate:

```text
.csv
```

### Excel Exports

Generate:

```text
.xlsx
```

### Future

PDF statements.

---

## Export Flow

```text
Generate Report
       │
       ▼
Cloud Storage
       │
       ▼
Signed URL
       │
       ▼
Download
```

### Storage

Use:

Google Cloud Storage

for all exported files.

---

## 11. Alerts Domain

### Purpose

Notification management.

### Deliverables

### Email Alerts

Notify:

* Sync failures
* Reconciliation failures
* System issues

### Webhook Alerts

Integrations:

* Slack
* Teams
* Custom endpoints

### Alert Severity

Levels:

```text
INFO
WARNING
CRITICAL
```

---

## 12. Audit Domain

### Purpose

Compliance and traceability.

### Deliverables

Track:

* User logins
* Report downloads
* Sync executions
* Reconciliation runs
* Configuration changes

### Audit Requirements

Every critical action must be attributable to:

```text
Who
What
When
```

---

## 13. Observability Domain

### Purpose

Operational monitoring.

### Deliverables

### Structured Logging

Track:

* API requests
* Sync jobs
* Reconciliation jobs

### Metrics

Track:

* Sync duration
* Reconciliation duration
* Match rate
* Failure rate

### Tracing

Follow:

```text
Sync
 ↓
Normalize
 ↓
Store
 ↓
Reconcile
```

### GCP Services

Use:

* Cloud Logging
* Cloud Monitoring

---

## 14. Config Domain

### Purpose

Centralized configuration management.

### Deliverables

Manage:

* Database configuration
* Provider configuration
* Alert configuration
* Scheduler configuration

### Secrets

Store in:

```text
Google Secret Manager
```

Never store credentials in source control.

---

# Infrastructure Deliverables

## Compute

Use:

```text
Cloud Run
```

Benefits:

* Serverless
* Auto-scaling
* Pay-per-use

---

## Database

Use:

```text
Cloud SQL PostgreSQL
```

Store:

* Transactions
* Merchants
* Users
* Reconciliation Results
* Audit Logs

---

## Scheduling

Use:

```text
Cloud Scheduler
```

for all automated sync jobs.

---

## File Storage

Use:

```text
Cloud Storage
```

for:

* Reports
* Exports
* Archived statements

---

## Secrets

Use:

```text
Secret Manager
```

for:

* API keys
* Provider credentials
* JWT secrets

---

# Future Scale Architecture

When transaction volume becomes significant:

```text
Cloud Scheduler
        │
        ▼
Cloud Run
        │
        ▼
Pub/Sub
        │
 ┌──────┼──────┐
 ▼      ▼      ▼

Sync  Normalize Reconcile
Workers Workers Workers
```

This allows independent scaling of ingestion, normalization, and reconciliation workloads.

---

# MVP Scope

## Build First

* Auth
* Users
* Merchants
* Provider Integrations
* Sync
* Transactions
* Normalizer
* Reconciler
* Reporting
* Exports
* Alerts
* Audit
* Cloud Run
* Cloud SQL
* Cloud Scheduler
* Cloud Storage

## Build Later

* Pub/Sub
* Cloud Tasks
* Rule Builder
* Bank Integrations
* AI-assisted Matching
* Real-Time Reconciliation
* Multi-Region Deployments

The primary objective of the MVP is not sophisticated reconciliation logic. The primary objective is proving that businesses will trust the platform with their transaction operations and pay to eliminate manual reconciliation work.
