# Contributing to AUTO AUDIT

First off, thank you for taking the time to contribute! Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

---

## Table of Contents
1. [Code of Conduct](#code-of-conduct)
2. [How Can I Contribute?](#how-can-i-contribute)
   - [Reporting Bugs](#reporting-bugs)
   - [Suggesting Enhancements](#suggesting-enhancements)
   - [Pull Requests](#pull-requests)
3. [Styleguides](#styleguides)
   - [Go Styleguide](#go-styleguide)
   - [Java Styleguide](#java-styleguide)
   - [Commit Messages](#commit-messages)
4. [Setting Up the Environment](#setting-up-the-environment)

---

## Code of Conduct

This project and everyone participating in it is governed by the [AUTO AUDIT Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior.

---

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please open an Issue using the **Bug Report** template.
* **Search first:** Check the existing Issues to see if the bug has already been reported.
* **Be specific:** Provide steps to reproduce, expected behavior, actual behavior, and include logs (both from the Go agent and Java Spring Boot application) where possible.

### Suggesting Enhancements

If you have an idea to improve the auditor agent (e.g., supporting a new vulnerability, adding a new effector/tool, or improving the POMDP math/decision loop), please open an Issue using the **Feature Request** template.
* Explain the *why* and *how* of your suggestion.
* If applicable, provide mockups, logs, or pseudocode.

### Pull Requests

Before submitting a Pull Request (PR):
1. Fork the repository and create your branch from `main`.
2. Ensure the code compiles and tests pass.
3. Make sure your code follows the styleguides below.
4. Reference the related issue(s) in your PR description.
5. Use the provided Pull Request Template to describe your changes.

---

## Styleguides

### Go Styleguide

* Follow standard Go styling conventions.
* Format your code using `gofmt` before committing.
* Run `go vet ./...` to check for common mistakes.
* Keep functions cohesive and document public APIs.

### Java Styleguide

* Follow standard Java Spring Boot and Maven layout conventions.
* Format code using standard IDE formatters.
* Keep class structures modular, obeying SOLID principles.

### Commit Messages

* Use clear and descriptive commit messages.
* We recommend starting with a prefix, such as:
  * `feat:` for new features
  * `fix:` for bug fixes
  * `docs:` for documentation updates
  * `refactor:` for code changes that neither fix bugs nor add features
  * `test:` for adding or modifying tests

---

## Setting Up the Environment

Refer to the [README.md](README.md) file for setup instructions. You can run the environment locally or using the Docker Compose setup provided in `deployments/docker-compose.yml`.
* **Go Agent:** Ensure Go 1.21+ is installed.
* **Java Vulnerable App:** Ensure JDK 17+ and Maven 3.8+ are installed.
