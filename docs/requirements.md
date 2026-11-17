# Requirements

Each requirement is a feature or use case. Status values: `COMPLETED`, `IN PROGRESS`, `PLANNED`, `NOT_IMPLEMENTED`.

## 1. Distribution and platform

| ID     | Requirement                                                                                                     | Status    |
|--------|-----------------------------------------------------------------------------------------------------------------|-----------|
| REQ-01 | Ship the CLI as a single offline binary for macOS, Linux and Windows, with no other language SDK needed          | COMPLETED |
| REQ-02 | Install the CLI with `go install`, from GitHub Releases, or with Homebrew (GoReleaser formula)                   | COMPLETED |
| REQ-03 | Ship Progen Desktop as a universal macOS app (Apple Silicon and Intel, macOS 12+), uploaded as a ZIP to each release | COMPLETED |
| REQ-04 | Release the CLI and Progen Desktop together from one `vX.Y.Z` tag that must match `internal/version/version.go`  | COMPLETED |

## 2. CLI use cases

| ID     | Requirement                                                                                          | Status    |
|--------|------------------------------------------------------------------------------------------------------|-----------|
| CLI-01 | Generate a project by answering interactive prompts (`progen`)                                        | COMPLETED |
| CLI-02 | Create a `.progen.json` config file (`progen init`)                                                   | COMPLETED |
| CLI-03 | Generate a project from a config file (`progen -c .progen.json`)                                      | COMPLETED |
| CLI-04 | Replace invalid values in a config file with defaults, printing a warning for each one                | COMPLETED |
| CLI-05 | Refuse to generate when a directory with the application's name already exists in the working directory | COMPLETED |
| CLI-06 | Print the version (`progen version`)                                                                  | COMPLETED |

## 3. Progen Desktop use cases

| ID     | Requirement                                                                                                            | Status    |
|--------|------------------------------------------------------------------------------------------------------------------------|-----------|
| DSK-01 | Walk the user through the options in a wizard: Welcome → Project → Frontend → Database → Features → Review                | COMPLETED |
| DSK-02 | Fill in the base package from the Group ID and Artifact ID until the user edits it                                       | COMPLETED |
| DSK-03 | Let the user edit the name, version, Group ID, Artifact ID and package on the Review step                               | COMPLETED |
| DSK-04 | Generate the project off the main thread and save it as `~/Downloads/<app-name>.zip`, adding a numeric suffix if the file exists | COMPLETED |
| DSK-05 | Ask for confirmation before Cancel clears the selections and returns to Welcome                                          | COMPLETED |
| DSK-06 | Save templates (a selection of every step's options), and use, update and delete them (deleting asks for confirmation)   | COMPLETED |
| DSK-07 | Include three built-in templates: REST API, Thymeleaf + HTMX, Angular                                                    | COMPLETED |
| DSK-08 | Let the user set the defaults for new projects (version, Group ID, build tool, task runner) and restore the built-in values | COMPLETED |
| DSK-09 | Store templates and defaults as JSON in `~/Library/Application Support/Progen Desktop/`                                  | COMPLETED |
| DSK-10 | Configure security from the chosen frontend: session login for Thymeleaf/JTE, JWT for REST API, Angular and React        | COMPLETED |
| DSK-11 | Preview the generated project (file tree and read-only file contents) from the Review step before saving the ZIP         | COMPLETED |

## 4. Generated project: project setup

| ID     | Requirement                                                                                  | Status    |
|--------|----------------------------------------------------------------------------------------------|-----------|
| GEN-01 | Generate a Spring Boot project from the application name, version, Group ID, Artifact ID and base package, targeting Java 25 | COMPLETED |
| GEN-02 | Validate the project config: the name must be a plain directory name, the required fields must be set, and the choices must be valid | COMPLETED |
| GEN-03 | Build tool: Maven or Gradle (with wrapper)                                                    | COMPLETED |
| GEN-04 | Code formatting with Spotless, applied after generation                   | COMPLETED |
| GEN-05 | Task runner: mise (default), or Taskfile with SDKMAN                                          | COMPLETED |
| GEN-06 | Write `.gitignore`, a project `README.md` and the `.progen.json` used, into the project        | COMPLETED |

## 5. Generated project: persistence

| ID     | Requirement                                              | Status    |
|--------|----------------------------------------------------------|-----------|
| DB-01  | Persistence: Spring Data JPA, Spring JdbcClient or jOOQ  | COMPLETED |
| DB-02  | Database: PostgreSQL, MySQL or MariaDB                   | COMPLETED |
| DB-03  | Database migrations: Flyway or Liquibase                 | COMPLETED |

## 6. Generated project: web and UI

| ID     | Requirement                                                                | Status    |
|--------|----------------------------------------------------------------------------|-----------|
| WEB-01 | REST API only (no frontend)                                                 | COMPLETED |
| WEB-02 | Server-rendered web app with Thymeleaf or JTE, styled with TailwindCSS       | COMPLETED |
| WEB-03 | Optional HTMX (with htmx-spring-boot) for Thymeleaf and JTE                 | COMPLETED |
| WEB-04 | Single-page app with Angular, or React with Vite, styled with TailwindCSS    | COMPLETED |
| WEB-05 | API documentation with Swagger UI (springdoc-openapi)                       | COMPLETED |

## 7. Generated project: security and architecture

| ID     | Requirement                                                     | Status    |
|--------|-----------------------------------------------------------------|-----------|
| SEC-01 | Spring Security with session-based login for web apps           | COMPLETED |
| SEC-02 | Spring Security with JWT for REST APIs and SPAs                 | COMPLETED |
| ARC-01 | Modular monolith with Spring Modulith                           | COMPLETED |

## 8. Generated project: optional integrations

| ID     | Requirement                                                                  | Status    |
|--------|------------------------------------------------------------------------------|-----------|
| INT-01 | Spring Cloud AWS with Floci for local development                             | COMPLETED |
| INT-02 | Email with Mailpit and Testcontainers Mailpit                                 | COMPLETED |
| INT-03 | Messaging with Spring RabbitMQ                                                | COMPLETED |
| INT-04 | Redis caching                                                                 | COMPLETED |
| INT-05 | OpenTelemetry                                                                 | COMPLETED |
| INT-06 | Spring AI with one or more chat models: Ollama, OpenAI, Anthropic             | COMPLETED |
| INT-07 | Axelix Spring Boot starter and build plugin (Maven and Gradle)                | COMPLETED |
| INT-08 | BootUI                                                                        | COMPLETED |

## 9. Generated project: testing, local development and CI

| ID     | Requirement                                                                     | Status    |
|--------|---------------------------------------------------------------------------------|-----------|
| TST-01 | Unit tests with JUnit                                                           | COMPLETED |
| TST-02 | Integration tests and local development setup with Testcontainers               | COMPLETED |
| DEV-01 | Docker Compose setup for local development, with the services the selected features need | COMPLETED |
| DEV-02 | Kubernetes manifests for the selected features and database (optional)          | COMPLETED |
| CI-01  | GitHub Actions CI workflow                                                      | COMPLETED |
| CI-02  | Dependency updates with Renovate                                                | COMPLETED |
