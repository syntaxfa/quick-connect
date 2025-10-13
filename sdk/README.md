You got it. Here's the revised `README.md` in English, with the "Available SDKs" section removed.

-----

# Quick Connect SDKs

This directory contains the Software Development Kits (SDKs) for the various services within the **Quick Connect** project. The purpose of these SDKs is to simplify and standardize the integration process with our services across different platforms.

## Our Philosophy

By co-locating the SDKs within the same monorepo, we ensure that any changes to the service APIs are simultaneously reflected in their corresponding SDKs within a single, atomic commit. This approach guarantees constant compatibility between the client and server, preventing any versioning conflicts.

## Directory Structure

The SDKs are organized by service and target platform:

```
sdk/
├── <service-name>/
│   ├── <platform-name>/
│   │   ├── README.md         // Documentation specific to this SDK
│   │   └── ...               // SDK source code
│   └── ...
└── ...
```

* **`<service-name>`:** The name of the microservice for which the SDK is built (e.g., `chat`, `notification`).
* **`<platform-name>`:** The platform for which the SDK is designed (e.g., `web`, `android`, `ios`, `go`).

## How to Use

Each SDK has its own dedicated `README.md` file within its directory, containing detailed instructions for installation, setup, and usage examples. To get started, please refer to the documentation of the specific SDK you wish to use.

## Contributing

We welcome contributions for developing new SDKs or improving existing ones. If you plan to add a new SDK, please follow these steps:

1.  Open a new **Issue** in the repository to discuss your proposal.
2.  Once approved, create a new directory following the structure above (`sdk/<service>/<platform>`).
3.  Implement the SDK code, including comprehensive unit tests.
4.  Write a thorough `README.md` file for your SDK, including installation guides, usage instructions, and practical examples.
5.  Submit a **Pull Request**.