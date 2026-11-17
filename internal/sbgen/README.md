# SBGen

SBGen is a macOS desktop app for generating production-ready Spring Boot applications.
A step-by-step wizard walks you through the options; the generated project is saved to your
`~/Downloads` folder as a ZIP archive.

## What you can generate

| Step     | Options                                                                                                 |
|----------|---------------------------------------------------------------------------------------------------------|
| Project  | Application name and version, Group ID, Artifact ID, base package, build tool (Maven or Gradle), task runner (Mise or Taskfile with SDKMAN) |
| Frontend | None (REST API), Thymeleaf, JTE, Angular, or React; optional HTMX for Thymeleaf and JTE                   |
| Database | PostgreSQL, MySQL, or MariaDB; Spring Data JPA, Spring JdbcClient, or jOOQ; Flyway or Liquibase          |
| Features | Spring Cloud AWS, Email, RabbitMQ, Redis Caching, OpenTelemetry, Kubernetes manifests, Axelix            |
| AI       | Spring AI chat model providers: Ollama, OpenAI, Anthropic                                                |

Security is configured automatically based on the frontend: session-based login for server-rendered
frontends (Thymeleaf, JTE) and JWT for REST API, Angular, and React applications.

## Install

1. Download `sbgen-<version>-macos.zip` from the [Releases](https://github.com/sivaprasadreddy/progen/releases) page.
2. Unzip it and move `SBGen.app` to your `Applications` folder.
3. Open the app. The first time, macOS will block it — see
   [Allowing the app to run on macOS](#allowing-the-app-to-run-on-macos) below.

### Prerequisites

- macOS 12 or later (universal build for Apple Silicon and Intel).
- A JDK installed and available on your `PATH`. Generated projects target **Java 25**.
  After generating the project, SBGen formats the code by running the project's Maven or Gradle wrapper
  (`spotless:apply` / `spotlessApply`), which needs Java and an internet connection the first time
  it downloads Maven/Gradle and plugins. Generation fails if this step fails.

## Usage

1. Launch **SBGen** and click **Get started**.
2. **Project** – enter the application name, version, Group ID and Artifact ID, and pick the build tool
   and task runner (Mise by default, or Taskfile with SDKMAN).
   The package name is filled in from the Group ID and Artifact ID until you change it.
3. **Frontend** – pick a frontend, and HTMX if you chose Thymeleaf or JTE.
4. **Database** – pick the database, persistence framework, and migration tool.
5. **Features** – select optional features and Spring AI model providers.
6. **Review** – check your selections and click **Generate**. You can change the application name, version,
   Group ID, Artifact ID and package right here, including for a project started from a template. You can change the application name, version, Group ID, Artifact ID and package right here, including for a project started from a template.
7. The project is saved as `~/Downloads/<application-name>.zip` (a numeric suffix is added if the
   file already exists). Unzip it and follow the generated project's `README.md`.

**Cancel**, on every step from **Project** to **Review**, asks for confirmation, then clears all your selections
and returns to **Welcome**.

### Defaults

The **Defaults** link below the steps sets the values a new project starts with on the **Project** step:
application version, Group ID (which also starts the package name), build tool, and task runner.
**Restore built-in values** fills in SBGen's own (1.0.0, `dev.sivalabs`, Maven, Mise); click
**Save defaults** to keep them. A project you have already started, and templates, keep their own values.
Defaults are stored in `~/Library/Application Support/SBGen/defaults.json`.

### Templates

A template is a saved selection of every step's options, so you can reuse a combination you
generate often. SBGen comes with three, all using Maven, Mise, PostgreSQL, Flyway and Spring Data JPA:
**REST API** (no frontend), **Thymeleaf + HTMX**, and **Angular**. They are yours to update or delete
like any other.

- **Save:** on the **Review** step, click **Save as template**, give it a name and click **Save template**.
  If you started from a template, you can **Update** that template or **Save as new**.
- **Use:** the **Welcome** step lists your templates, and the **Templates** link below the steps shows
  them all at any time; its **Back** button returns to the step you were on. Click a template to go
  straight to **Review**, and use **Back** to change any selection before you generate.
- **Update or delete:** each template in either list has **Update**, which opens it at the **Project**
  step so you can change its selections and save it again from **Review**, and **Delete**, which asks for
  confirmation first.

Templates are stored in `~/Library/Application Support/SBGen/templates.json`.

## Allowing the app to run on macOS

SBGen releases are not signed with an Apple Developer ID or notarized by Apple, so macOS Gatekeeper
blocks the app the first time you open it with a message such as
*"Apple could not verify “SBGen” is free of malware"*. You only need to do one of the following once.

### Option 1: System Settings (recommended)

1. Try to open `SBGen.app`, then click **Done** (or **OK**) on the warning dialog.
2. Open **System Settings → Privacy & Security**.
3. Scroll down to the **Security** section. You will see *"SBGen" was blocked to protect your Mac*.
   Click **Open Anyway**.
4. Confirm with **Open Anyway** again and authenticate with your password or Touch ID.

On macOS 14 (Sonoma) and earlier you can instead Control-click (right-click) `SBGen.app` in Finder,
choose **Open**, and then click **Open** in the dialog.

### Option 2: Terminal

Remove the quarantine attribute that macOS adds to downloaded files:

```shell
xattr -dr com.apple.quarantine /Applications/SBGen.app
```

If macOS reports that the app *"is damaged and can't be opened"*, this command fixes it as well.

> Only bypass Gatekeeper for apps you trust. If you prefer, build the app from source (see below).

