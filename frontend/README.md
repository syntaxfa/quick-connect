# Quick Connect Frontend Applications

This directory contains all user-facing applications for the Quick Connect project. This includes web-based admin panels, management dashboards, and potentially mobile or desktop clients.

Each application in this directory is a standalone project that communicates with our backend services via their respective APIs.

## Our Philosophy

By separating frontend applications from the backend services (`app/`) and SDKs (`sdk/`), we maintain a clean and scalable architecture. This separation allows for:

* **Independent Development:** Frontend and backend teams can work in parallel.
* **Independent Deployment:** A change in an admin panel doesn't require redeploying a backend service.
* **Technological Freedom:** Each frontend application can use the technology best suited for its purpose (e.g., React, Vue, Svelte, or even plain HTML for web; Swift/Kotlin for mobile).

## Directory Structure

All applications are organized by their specific function:

```
frontend/
├── <application-name>/
│   ├── README.md         // Setup and usage instructions for this application
│   └── ...               // Source code, assets, etc.
└── ...
```

* **`<application-name>`:** A descriptive name for the application (e.g., `chat-admin-panel`, `ios-client`).

## Getting Started

Each application is self-contained and has its own setup, development, and build process.

To get started with a specific application, **please navigate to its directory and follow the instructions in its dedicated `README.md` file.**

## Contributing

If you wish to add a new frontend application, please follow these guidelines:

1.  Create a new directory inside `frontend/` with a descriptive name.
2.  Initialize your project using the tools and structure appropriate for its platform.
3.  Ensure you include a `README.md` file with clear instructions on how to set up, run, and build the application.
4.  Configure the application to communicate with the backend services, preferably using environment variables for API endpoints.
5.  Submit a Pull Request for review.