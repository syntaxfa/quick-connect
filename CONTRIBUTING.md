# Contributing to Quick Connect

Thank you for your interest in contributing to **Quick Connect**!
We take code quality, maintainability, and scalability very seriously.
Please carefully follow these guidelines to ensure a smooth contribution process.

---

## 🛠 Workflow for Contributing

1. **Find or Propose an Issue:**
    - Check the open [Issues](../../issues) for available tasks.
    - If proposing a new task, **create an issue** and wait for approval from a maintainer.
    - Issues must include:
        - Clear title and description.
        - Proper labeling (`feature`, `bug`, `enhancement`, etc.).
        - Approval label or comment from a maintainer.

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

3. **Pre-Commit Checks:**
    - Install and configure the provided **pre-commit hooks** and **linters**.
    - Contributions that fail formatting, linting, or validation checks **will be rejected**.

4. **Code Changes:**
    - Follow the project's architecture and directory structure precisely.
    - Write clean, maintainable, and well-tested code.
    - Update relevant documentation and tests if necessary.

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
- `revert`: Revert a previous commit

**Example:**
`feat(userapp): add JWT authentication`

---

## 🚀 Pull Request Process

- Submit a Pull Request only when your work is complete and ready for review.
- Pull Request title must follow the **Conventional Commits** format.
- Pull Request description must include:
    - Reference to the related issue (`Closes #issue_number`).
    - A clear explanation of the changes made.
    - List of any breaking changes if applicable.

- Requirements before submitting:
    - All pre-commit hooks must pass.
    - All tests must pass.
    - Code must be properly formatted.

---

## ✅ Review Standards

- Code must be clear, consistent, and easy to understand.
- Follow SOLID principles where applicable.
- No dead code, debug prints, or commented-out sections should be left.
- Be responsive and cooperative during code reviews.
- Reviews might require multiple iterations — be patient and professional.

---

## 📜 Code of Conduct

We are committed to providing a welcoming and respectful environment for all contributors.
Please:
- Be respectful and considerate.
- Provide constructive feedback.
- Maintain a positive and supportive tone.

Toxic behavior, disrespect, or harassment will not be tolerated.

---

## 📚 Useful Commands

Before submitting a pull request, run:

```bash
make build    # Build the project
make test     # Run all tests
make lint     # Run linters and format checks
