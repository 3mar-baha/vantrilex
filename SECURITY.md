# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| 1.x     | Yes       |
| 0.x     | No        |

## Reporting a vulnerability

Open a private security advisory on GitHub or contact the maintainers.
Do not file public issues for suspected vulnerabilities.

## Scope notes

- Runner CLIs execute with your user privileges; review workspace
  provisioning output before launching sessions.
- API keys are read from process environment at runtime and are never
  written to disk by the launcher.
