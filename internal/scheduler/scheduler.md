# Go Scheduler for Automated Transaction Fetching

## Overview

This document explains the scheduler implementation used for the automated transaction fetching service in the reconciliation platform.

The scheduler is responsible for triggering the reconciliation process at configurable intervals without requiring manual intervention.

---

# Objective

The scheduler enables the application to automatically fetch financial transactions from external institutions according to a predefined schedule.

Supported schedules include:

* Daily
* Weekly
* Monthly

Each schedule can execute at a configurable time.

For example:

* Every day at **06:00 AM**
* Every Sunday at **06:00 AM**
* On the first day of every month at **06:00 AM**

---

# Architecture

```
                 Application Starts
                        │
                        ▼
              Load Scheduler Configuration
                        │
                        ▼
             Generate Cron Expression
                        │
                        ▼
          Register Job with Cron Scheduler
                        │
                        ▼
                Wait for Scheduled Time
                        │
                        ▼
           Execute Reconciliation Service
                        │
                        ▼
            Fetch Financial Transactions
                        │
                        ▼
              Validate Transactions
                        │
                        ▼
              Save New Transactions
                        │
                        ▼
             Begin Reconciliation Process
```

---

# Project Structure

```
internal/
│
├── scheduler/
│   ├── config.go
│   ├── cron.go
│   └── scheduler.go
│
├── service/
│   └── reconciliation.go
│
cmd/
└── server/
    └── main.go
```

---

# Components

## config.go

This file defines the scheduling configuration.

It contains:

* Available scheduling frequencies
* Schedule structure
* Time configuration

Example:

```go
type Frequency string

const (
    Daily  Frequency = "daily"
    Weekly Frequency = "weekly"
    Monthly Frequency = "monthly"
)

type Schedule struct {
    Frequency Frequency

    Hour int
    Minute int

    Weekday int
    Day int
}
```

---

## cron.go

This component converts a schedule into a cron expression understood by the scheduler.

Example conversions:

| Schedule           | Cron Expression |
| ------------------ | --------------- |
| Daily at 6:00 AM   | `0 6 * * *`     |
| Weekly at 6:00 AM  | `0 6 * * 0`     |
| Monthly at 6:00 AM | `0 6 1 * *`     |

The implementation uses Go's `fmt.Sprintf()` to dynamically generate the correct cron expression.

---

## scheduler.go

This component manages the cron scheduler.

Responsibilities include:

* Creating the scheduler
* Registering jobs
* Starting the scheduler
* Gracefully stopping the scheduler

Example workflow:

```
Create Scheduler
        │
        ▼
Register Job
        │
        ▼
Start Scheduler
        │
        ▼
Wait for Scheduled Time
        │
        ▼
Execute Job
```

---

## reconciliation.go

This service contains the business logic executed by the scheduler.

Typical responsibilities include:

1. Fetch transactions from financial institutions.
2. Validate incoming transactions.
3. Remove duplicate records.
4. Save transactions.
5. Trigger reconciliation logic.
6. Generate logs and notifications if required.

The scheduler should only trigger this service and should not contain business logic itself.

---

## main.go

The application's entry point.

Responsibilities:

1. Create the reconciliation service.
2. Create the scheduler.
3. Register the scheduled job.
4. Start the scheduler.
5. Keep the application running.

Example startup flow:

```
Application Starts
        │
        ▼
Create Services
        │
        ▼
Register Scheduler
        │
        ▼
Start Scheduler
        │
        ▼
Scheduler Waits
```

---

# Supported Scheduling Options

## Daily

Runs once every day.

Example:

```
06:00 AM
```

Cron:

```
0 6 * * *
```

---

## Weekly

Runs once every week on a specified weekday.

Example:

```
Every Sunday at 06:00 AM
```

Cron:

```
0 6 * * 0
```

---

## Monthly

Runs once every month on a specified day.

Example:

```
First day of the month at 06:00 AM
```

Cron:

```
0 6 1 * *
```

---

# Execution Flow

```
Scheduler Triggered
        │
        ▼
Generate Cron Schedule
        │
        ▼
Cron Executes Job
        │
        ▼
Reconciliation Service Starts
        │
        ▼
Fetch Transactions
        │
        ▼
Validate Transactions
        │
        ▼
Store Transactions
        │
        ▼
Update Last Fetch Timestamp
        │
        ▼
Finish Execution
```

---

# Advantages of the Design

* Separation of scheduling and business logic.
* Easy to add new scheduling frequencies.
* Simple to maintain and extend.
* Supports configurable execution times.
* Suitable for production environments.
* Scales well as more financial institutions are added.

---

# Future Improvements

The scheduler can be enhanced by:

* Loading schedules dynamically from a database.
* Allowing administrators to update schedules without restarting the application.
* Supporting multiple concurrent scheduled jobs.
* Adding retry logic for failed transaction fetches.
* Integrating distributed locks to prevent duplicate execution in multi-instance deployments.
* Publishing events to a message broker for asynchronous reconciliation.
* Recording execution history, duration, and failures for auditing and monitoring.

---

# Conclusion

The scheduler provides a reliable mechanism for automating transaction retrieval in the reconciliation platform. By separating scheduling concerns from business logic, the system remains modular, maintainable, and scalable. This design supports configurable daily, weekly, and monthly execution schedules while providing a solid foundation for future enhancements such as distributed scheduling, monitoring, and event-driven processing.
