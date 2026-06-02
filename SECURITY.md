# Security Policy

## Supported Versions

Currently, security updates and patches are active for the following versions:

| Version | Supported |
| ------- | --------- |
| main    | :white_check_mark: |
| < 1.0.0 | :x: |

---

## Reporting a Vulnerability

**IMPORTANT:** This repository is an educational/demonstrative pentest agent stand illustrating Log4Shell exploitation and remediation. However, security vulnerabilities in the agent itself, the OOB listeners, or the surrounding architecture should still be reported responsibly.

If you discover a security vulnerability, please do **not** report it via a public issue. Instead:
1. Open a Draft Security Advisory on GitHub (if available).
2. Or contact the project maintainers directly via their public profiles or by opening a general issue requesting private contact instructions.

We will acknowledge your report within 48 hours and work with you to analyze and remediate the issue promptly.

---

## Safety and Educational Disclaimer

This project is intended strictly for educational purposes, security research, and compliant vulnerability auditing. Running the target container exposes a vulnerable Java Spring Boot application (CVE-2021-44228). Do not deploy the vulnerable application image in a production or public network environment without proper network segmentation and access control.
