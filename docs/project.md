# Progen

Progen generates production-ready Spring Boot projects. It works offline, needs no other language SDK (no Node.js or Python), and puts in only the features and configuration that most projects use.

## Two ways to use it

- **Progen CLI** (`progen`): a single Go binary for macOS, Linux and Windows. It asks its questions as interactive prompts, or reads them from a `.progen.json` config file.
- **Progen Desktop**: a native desktop app for macOS, Windows and Linux, built with [MyGo](https://mygo.run/), that asks the same questions in a step-by-step wizard and saves the project as a ZIP in the `Downloads` folder of the user's home directory.

Both use the same generator (`internal/generators/springboot`), so they produce the same projects. They are released together under one version.

## What a generated project includes

- A Maven or Gradle build with Spotless formatting, and a mise or Taskfile (with SDKMAN) task runner.
- Persistence with Spring Data JPA, JdbcClient or jOOQ, on PostgreSQL, MySQL or MariaDB, with Flyway or Liquibase migrations.
- A REST API, or a frontend: Thymeleaf or JTE (optionally with HTMX), or an Angular or React SPA, styled with TailwindCSS.
- Spring Security: session-based login for server-rendered apps, JWT for REST APIs and SPAs.
- Spring Modulith, springdoc-openapi (Swagger UI) and BootUI.
- Optional integrations: Spring Cloud AWS, email (Mailpit), RabbitMQ, Redis caching, OpenTelemetry, Spring AI, Axelix and Kubernetes manifests.
- JUnit and Testcontainers tests, a Docker Compose setup for local development, GitHub Actions CI and Renovate.

## Goals

- Generate as little code as possible, and only the features most projects need.
- Work fully offline.
- Ship as one binary with no other runtime to install.
