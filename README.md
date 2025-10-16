# Project Generator CLI
A CLI tool for generating Spring Boot projects.

## Motivation
There are a few good project generator tools exist, such as [JHipster](https://www.jhipster.tech/), [Bootify](https://bootify.io/), [generator-springboot](https://github.com/sivaprasadreddy/generator-springboot).

The reasons for creating this new **progen** CLI tool are:
* To generate the application with minimal code and include only the most commonly used features and configurations.
* We should be able to use the tool offline to create the projects.
* We should be able to use the tool without requiring the installation of other languages SDKs.
  For example, for generating a Java application, we shouldn't need to install Node.js or Python, etc.

## features
The progen generates a Spring Boot application with the following features configured:

* [Spring Boot](https://spring.io/projects/spring-boot) project with Maven and Gradle support
* [Spring Data JPA](https://spring.io/projects/spring-data-jpa) integration with MySQL, Postgresql, MariaDB support
* [Flyway](https://www.red-gate.com/products/flyway/community/) and [Liquibase](https://www.liquibase.com/) database migration support
* [Spring Security](https://spring.io/projects/spring-security) Configuration for WebApp
* JWT based Spring Security Configuration for REST APIs
* [Spring Modulith](https://spring.io/projects/spring-modulith) Configuration
* [Spring Cloud AWS](https://awspring.io/) support with [LocalStack](https://www.localstack.cloud/) configuration
* Swagger UI Integration using [springdoc-openapi](https://springdoc.org/)
* [Testcontainers](https://java.testcontainers.org/) based Testing and Local Development Setup
* [Docker Compose](https://docs.docker.com/compose/) based Local Development Setup
* GitHub Actions Configuration
* Code formatting using [Spotless](https://github.com/diffplug/spotless)
* [JUnit 5](https://junit.org/junit5/)
* [Taskfile](https://taskfile.dev/) to execute commonly used commands


## Installation

### Using OS specific binary
Download the latest binary for your OS and Architecture(arm64, x86_64) 
from https://github.com/sivaprasadreddy/progen/releases.

### Using Go Install
If you have Go installed on your machine, you can use the following command to install `progen`:

```shell
$ go install github.com/sivaprasadreddy/progen@latest
$ progen --help
```

> [!IMPORTANT]
> On MacOS, you may get the error **"progen cannot be opened because the developer cannot be verified"**.
> To fix this, you can select the progen binary in Finder, hold control and click, and select Open.
> You will then be prompted to confirm that you want to open the binary. After confirming, you can run progen as usual.

## Usage:

```shell
$ progen
```

## License
The **progen** is an Open Source software
released under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html).
