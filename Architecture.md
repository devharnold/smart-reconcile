# Smart Reconciliation Platform (Go)

## Recommended Stack

-   Go + chi/gin
-   PostgreSQL
-   Apache Kafka
-   Redis
-   pgx or sqlc
-   robfig/cron
-   WebSockets
-   zap/zerolog
-   Prometheus + Grafana

## Core Services

-   API Service
-   Scheduler
-   Fetch Worker
-   Normalize Worker
-   Reconciliation Worker
-   Notification Worker
-   Workflow Orchestrator

## Event Flow

Scheduler -\> JobStarted -\> Orchestrator -\> StartFetch Fetch Worker
-\> TransactionsFetched Orchestrator -\> StartNormalization or
StartReconciliation Normalize Worker -\> TransactionsNormalized
Reconcile Worker -\> ReconciliationCompleted Notification Worker listens
for completion/failure events and: - Stores notification - Pushes via
WebSocket - Sends email if enabled

## Job Status

QUEUED -\> FETCHING -\> NORMALIZING -\> RECONCILING -\> COMPLETED or
FAILED

## Kafka Topics

jobs transactions normalization reconciliation notifications
