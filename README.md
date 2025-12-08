<div align="center">
  <img src="https://private-user-images.githubusercontent.com/86611004/521466724-6c2c7486-5c6f-46f2-ab1f-da4917082ff6.webp?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NjQ2OTU5MTAsIm5iZiI6MTc2NDY5NTYxMCwicGF0aCI6Ii84NjYxMTAwNC81MjE0NjY3MjQtNmMyYzc0ODYtNWM2Zi00NmYyLWFiMWYtZGE0OTE3MDgyZmY2LndlYnA_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjUxMjAyJTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI1MTIwMlQxNzEzMzBaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT1lMWE5ZjgyM2VmNjFiMDJiNzJkNDg5NDI2OGNlNzcxNDI1Y2IyOTE2YWY3ZTdjMGRkODZkNDBjODY3ZDlkMWNlJlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCJ9.iWRDyPw0XCuPgc1t0FFCSGw0vkXPX2Y6hdTmBSWeLVo" alt="Quick Connect Logo" width="200"/>

  <h1>Quick Connect</h1>

  <p>
    <b>Better Communication, More Customers. Quick Connect Cloud Platform.</b>
  </p>

  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/syntaxfa/quick-connect?style=flat-square&color=blue" alt="License">
  </a>
  <a href="https://go.dev">
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/Cloud--Native-Yes-success?style=flat-square&logo=docker" alt="Cloud Native">
  </a>
</div>

---

## 📖 About Quick Connect

**Quick Connect** is an open-source, cloud-native, and lightweight platform designed to enhance and optimize your customer engagement. Unlike similar third-party services, Quick Connect is **self-hosted**, ensuring that 100% of your data remains in your control.

### ⚡ The Quick Connect Difference
While many similar open-source alternatives are built using heavier frameworks like **Ruby on Rails** or **JavaScript**, Quick Connect is engineered with modern, high-performance technologies. This delivers a significantly faster, more responsive, and interactive experience for both you and your customers.

### 🪶 How Lightweight?
When we say lightweight, we mean it. You can deploy the entire Quick Connect stack using our [all-in-one](./deploy/all-in-one/deploy/compose.yml) Docker image, which is **less than 30MB** in size!

### 🌐 Live Demo
Experience the speed and features yourself on our live demo:
* **Website:** [Link to Website](https://quick-connect.syntaxfa.ir)
* **Username:** `quickconnect`
* **Password:** `quickconnect`

<div align="center">
  <br>
  <p><b>👇 Preview: Real-time Chat Service</b></p>
  <img src="https://github.com/user-attachments/assets/24201987-9382-492c-b3ad-d4e00a69076d" alt="Quick Connect Real-time Chat Demo" width="100%" style="border-radius: 10px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);"/>
</div>

## ✨ Key Features

Quick Connect comes packed with everything you need to build a modern engagement platform:

* **💬 Real-time Support Chat**
  Lightning-fast messaging powered by **WebSockets** and **Redis**. It ensures zero-latency communication between users and support agents.

* **📸 Interactive Stories (Trending Feature)**
  Boost user engagement by adding "Stories" to your app (similar to Instagram/Snapchat). Share ephemeral updates, news, or promotions directly with your users.

* **🔔 Smart & Multi-channel Notifications**
  A robust notification engine that supports **Email, SMS, and Push**.
    * **Smart Routing:** Automatically detects if a user is online (sends via WebSocket) or offline (fallbacks to Email/SMS).
    * **Multi-language:** Built-in i18n support for global applications.

* **📂 Flexible File Handler**
  A dedicated microservice for managing media uploads.
    * **Storage Agnostic:** Supports both **Local File System** and **S3-compatible** object storage (AWS S3, MinIO, etc.).

* **🎛️ Lightweight Admin Dashboard**
  Manage agents, users, and settings with a modern dashboard built using **Go Templates + HTMX**.
    * **No heavy SPA frameworks:** Extremely fast page loads and low resource usage.

## 🚧 Roadmap & Upcoming Features

We have ambitious plans for Quick Connect! Here is a glimpse of what's coming next:

### 📱 Client SDKs
- [ ] **Mobile SDKs:** Native libraries for **Android (Kotlin)** and **iOS (Swift)** to easily integrate chat into mobile apps.
- [ ] **Flutter Plugin:** A dedicated package for cross-platform mobile development.
- [ ] **React Native Component:** Plug-and-play component for RN apps.

### 🤖 AI & Automation
- [ ] **AI-First Support (RAG):** Upload your documents and FAQs. The system vectorizes your data to let AI answer incoming messages first, aiming to resolve **~90% of queries** instantly.
- [ ] **Smart Handover:** AI acts as the first line of defense. If the confidence score is low, the conversation is seamlessly transferred to a human agent with full context.
- [ ] **Sentiment Analysis:** Automatically analyze user mood (e.g., Angry, Neutral) to prioritize urgent tickets for human review.

### 📞 Communication & Media
- [ ] **Voice Messages:** Ability to record and send voice notes.
- [ ] **Video/Audio Calls:** Peer-to-peer calls using **WebRTC**.
- [ ] **File Preview:** Better preview for PDF and Office documents directly in chat.

### 🔌 Integrations
- [ ] **Telegram & WhatsApp Bridge:** Manage messages from Telegram Bot and WhatsApp Business directly in the Quick Connect dashboard.
- [ ] **Slack Integration:** Forward notifications to your team's Slack channel.
- [ ] **CRM Sync:** Sync user data with external CRMs like HubSpot or Salesforce.

### 📊 Analytics & Ops
- [ ] **Advanced Reporting:** Charts for agent response time, resolution rate, and busy hours.
- [ ] **Kubernetes Helm Charts:** Production-ready Helm charts for easy K8s deployment.

## 🏗️ Architecture

Quick Connect is architected as a **Modular Monolith**, giving you the ultimate flexibility in deployment. You are not forced into complex microservices if you don't need them.

### 🔄 Dual Deployment Modes
One of the unique features of Quick Connect is its **"Code-Level Monolith"** design. You can run the platform in two modes using the exact same codebase:

1.  **Microservices Mode (Scale):** Each component (Chat, Manager, Notification) runs as an independent container. Services communicate over the network via **gRPC**. Ideal for high-traffic, distributed environments (Kubernetes).
2.  **Monolith Mode (Speed & Simplicity):** All services run within a **single binary** (All-in-One). In this mode, inter-service communication bypasses the network completely and occurs via **direct function calls** (in-memory).
  * **Zero Network Latency:** No gRPC overhead between internal services.
  * **Easy Ops:** Deploy just one container/binary.

### 🧩 Service Modules

| Module | Responsibility | Key Tech Stack |
| :--- | :--- | :--- |
| **Manager** | The core identity provider handling **Authentication (JWT)**, User Management, and RBAC. | PostgreSQL |
| **Chat** | Manages real-time conversations, message persistence, and **WebSocket** connections. | Redis, PostgreSQL |
| **Notification** | A centralized engine for dispatching emails, SMS, and push notifications using the **Outbox Pattern**. | Redis Streams, Workers |
| **File Handler** | Handles secure media uploads (Local/S3). | S3 API |
| **Admin** | A server-side rendered dashboard for system management. | **HTMX**, Go Templates |

### 📐 Design Patterns & Best Practices
* **Hexagonal Architecture (Ports & Adapters):** Keeps the business logic isolated from external concerns (DB, API).
* **Outbox Pattern:** Ensures eventual consistency for notifications and events.
* **Abstracted Communication:** The code automatically switches between **gRPC** (remote) and **Function Calls** (local) based on the deployment configuration.

## 🛠️ Tech Stack

Quick Connect utilizes a modern, performance-oriented technology stack to ensure scalability and ease of maintenance.

### 🔙 Backend
* **Language:** [Go (Golang)](https://go.dev/) `1.25+` - For high-performance concurrency.
* **Framework:** [Echo v4](https://echo.labstack.com/) - High performance, extensible web framework.
* **Communication:**
  * **gRPC & Protobuf:** For efficient inter-service communication.
  * **WebSocket:** For real-time bi-directional events (Chat).
* **Database & Storage:**
  * **PostgreSQL:** Primary relational database (using `pgx` driver).
  * **Redis:** For caching, Pub/Sub, and session management.
  * **S3-Compatible Storage:** For file persistence (MinIO/AWS).
* **Key Libraries:**
  * `koanf`: Configuration management.
  * `ozzo-validation`: Data validation.
  * `sql-migrate`: Database migrations.
  * `cobra`: CLI command management.

### 🎨 Frontend (Admin Dashboard)
* **Architecture:** Server-Side Rendered (SSR).
* **Core:** [HTMX](https://htmx.org/) - For dynamic interactions without complex JS bundles.
* **Templating:** Go `html/template`.
* **Styling:** Custom CSS (No heavy CSS frameworks).

### ⚙️ DevOps & Infrastructure
* **Containerization:** Docker & Docker Compose.
* **Orchestration:** Kubernetes ready.
* **Observability:**
  * **OpenTelemetry (OTel):** Distributed tracing and metrics.
  * **Prometheus:** Metrics collection.
  * **Grafana:** Visualization (optional integration).
* **CI/CD:** GitHub Actions.
* **Dev Tools:** `Hybrid Development Environment`

## 🚀 Get Started

Quick Connect offers two deployment modes: **All-in-One** (recommended for testing & small setups) and **Microservices** (for scalable production).

### Option 1: All-in-One (Fastest Way) ⚡
Run the entire platform as a single monolithic container with minimal resource usage (<30MB image). You don't even need to build the code!

#### Prerequisites
* [Docker](https://www.docker.com/) & Docker Compose

#### Steps
1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/syntaxfa/quick-connect.git](https://github.com/syntaxfa/quick-connect.git)
    cd quick-connect/deploy/all-in-one/deploy
    ```

2.  **Setup Environment:**
    Copy the example configuration file. You can use the default values for a quick start.
    ```bash
    cp .env.example .env
    ```

3.  **Run with Docker Compose:**
    ```bash
    docker compose up -d
    ```
    *This will pull the `syntaxfa/quickconnect-all-in-one` image, setup Postgres & Redis, run migrations, and create a default superuser.*

4.  **Access the Dashboard:**
    Once the containers are healthy, open your browser:
    * **URL:** `http://localhost:2560`
    * **Default Username:** `alireza`
    * **Default Password:** `password`

---

### Option 2: Microservices (Production Architecture) 🏗️

For scalable environments, you can run Quick Connect as a set of distributed microservices where each component (Chat, Manager, Admin) runs in its own isolated container.

> **Note:** This Docker Compose setup demonstrates the full architecture running on a single machine for development/testing purposes. In a real production environment, these services would typically be orchestrated via Kubernetes across multiple nodes.

#### Steps

1.  **Navigate to the deployment directory:**
    If you haven't cloned the repo yet:
    ```bash
    git clone [https://github.com/syntaxfa/quick-connect.git](https://github.com/syntaxfa/quick-connect.git)
    cd quick-connect/deploy/microservice
    ```

2.  **Setup Environment:**
    Create the environment file from the sample. The default configuration connects all services automatically.
    ```bash
    cp .env.example .env
    ```

3.  **Run the Stack:**
    This command orchestrates all services (Manager, Chat, Admin, Infra) using the modular compose files.
    ```bash
    docker compose up -d
    ```

4.  **Access the Services:**
    Once the containers are healthy (it might take a few seconds for migrations to finish):

    * **Admin Dashboard:** `http://localhost:2560`
    * **Chat API:** `http://localhost:2530`
    * **Manager API:** `http://localhost:2520`

    > **Default Credentials:**
    > * **Username:** `alireza`
    > * **Password:** `password`

### 💻 For Developers (Build from Source)

If you want to contribute or modify the code, we recommend the **Hybrid Workflow**: run infrastructure via Docker and services locally via Go.

#### 1. Start Infrastructure (DB & Redis)
First, spin up the required databases (Postgres & Redis) without the application containers:
```bash
cd deploy/all-in-one/development
docker compose up -d
````

*This ensures you have a ready-to-use database environment mapped to your localhost ports.*

#### 2\. Configuration ⚙️

Quick Connect uses a layered configuration system. The priority order is:

1.  **Environment Variables** (Highest Priority - Overrides everything)
2.  **YAML Config Files** (Located in `deploy/<service>/config.yml`)
3.  **Default Values** (Lowest Priority)

> **Tip:** For local development, the services are pre-configured to read from `deploy/<service>/config.yml`. You can modify these files directly.

#### 3\. Run Services

Install dependencies and run each microservice individually. For example, to start the **Manager Service**:

```bash
# 1. Download Dependencies
go mod download

# 2. Run Database Migrations
go run cmd/manager/main.go migrate up

# 3. Start the Server
go run cmd/manager/main.go server
```

You can repeat this process for `cmd/chat/main.go`, `cmd/notification/main.go`, etc.

## 📚 Documentation

Explore the detailed documentation to understand how to integrate, customize, and extend Quick Connect:

* **[Project Structure](./docs/structure.md):** A deep dive into the directory layout, hexagonal architecture, and module organization.
* **[Client SDKs](./sdk/README.md):** Official SDKs.
* **API Reference (OpenAPI/Swagger):**
    * The API definitions are located within each service's directory (e.g., `app/chat/docs/chat_swagger.yaml`).
    * You can import these files into Postman or Swagger UI to inspect endpoints and schemas.
* **[Deployment Examples](./example/deploy):** Advanced configurations for different environments, including Kubernetes or separate VM setups.

## 🤝 Contributing

We enthusiastically welcome contributions from the community! Whether it's fixing bugs, improving documentation, or suggesting new features, your help is appreciated.

To get started:
1.  Read our **[Contributing Guide](CONTRIBUTING.md)** to understand the workflow and coding standards.
2.  Check out the [Open Issues](https://github.com/syntaxfa/quick-connect/issues) for tasks that interest you.
3.  Please review our **[Code of Conduct](CODE_OF_CONDUCT.md)** to ensure a welcoming environment for everyone.

> **Need Help?** If you have questions or need coordination before starting a large feature, feel free to reach out to the maintainer on Telegram: **[@Ayeef](https://t.me/Ayeef)**.

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0**. See the [LICENSE](LICENSE) file for details.
