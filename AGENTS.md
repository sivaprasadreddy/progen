# AGENTS.md

This file provides guidance to AI Agents when working with code in this repository.

## What This Project Is

**Progen** is a CLI tool (written in Go) that generates complete, production-ready Spring Boot projects. 
It works offline, ships as a single binary, and uses interactive prompts or a JSON config file to scaffold a full Java project — build files, Java source, tests, Docker Compose, GitHub Actions, database migrations, security setup, and more.

The repository also contains **Progen Desktop**, a desktop app for macOS, Windows and Linux built with [MyGo](https://mygo.run/) (native UI, no web view) that generates the same projects through a wizard. Both are released together with the version in `internal/version/version.go`.

## Layout

- `cmd/progen/` – the CLI entry point; `internal/cli/` – Cobra commands and interactive prompts.
- `internal/generators/springboot/` – the Spring Boot generator (`springboot.Generate`) and its templates, shared by the CLI and Progen Desktop.
- `cmd/desktop/` – the Progen Desktop `main` package and its MyGo project directory: `mygo.json`, `resources/icon.png`, `AppIcon.icon`, and `README.md` (user docs).
- `.github/scripts/package-desktop.sh` – packs a Progen Desktop build into its release archive (`macos`, `windows-amd64`, `windows-arm64`, `linux-amd64`, `linux-arm64`); CI and the release workflow use it.
- `internal/desktop/` – Progen Desktop's packages: `config` (the user's choices and `Validate`), `generator` (maps them to `springboot.ProjectConfig`, generates and zips into `~/Downloads`), `store` (templates and defaults in the user data directory), and `app` (the wizard's state in `wizard.go`, the screens in `view.go`, the look in `style.go`, choices in `options.go`). `config`, `generator` and `store` must not import `app` or MyGo's `ui`.
- The `progen` CLI must not import `internal/desktop` or MyGo (`go list -deps ./cmd/progen | grep mygo` should print nothing).

## Commands

```bash
make run         # Run interactively (go run ./cmd/progen)
make build       # Compile binaries for macOS, Linux, Windows → dist/
make test        # Run tests in short mode (fast, skips permutation tests)
make all_tests   # Run all tests including slow integration tests
make fmt         # Format Go source code
make clean       # Remove dist/
```

Progen Desktop (the MyGo project directory is `cmd/desktop`; pass it to every `mygo` command):
```bash
go tool mygo dev ./cmd/desktop                               # Run the app, rebuilding on changes (make desktop_dev)
go tool mygo build ./cmd/desktop                             # Build Progen Desktop.app and a DMG into cmd/desktop/build/ (make desktop_build)
go tool mygo build -platform darwin/universal ./cmd/desktop  # Universal macOS build
go tool mygo build -platform windows/amd64,windows/arm64,linux/amd64,linux/arm64 -skip-installer ./cmd/desktop  # Windows and Linux builds (no cgo, any machine)
go test -short ./internal/desktop/...                        # Progen Desktop tests, including ui.Tester view tests
go tool mygo vet ./internal/desktop/app                      # MyGo-specific checks; one package directory at a time
```
The Progen Desktop view tests need a text engine: they run on macOS (CI runs them on `macos-latest`; the Ubuntu job skips `internal/desktop`). CI also cross-compiles and packages the Windows and Linux apps on Ubuntu, without running tests there.

## Progen Desktop Native UI Practices

- The view runs on the main thread, possibly several times a frame. Change state in actions (`Clicked`, `OnClick`, a widget's `Changed`), not unconditionally while building.
- Don't keep `ui.Element` values or a `*ui.Context` in the wizard or in goroutines; `go tool mygo vet` reports it.
- Keep slow work off the main thread: generation runs through `wizard.async`, and its result is applied with `Window.Update`.
- Give a widget a key (`c.Key(...)`) where siblings would otherwise share state, as option cards in a loop.
- For alerts whose action discards data, list the action first and Cancel last: `ui.AlertDialog`'s last button is the default for Enter.
- Add a choice to a step by adding its constant in `internal/desktop/config/config.go`, accepting it in `Validate`, mapping it in `toProjectConfig` (`internal/desktop/generator/generator.go`), and listing it in `internal/desktop/app/options.go`.
- Cover new view behavior with a `ui.Tester` test in `internal/desktop/app/view_test.go`.
- MyGo docs for the pinned version: `go list -m -f '{{.Dir}}' github.com/egoist/mygo`, under `docs/` (native UI under `docs/ui/`). `.agents/skills/` has the `mygo-maintenance` and `mygo-native-modules` skills.

## Releasing

Bump `Version` in `internal/version/version.go`, commit, and push a matching `vX.Y.Z` tag. `.github/workflows/release.yml` checks that the tag matches, runs GoReleaser for the `progen` binaries and Homebrew formula, and builds Progen Desktop for macOS (universal), Windows (amd64, arm64) and Linux (amd64, arm64), writing the tag's version into `cmd/desktop/mygo.json`. Only when every desktop build succeeds does it upload `progen-desktop-vX.Y.Z-macos.zip`, `-windows-amd64.zip`, `-windows-arm64.zip`, `-linux-amd64.tar.gz` and `-linux-arm64.tar.gz` to the same release.

