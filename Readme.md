# Distributed Job Scheduler

A scalable, fault-tolerant distributed job scheduling system built with
Go.

This project demonstrates leader election, worker leasing, job execution
tracking, TTL expiry handling, and horizontal scalability using Redis
and PostgreSQL.

------------------------------------------------------------------------

## 🚀 Overview

Distributed Job Scheduler is designed to coordinate background jobs
across multiple worker nodes while ensuring:

-   Exactly-once job leasing
-   Leader-based orchestration
-   Fault tolerance via TTL expiry
-   Horizontal worker scaling
-   Clean separation of concerns

It is production-inspired and Kubernetes-ready.

------------------------------------------------------------------------

## 🏗 Architecture

The system consists of three primary services:

  -----------------------------------------------------------------------
  Service                     Responsibility
  --------------------------- -------------------------------------------
  Scheduler Service           Core orchestration, leader election, job
                              lifecycle management

  Worker Service              Job execution and heartbeat reporting

  Proxy Service               Routes requests to the active scheduler
                              leader
  -----------------------------------------------------------------------

Shared contracts and constants live under `common/`.

------------------------------------------------------------------------

## 📦 Project Structure

    distributed-job-scheduler/
    │
    ├── common/                # Shared types and contracts
    ├── scheduler-service/     # Core scheduling engine
    ├── worker-service/        # Distributed job workers
    ├── proxy-service/         # Leader-aware request router
    ├── go.work                # Go workspace definition
    └── docker-compose.yml     # Local dependency setup

------------------------------------------------------------------------

## ⚙️ Tech Stack

-   Go
-   PostgreSQL
-   Redis
-   Docker
-   Swagger
-   Kubernetes-ready design

------------------------------------------------------------------------

## 🔄 System Flow

1.  Client submits a job
2.  Proxy routes request to active scheduler leader
3.  Scheduler persists job in PostgreSQL
4.  Worker leases job via Redis
5.  Worker executes and updates status
6.  Scheduler tracks execution and cleans expired workers via TTL
    listeners

------------------------------------------------------------------------

## 🧠 Core Concepts

### Leader Election

Only one scheduler instance actively assigns jobs at a time.

### Worker Leasing

Workers lease jobs via Redis to prevent duplicate execution.

### Event-Driven Coordination

Redis Pub/Sub is used for distributed signaling.

------------------------------------------------------------------------

## 🛠 Running Locally

### 1️⃣ Start Dependencies

    docker-compose up

### 2️⃣ Start Scheduler

    cd scheduler-service
    go run cmd/server/main.go

### 3️⃣ Start Worker

    cd worker-service
    go run cmd/main.go

### 4️⃣ Start Proxy

    cd proxy-service
    go run main.go

------------------------------------------------------------------------

## 📘 API Documentation

Swagger documentation is available inside:

    scheduler-service/docs

------------------------------------------------------------------------

## 🧪 Scalability & Fault Tolerance

-   Multiple workers can run concurrently
-   Multiple scheduler instances can run (leader election ensures single
    active leader)
-   Worker crashes are handled via TTL expiry
-   Safe job leasing prevents duplicate execution

.
