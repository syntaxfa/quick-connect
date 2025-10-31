---
name: Bug report
about: Create a report to help us improve Quick Connect
title: ''
labels: bug
assignees: alireza-fa

---

---
name: "Bug Report"
about: "Create a report to help us improve Quick Connect"
title: "[Bug]: "
labels: ["bug", "triage"]
body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report!

  - type: textarea
    id: what-happened
    attributes:
      label: "What happened?"
      description: "A clear and concise description of what the bug is. Please provide as much detail as possible."
    validations:
      required: true

  - type: textarea
    id: steps-to-reproduce
    attributes:
      label: "Steps to Reproduce"
      description: "How can we reproduce the issue?"
      placeholder: |
        1. Go to '...'
        2. Click on '....'
        3. Scroll down to '....'
        4. See error...
    validations:
      required: true

  - type: textarea
    id: expected-behavior
    attributes:
      label: "Expected Behavior"
      description: "A clear and concise description of what you expected to happen."
    validations:
      required: true

  - type: dropdown
    id: service
    attributes:
      label: "Which service is affected? (if known)"
      description: "Understanding which microservice has the bug helps us fix it faster."
      options:
        - "manager"
        - "chat"
        - "notification"
        - "admin"
        - "filehandler"
        - "Other / I don't know"
    validations:
      required: false

  - type: textarea
    id: environment
    attributes:
      label: "Environment (optional)"
      description: "Please provide details about the environment where you saw this bug."
      placeholder: |
        - OS: [e.g. Ubuntu 22.04]
        - Browser: [e.g. Chrome 105]
        - Go Version: [e.g. 1.25.2]
        - Docker Version: [e.g. 20.10.7]

  - type: checkboxes
    id: security
    attributes:
      label: "Is this a security vulnerability?"
      options:
        - label: "Yes, this is a security vulnerability."

  - type: markdown
    attributes:
      value: |
        **IMPORTANT:** If you checked the box above, **DO NOT POST DETAILS HERE.**
        Please report security vulnerabilities privately according to our [Security Policy](https://github.com/syntaxfa/quick-connect/security/policy), for example by contacting **Alireza Feizi** on Telegram at **[https://t.me/Ayeef](https://t.me/Ayeef)**.
---
