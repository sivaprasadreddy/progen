# Progen Desktop

Progen Desktop is a macOS desktop app for generating production-ready Spring Boot applications.
A step-by-step wizard walks you through the options; the generated project is saved to your `~/Downloads` folder as a ZIP archive.

Progen Desktop uses the same Spring Boot generator as the `progen` CLI, so the projects it creates are the same.
Progen Desktop gives you a graphical alternative to answering prompts in the terminal.

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

1. Download `progen-desktop-<version>-macos.zip` from the [Releases](https://github.com/sivaprasadreddy/progen/releases) page.
2. Unzip it and move `Progen Desktop.app` to your `Applications` folder.
3. Open the app. The first time, macOS will block it — see
   [Allowing the app to run on macOS](#allowing-the-app-to-run-on-macos) below.

### Prerequisites

- macOS 12 or later (universal build for Apple Silicon and Intel).
- A JDK installed and available on your `PATH`. Generated projects target **Java 25**.
  After generating the project, Progen Desktop formats the code by running the project's Maven or Gradle wrapper
  (`spotless:apply` / `spotlessApply`), which needs Java and an internet connection the first time
  it downloads Maven/Gradle and plugins. Generation fails if this step fails.

## Usage

1. Launch **Progen Desktop** and click **Get started**.
2. **Project** – enter the application name, version, Group ID and Artifact ID, and pick the build tool
   and task runner (Mise by default, or Taskfile with SDKMAN).
   The package name is filled in from the Group ID and Artifact ID until you change it.
3. **Frontend** – pick a frontend, and HTMX if you chose Thymeleaf or JTE.
4. **Database** – pick the database, persistence framework, and migration tool.
5. **Features** – select optional features and Spring AI model providers.
6. **Review** – check your selections and click **Generate**. You can change the application name, version,
   Group ID, Artifact ID and package right here, including for a project started from a template.
   Click **Preview** to see every generated file before saving: a file tree beside the selected file, read-only,
   with line numbers and syntax highlighting. **Generate** there saves exactly those files; **Back to Review**
   keeps the preview until you change a selection.
7. The project is saved as `~/Downloads/<application-name>.zip` (a numeric suffix is added if the
   file already exists). Unzip it and follow the generated project's `README.md`.

**Cancel**, on every step from **Project** to **Review**, asks for confirmation, then clears all your selections
and returns to **Welcome**.

### Defaults

The **Defaults** link below the steps sets the values a new project starts with on the **Project** step:
application version, Group ID (which also starts the package name), build tool, and task runner.
**Restore built-in values** fills in Progen Desktop's own (1.0.0, `dev.sivalabs`, Maven, Mise); click
**Save defaults** to keep them. A project you have already started, and templates, keep their own values.
Defaults are stored in `~/Library/Application Support/Progen Desktop/defaults.json`.

### Templates

A template is a saved selection of every step's options, so you can reuse a combination you
generate often. Progen Desktop comes with three, all using Maven, Mise, PostgreSQL, Flyway and Spring Data JPA:
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

Templates are stored in `~/Library/Application Support/Progen Desktop/templates.json`.

## Allowing the app to run on macOS

Progen Desktop releases are not signed with an Apple Developer ID or notarized by Apple, so macOS Gatekeeper
blocks the app the first time you open it with a message such as
*"Apple could not verify “Progen Desktop” is free of malware"*. You only need to do one of the following once.

### Option 1: System Settings (recommended)

1. Try to open `Progen Desktop.app`, then click **Done** (or **OK**) on the warning dialog.
2. Open **System Settings → Privacy & Security**.
3. Scroll down to the **Security** section. You will see *"Progen Desktop" was blocked to protect your Mac*.
   Click **Open Anyway**.
4. Confirm with **Open Anyway** again and authenticate with your password or Touch ID.

On macOS 14 (Sonoma) and earlier you can instead Control-click (right-click) `Progen Desktop.app` in Finder,
choose **Open**, and then click **Open** in the dialog.

### Option 2: Terminal

Remove the quarantine attribute that macOS adds to downloaded files:

```shell
xattr -dr com.apple.quarantine "/Applications/Progen Desktop.app"
```

If macOS reports that the app *"is damaged and can't be opened"*, this command fixes it as well.

> Only bypass Gatekeeper for apps you trust. If you prefer, build the app from source (see below).

## Building from source

Progen Desktop is built with [MyGo](https://mygo.run/) and its native UI: the interface is written in Go with
package `ui` and drawn by MyGo on the GPU, with no web view, HTML or JavaScript.

Requirements: [Go](https://go.dev/dl/) 1.27+. The MyGo CLI is a tool of the Go module, so
`go tool mygo` runs it with no separate install. Compiling the macOS 26 layered icon needs Xcode 26
or later; without it, the build falls back to the PNG icon.

Run these from the root of the progen repository; `cmd/desktop` is the MyGo project directory
(it holds `mygo.json` and the icons).

```shell
git clone https://github.com/sivaprasadreddy/progen.git
cd progen

# Run in development mode, rebuilding as the code changes (or: mise run desktop_dev)
go tool mygo dev ./cmd/desktop

# Build Progen Desktop.app and a DMG into cmd/desktop/build/darwin-arm64/ (or: mise run desktop_build)
go tool mygo build ./cmd/desktop

# Build a universal (Apple Silicon + Intel) app and DMG under cmd/desktop/build/
go tool mygo build -platform darwin/universal ./cmd/desktop

# Run the Progen Desktop tests, including the view tests that click through the wizard without a window
# (add -short to skip the test that generates a full project)
go test ./internal/desktop/...
```

Apps you build yourself are signed locally (ad-hoc) and don't trigger the Gatekeeper warning on your own Mac.

### Project structure

- `cmd/desktop/main.go` – the `desktop` command: wires the services into the wizard and opens the window, whose
  content is the wizard's native view.
- `cmd/desktop/mygo.json` – app name, identifier, version, and icons. The release workflow sets the version from the tag.
- `cmd/desktop/resources/icon.png` and `cmd/desktop/AppIcon.icon` – the app icon, and the layered icon that macOS 26 uses.
- `internal/desktop/config/` – `GeneratorConfiguration`, the project the user chooses (build tool, frontend, database,
  features, AI models), and `Validate`, which checks it.
- `internal/desktop/generator/` – `generator.Service`, which validates a configuration, maps it to a progen
  `springboot.ProjectConfig`, generates the project with `springboot.Generate` into a temporary directory
  (`Prepare`, which the preview shows), and zips it into `~/Downloads` (`Save`).
- `internal/desktop/store/` – `TemplateService` and `DefaultsService`, which keep templates and the Project step's
  defaults in the app's user data directory; `json.go` has the JSON file helpers they share.
- `internal/desktop/app/` – the window:
  - `wizard.go` – the wizard's state and actions: the selection, the steps, templates, defaults, and generating.
  - `view.go` – the screens: the sidebar, each step, and the Templates and Defaults pages, with their alerts.
  - `style.go` – the look: the Inter font, colors, the logo and icons, and the shared cards, fields, and buttons.
  - `options.go` – the choices each step offers, with their names and descriptions.
  - `view_test.go` – view tests that drive the wizard with MyGo's `ui.Tester`;
    `DESKTOP_SNAPSHOT_DIR=<dir> go test -run TestSnapshots ./internal/desktop/app` saves a picture of every screen.
  - `assets/` – the Inter font (with its license) and the logo, embedded into the app.
- `cmd/desktop/build/` – build output; generated by `mygo build`.

## Author

[SivaLabs](https://sivalabs.in)
