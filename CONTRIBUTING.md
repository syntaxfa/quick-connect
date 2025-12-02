# Contributing to Quick Connect

Thank you for your interest in contributing to **Quick Connect**!
We take code quality, maintainability, and scalability very seriously.
Please carefully follow these guidelines to ensure a smooth contribution process.

---

## 🛠 Prerequisites

Before starting, ensure you have the following installed:
- **Go**: Version `1.25` or higher.
- **Docker & Docker Compose**: For running the infrastructure.
- **Pre-commit**: For managing git hooks (`pip install pre-commit`).
- **Make**: To run project commands.

---

## 🛠 Workflow for Contributing

1. **Find or Propose an Issue:**
    - Check the open [Issues](../../issues) for available tasks.
    - If proposing a new task, **create an issue** and wait for approval from a maintainer.
    - Issues must include:
        - Clear title and description.
        - Proper labeling (`feature`, `bug`, `enhancement`, etc.).

2. **Branch Creation:**
    - Always branch from the `main` branch.
    - Branch naming convention:
      ```
      {type}/{short-description}
      ```
        - Example:
          `feat/add-user-authentication`
          `fix/fix-order-validation`
          `docs/update-api-docs`

3. **Local Setup & Pre-Commit:**
    - Clone the repository.
    - Install dependencies:
      ```bash
      go mod download
      ```
    - **Important:** Install pre-commit hooks to ensure your code passes checks before pushing:
      ```bash
      pre-commit install
      ```
    - Contributions that fail formatting (`golangci-lint`), linting, or validation checks **will be rejected** by CI.

4. **Code Changes:**
    - Follow the project's **Hexagonal Architecture** and directory structure (`app/`, `adapter/`, `pkg/`).
    - If you modify `.proto` files, ensure you regenerate the Go code (check `Makefile` or use `buf`).
    - Write clean, maintainable, and well-tested code.

---

## ✏️ Commit Message Guidelines

We strictly follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.

### Commit Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, white-space)
- `refactor`: Code refactoring without behavior change
- `perf`: Performance improvements
- `test`: Adding or fixing tests
- `build`: Changes to build system or dependencies
- `ci`: Changes to CI configuration
- `chore`: Miscellaneous chores

**Example:**
`feat(manager): implement jwt token validation`

---

## 🚀 Pull Request Process

- Submit a Pull Request only when your work is complete and ready for review.
- Pull Request title must follow the **Conventional Commits** format.
- Pull Request description must include:
    - Reference to the related issue (`Closes #issue_number`).
    - A clear explanation of the changes made.
    - List of any breaking changes if applicable.

- **Checklist before submitting:**
    - [ ] ran `make lint` and fixed all issues.
    - [ ] ran `make test-general` and all tests passed.
    - [ ] pre-commit hooks passed locally.

---

## 📚 Useful Commands (Makefile)

We use `make` to automate common tasks. Here are the most useful commands:

| Command | Description |
| :--- | :--- |
| `make help` | Show all available make commands. |
| `make build` | Build the project binaries. |
| `make test-general` | Run unit tests for core packages. |
| `make lint` | Run `golangci-lint` to check code style. |
| `make docker-up` | Start infrastructure (Postgres, Redis) via Docker. |
| `make proto-lint` | Lint Protocol Buffer files. |

---

## ✅ Review Standards

- Code must be clear, consistent, and easy to understand.
- Follow **SOLID principles** where applicable.
- No dead code, debug prints, or commented-out sections should be left.
- Be responsive and cooperative during code reviews.

---

## 📜 Code of Conduct

We are committed to providing a welcoming and respectful environment for all contributors.
Please read our [Code of Conduct](CODE_OF_CONDUCT.md).

Toxic behavior, disrespect, or harassment will not be tolerated.