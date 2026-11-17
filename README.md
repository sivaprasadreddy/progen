# Project Generator CLI
A CLI tool for generating Spring Boot projects.

This repository also contains **Progen Desktop**, a macOS desktop app that generates the same Spring Boot projects
through a step-by-step wizard. See [Progen Desktop](#progen-desktop-app) below.

[![Build](https://github.com/sivaprasadreddy/progen/actions/workflows/build.yml/badge.svg)](https://github.com/sivaprasadreddy/progen/actions/workflows/build.yml)
[![GoReleaser](https://github.com/sivaprasadreddy/progen/actions/workflows/release.yml/badge.svg)](https://github.com/sivaprasadreddy/progen/actions/workflows/release.yml)

## Motivation
There are a few good project generator tools exist, such as [JHipster](https://www.jhipster.tech/), [Bootify](https://bootify.io/), [generator-springboot](https://github.com/sivaprasadreddy/generator-springboot).

The reasons for creating this new **progen** CLI tool are:
* To generate the application with minimal code and include only the most commonly used features and configurations.
* We should be able to use the tool offline to create the projects.
* We should be able to use the tool without requiring the installation of other languages SDKs.
  For example, for generating a Java application, we shouldn't need to install Node.js or Python, etc.

## How to use Progen?
Progen provides a CLI tool and a desktop app to generate Spring Boot applications.

### Progen CLI
You can download the progen CLI binary from https://github.com/sivaprasadreddy/progen/releases.
If you have Go installed on your machine, you can use the following command to install `progen`:

```shell
$ go install github.com/sivaprasadreddy/progen/cmd/progen@latest
$ progen --help

# To use interactive shell
$ progen

# To create progen config file
$ progen init
Generated .progen.json successfully

# Generating a new project using .progen.json configuration
$ progen -c .progen.json
```

> [!IMPORTANT]
#### Allowing the app to run on macOS
> On MacOS, you may get the error **"progen cannot be opened because the developer cannot be verified"**.
> To fix this, you can select the progen binary in Finder, hold control and click, and select Open.
> You will then be prompted to confirm that you want to open the binary. After confirming, you can run progen as usual.

### Progen Desktop App
Progen Desktop is a native macOS app (built with [MyGo](https://mygo.run/)) that uses the same generator as the `progen` CLI.
You can download `progen-desktop-<version>-macos.zip` from https://github.com/sivaprasadreddy/progen/releases.

> [!IMPORTANT]
> Progen Desktop releases are not signed with an Apple Developer ID or notarized by Apple, so macOS Gatekeeper
> blocks the app the first time you open it with a message such as
> *"Apple could not verify “Progen Desktop” is free of malware"*.

#### Allowing the app to run on macOS
1. Try to open `Progen Desktop.app`, then click **Done** (or **OK**) on the warning dialog.
2. Open **System Settings → Privacy & Security**.
3. Scroll down to the **Security** section. You will see *"Progen Desktop" was blocked to protect your Mac*. Click **Open Anyway**.
4. Confirm with **Open Anyway** again and authenticate with your password or Touch ID.

On macOS 14 (Sonoma) and earlier you can instead Control-click (right-click) `Progen Desktop.app` in Finder,
choose **Open**, and then click **Open** in the dialog.

## Features
The progen generates a Spring Boot application with the following features configured:

| Category                 | Feature                | Details                                                                                                                                                                                 |
|--------------------------|------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Build**                | Project setup          | [Spring Boot](https://spring.io/projects/spring-boot) project with **Maven** or **Gradle**                                                                                              |
|                          | Code formatting        | [Spotless](https://github.com/diffplug/spotless)                                                                                                                                        |
|                          | Task runner            | [mise](https://mise.jdx.dev/) (default) or [Taskfile](https://taskfile.dev/) (with [SDKMAN](https://sdkman.io/) configuration) to execute commonly used commands                        |
| **Persistence**          | Persistence frameworks | **Spring Data JPA**, **Spring JdbcClient**, **jOOQ**                                                                                                                                    |
|                          | Databases              | **PostgreSQL**, **MySQL**, **MariaDB**                                                                                                                                                  |
|                          | Database migrations    | [Flyway](https://www.red-gate.com/products/flyway/community/) or [Liquibase](https://www.liquibase.com/)                                                                                |
| **Web & UI**             | Single Page Apps       | [Angular](https://angular.dev/) or [React](https://react.dev/) (with [Vite](https://vite.dev/)) and [TailwindCSS](https://tailwindcss.com/)                                             |
|                          | Server-side Web Apps   | [Thymeleaf](https://www.thymeleaf.org/) or [JTE](https://jte.gg/) with [TailwindCSS](https://tailwindcss.com/)                                                                          |
|                          | HTMX                   | [HTMX](https://htmx.org/) with [htmx-spring-boot](https://github.com/wimdeblauwe/htmx-spring-boot)                                                                                      |
|                          | API documentation      | Swagger UI using [springdoc-openapi](https://springdoc.org/)                                                                                                                            |
| **Security**             | Web App security       | [Spring Security](https://spring.io/projects/spring-security) configuration                                                                                                             |
|                          | REST API security      | **Spring Security JWT** configuration                                                                                                                                                   |
| **Architecture**         | Modular monolith       | [Spring Modulith](https://spring.io/projects/spring-modulith) configuration                                                                                                             |
| **Integrations**         | Messaging              | [Spring RabbitMQ](https://spring.io/projects/spring-amqp) configuration                                                                                                                 |
|                          | AWS                    | [Spring Cloud AWS](https://awspring.io/) with [Floci](https://floci.io/floci/) configuration                                                                                            |
|                          | Caching                | [Redis Caching](https://docs.spring.io/spring-boot/reference/io/caching.html#io.caching.provider.redis) configuration                                                                   |
|                          | Email                  | [MailPit](https://mailpit.axllent.org/) and [Testcontainers Mailpit](https://github.com/martinellich/testcontainers-mailpit)                                                            |
|                          | AI                     | [Spring AI](https://spring.io/projects/spring-ai) with **Ollama**, **OpenAI** and/or **Anthropic** chat models                                                                          |
| **Testing**              | Unit testing           | [JUnit](https://junit.org/)                                                                                                                                                             |
|                          | Integration testing    | [Testcontainers](https://java.testcontainers.org/) based testing and local development setup                                                                                            |
| **Local Dev & Deployment** | Docker Compose       | [Docker Compose](https://docs.docker.com/compose/) based local development setup                                                                                                        |
|                          | Kubernetes             | [Kubernetes](https://kubernetes.io/) manifests generated based on the selected features and databases                                                                                   |
| **Monitoring & Tooling** | Axelix                 | [Axelix](https://axelix.io/) Spring Boot starter and build plugin (optional, Maven and Gradle). |
|                          | BootUI                 | [BootUI](https://www.julien-dubois.com/boot-ui/)                                                                                                                                        |
| **CI & Maintenance**     | Continuous Integration | GitHub Actions configuration                                                                                                                                                            |
|                          | Dependency updates     | [Renovate](https://docs.renovatebot.com/)                                                                                                                                               |

## License
The **progen** is Open Source software released under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html).
