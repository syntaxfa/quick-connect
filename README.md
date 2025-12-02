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
* **Website:** [Link to Demo Website](https://demo-quick-connect.syntaxfa.ir)
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
* **Dev Tools:** `Air` (Live reload), `Make` (Build automation).

## 🚀 Get Started
## 📚 Documentation
## 🤝 Contributing
## 📄 License