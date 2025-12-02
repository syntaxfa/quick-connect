# Project Structure

This document outlines the confirmed structure of the **Quick Connect** project, designed for clarity, scalability, and high maintainability, based on the provided codebase and the architectural guidelines.

-----

## Directory Structure

```
.
├── adapter           # External Interfaces and Clients (Redis, gRPC Clients, DB, Observability)
├── app               # Core Microservices (Admin, Chat, Manager, Notification, File Handler)
├── cli               # Client Command-Line Interface tools
├── cmd               # Entry Points (main.go) and service execution commands
├── config            # Shared configuration loading logic and types
├── deploy            # Deployment files (Dockerfiles, Docker Compose)
├── docs              # Project documentation, HLD, and developer guides
├── event             # Event Contracts (Reserved for future asynchronous communication)
├── example           # Runnable examples for core technologies
├── frontend          # Separate User Interface projects (Admin Panels and Website)
├── logs              # Storage location for service logs
├── outbox            # Implementation of the Transactional Outbox pattern
├── pkg               # Shared and reusable packages/utilities
├── protobuf          # Protocol Buffer files and generated gRPC code
├── sdk               # Client-side Software Development Kits (SDKs)
└── types             # Shared foundational data types
```

-----

## Directory Descriptions

### 1\. `adapter` (External Interfaces & Clients)

**Purpose:** Contains clients and interfaces that manage communication with **external infrastructure** (e.g., databases, caching systems) and **other microservices** via gRPC. This layer isolates network and client concerns from the core business logic, serving as the **Driven/Secondary Adapters** in a Hexagonal architecture.
**Observed Structure:**

* `adapter/manager`, `adapter/chat`: gRPC clients for inter-service communication with the Manager and Chat services. These directories also contain corresponding `*_local.go` files that implement the same interface for use in Monolith mode.
* `adapter/postgres`, `adapter/redis`: Clients for PostgreSQL and Redis (caching).
* `adapter/pubsub`: Pub/Sub implementation (currently Redis Pub/Sub).
* `adapter/observability`: OpenTelemetry instrumentation setup for tracing and metrics.

-----

### 2\. `app` (Core Microservices)

**Purpose:** Contains the core code for each microservice. Each subdirectory represents an independent, deployable service, strictly organized into layers following the **Hexagonal Architecture** pattern:

* `delivery`: Protocol implementation layer (**Primary/Driving Adapters**) handling inbound requests (HTTP, gRPC, WebSocket).
* `service`: Core business/domain logic and defines the **Ports (Interfaces)**.
* `repository`: Data access layer (**Secondary/Driven Adapter**) for the local database, including database `migrations`.

-----

### 3\. `cmd` (Service Entry Points)

**Purpose:** Holds the main entry points (`main.go`) for running each microservice server and utility CLI commands.
**Observed Structure:**

* `cmd/{service}/main.go`: The main server executor.
* `cmd/{service}/command`: Contains auxiliary commands like `migrate` and `create_user`.
* `cmd/all-in-one`: A consolidated entry point for quickly launching all core services in a single binary/container for development or small-scale deployment.

-----

### 4\. `deploy` (Deployment)

**Purpose:** Contains all files related to building, packaging, and deploying the services. This includes infrastructure-as-code definitions.
**Key Components:** **Dockerfiles**, **Docker Compose** files (`compose.yml` for various environments/setups, including microservices orchestration), and environment-specific configuration files (`config.yml`).

-----

### 5\. `frontend` (Separate User Interface Projects)

**Purpose:** Contains the source code for larger, dedicated User Interface projects and web clients that are developed separately from the main Go backend.
**Observed Structure:** Includes dedicated admin panels (`chat-admin-panel`, `notification-admin-panel`) and the main website project (`quick-connect-website`).

-----

### 6\. `sdk` (Software Development Kits)

**Purpose:** Provides client-side libraries for developers to integrate Quick Connect services (e.g., chat widget, file uploads) into their own applications.
**Observed Structure:** Includes the Chat SDK (for Web/JS) and placeholder directories for `file` and `notification` SDKs.

-----

### 7\. `pkg` (Shared Packages)

**Purpose:** Low-level, general-purpose packages and utilities reused across multiple services. This code is strictly agnostic to the specific business logic of any single application.
**Key Examples:** Networking (`httpserver`, `grpcserver`), Error Handling (`richerror`), Security (`jwtvalidator`, `ratelimit`), and utility functions (`logger`, `paginate`, `translation`).

-----

### 8\. `protobuf` (Protocol Buffers)

**Purpose:** Stores the Protocol Buffer definition files (`.proto`) and the generated Go bindings used for gRPC communication.
**Observed Structure:** Files are organized by service (`chat`, `manager`) and contain the source (`proto`) and generated Go code (`golang`).

-----

### 9\. `outbox` (Transactional Outbox Pattern)

**Purpose:** Implements the core components of the **Transactional Outbox Pattern** to ensure reliable, asynchronous event publishing. This mechanism relies on saving events within the database transaction before publishing them.

-----

## 10\. Core Architectural Design: Hexagonal & Dual Deployment

The project strictly adheres to **Hexagonal Architecture (Ports & Adapters)** principles to ensure high decoupling and testability. This architecture is key to enabling the unique "Code-Level Monolith" or "Dual Deployment" capability:

### A. Hexagonal Architecture Implementation

1.  **Ports (Interfaces):** Defined within the `app/{service}/service` packages. These are the contracts that the business logic requires, independent of who calls them or how data is stored.
2.  **Adapters:** Implement the Ports.
    * **Primary/Driving Adapters:** Located in `app/{service}/delivery` (handling incoming traffic like HTTP, gRPC).
    * **Secondary/Driven Adapters:** Located in `app/{service}/repository` (Database access) and `adapter/` (External service clients).

### B. The Dual Deployment Mechanism

The project uses the same codebase to run in two modes by switching how the Secondary Adapters connect to other services:

| Deployment Mode | Inter-Service Adapter Used | Connection Type |
| :--- | :--- | :--- |
| **Microservices Mode** | Dedicated gRPC Clients (e.g., `adapter/manager/auth.go`) | Network Communication (Remote) |
| **Monolith (All-in-One) Mode** | Local Adapters (e.g., `adapter/manager/auth_local.go`) | Direct Function Call (In-Memory) |

This structural choice ensures the Core Business Logic remains oblivious to the deployment mode, making the system inherently flexible and scalable.

-----

## 11\. Other Directories

| Directory | Description |
| :--- | :--- |
| **`config`** | Centralized logic for reading, loading, and validating configuration settings. |
| **`docs`** | Stores technical documentation, High-Level Design (HLD) diagrams, and developer guides. |
| **`event`** | **(Reserved)** The designated location for defining versioned event contracts (Go structs) for future asynchronous communication. Currently empty, awaiting implementation. |
| **`example`** | Contains runnable, isolated examples demonstrating core technology usage (e.g., gRPC client calls, Observability setup, Pub/Sub flow). |
| **`logs`** | Runtime storage location for structured service log files (e.g., `logs.json`). |
| **`types`** | Defines fundamental and highly shared data types (e.g., custom ID types, Context keys, UserInfo structs) used universally across services. |

-----

### 12\. Root Files and Automation

| File | Purpose |
| :--- | :--- |
| **`Makefile`** | Automation of common developer tasks: `build`, `test`, generating Protobuf code, and deployment tasks. |
| **`go.mod`, `go.sum`** | Go module and dependency management. |
| **`renovate.json5`** | Configuration file for Renovate, used for automated dependency updates. |