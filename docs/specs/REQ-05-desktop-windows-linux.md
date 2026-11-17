---
UseCase ID : REQ-05
Title: Ship Progen Desktop for Windows and Linux
Depends On: [REQ-03, REQ-04]
Status: IN PROGRESS
---

## Actors

- **Progen Desktop user on Windows or Linux**: someone who wants to generate Spring Boot projects through the Progen Desktop wizard, not only through the CLI, on a non-macOS machine.
- **Maintainer**: releases Progen Desktop and the CLI together from one `vX.Y.Z` tag (REQ-04), and checks build artifacts on pushes and pull requests.

## Preconditions

- Progen Desktop ships as a universal macOS app (REQ-03), released with the CLI from one tag (REQ-04).
- For a release: a `vX.Y.Z` tag matching `internal/version/version.go` has been pushed, and GoReleaser has published the GitHub release with the CLI binaries.
- For a user generating a project: the same prerequisites as on macOS. A JDK is on `PATH`, and an internet connection is available the first time Spotless formats a project.

## UseCase Details

### Goal

Make Progen Desktop available on the platforms the `progen` CLI is released for, other than macOS, which keeps its universal app. On those platforms it should behave the same as it does on macOS.

### Targets

The CLI's targets in `.goreleaser.yaml`, apart from macOS:

| Target          | Release archive                                  |
|-----------------|--------------------------------------------------|
| `windows/amd64` | `progen-desktop-vX.Y.Z-windows-amd64.zip`        |
| `linux/amd64`   | `progen-desktop-vX.Y.Z-linux-amd64.tar.gz`       |
| `linux/arm64`   | `progen-desktop-vX.Y.Z-linux-arm64.tar.gz`       |
| macOS (universal, unchanged) | `progen-desktop-vX.Y.Z-macos.zip`   |

### Trigger: release

The maintainer pushes a `vX.Y.Z` tag.

### Main success path: release

1. The release workflow checks the tag against `internal/version/version.go` and GoReleaser publishes the CLI release, as today.
2. Progen Desktop is built for macOS (universal), `windows/amd64`, `linux/amd64` and `linux/arm64`, each stamped with the tag's version (`X.Y.Z`).
3. Each build is packed into one archive, named as in the table above:
   - **Windows:** a ZIP of the app folder: `Progen Desktop.exe` and the files it needs.
   - **Linux:** a `.tar.gz` holding the executable, desktop entry, icon and `install.sh`.
   - **macOS:** the ZIP of `Progen Desktop.app`, unchanged.
4. Once every desktop build has succeeded, all four archives are uploaded to the tag's GitHub release.

### Failure path: a desktop build fails during a release

- If any desktop build fails (macOS, Windows or Linux), **no** desktop archive is uploaded to the release, and the workflow is marked failed.
- The CLI release GoReleaser already published stays as it is.

### Trigger and success path: CI

On every push and pull request, besides the existing macOS build and tests, CI cross-compiles and packages Progen Desktop for `windows/amd64`, `linux/amd64` and `linux/arm64`. It uploads the results as build artifacts. The Progen Desktop tests still run only on macOS.

### Main success path: using the app on Windows or Linux

1. The user downloads the archive for their platform from the release.
2. The user extracts it:
   - **Windows:** extract the ZIP, then run `Progen Desktop.exe`. The build is unsigned, so on first launch SmartScreen may show "Windows protected your PC"; the user chooses *More info → Run anyway*.
   - **Linux:** extract the archive, then run `install.sh` or the executable.
3. The wizard, Review step, preview (DSK-11), templates (DSK-06, DSK-07) and defaults (DSK-08) work as they do on macOS.
4. Generate saves the project as `<home>/Downloads/<app-name>.zip`. This is the same rule on every OS (DSK-04): the folder is created if it's missing, and a numeric suffix is added if the file exists. The system's localized or relocated Downloads folder (for example `XDG_DOWNLOAD_DIR`) is not consulted.
5. Templates and defaults are stored as JSON in the OS user-data directory (DSK-09):
   - **Windows:** `%AppData%\Progen Desktop\`
   - **Linux:** `~/.config/Progen Desktop/`, or `$XDG_CONFIG_HOME/Progen Desktop/` when `XDG_CONFIG_HOME` is set
   - **macOS:** `~/Library/Application Support/Progen Desktop/`, unchanged

### Edge cases

- **Missing Downloads folder** (common on minimal Linux setups): `<home>/Downloads` is created, and the ZIP is saved there.
- **Missing JDK or failed formatting:** generation fails and the user sees the same error as on macOS. On Windows the wrappers run as `mvnw.cmd` / `gradlew.bat`.
- **Missing Linux runtime libraries** (GTK 3 / Pango): the app can't start. The README lists these libraries as a requirement.

## Business Rules

- Progen Desktop is released for the platforms the CLI is released for. macOS ships as one universal app; Windows and Linux ship one archive per architecture.
- Every desktop archive in a release carries the tag's version, and all of them come from the same tag as the CLI (REQ-04).
- A release gets either all of its desktop archives or none of them.
- Archive names follow `progen-desktop-vX.Y.Z-<os>-<arch>.<ext>`. The macOS archive keeps `progen-desktop-vX.Y.Z-macos.zip`, so existing links keep working.
- Each release offers one archive per Windows/Linux target: no installers or system packages.
- The Windows build is shipped unsigned.
- The generated project is saved to `<home>/Downloads` on every OS.
- The generated projects are the same on every OS. Only the app's packaging and data directory differ.

## Acceptance Criteria

1. Pushing tag `vX.Y.Z` uploads exactly these desktop archives to that release:
   - `progen-desktop-vX.Y.Z-macos.zip`
   - `progen-desktop-vX.Y.Z-windows-amd64.zip`
   - `progen-desktop-vX.Y.Z-linux-amd64.tar.gz`
   - `progen-desktop-vX.Y.Z-linux-arm64.tar.gz`
2. The app in each archive reports version `X.Y.Z`.
3. If any one desktop build fails during a release, the release has no `progen-desktop-*` archives, the workflow is marked failed, and the CLI assets remain.
4. A push or pull request produces CI artifacts with the packaged Windows (`windows/amd64`) and Linux (`linux/amd64`, `linux/arm64`) apps, besides the macOS app. The Progen Desktop tests run on macOS only.
5. On Windows 10/11 (amd64), the extracted `Progen Desktop.exe` starts (after the SmartScreen prompt, if shown). The user can complete the wizard, and Generate saves `C:\Users\<user>\Downloads\<app-name>.zip`, with a suffix such as `<app-name>-1.zip` when the file exists.
6. On Linux (amd64 and arm64) with GTK 3 / Pango installed, the extracted app starts, the wizard completes, and Generate saves `~/Downloads/<app-name>.zip`. `~/Downloads` is created if it doesn't exist.
7. Templates and defaults saved on Windows appear under `%AppData%\Progen Desktop\`. On Linux they appear under `~/.config/Progen Desktop/`, or `$XDG_CONFIG_HOME/Progen Desktop/` when it's set. They're still there after the app restarts.
8. The preview on the Review step shows the same file tree on Windows and Linux as on macOS for the same choices.
9. The Progen Desktop README describes, for Windows and Linux:
   - how to download and install the app
   - the SmartScreen workaround
   - the Linux runtime libraries
   - where templates and defaults are stored on each OS
10. `docs/requirements.md` lists REQ-05. DSK-09 describes the OS user-data directory instead of a macOS-only path.

## In Scope

- Release builds and archives for `windows/amd64`, `linux/amd64` and `linux/arm64`, stamped with the tag's version.
- All-or-nothing upload of the desktop archives to the release.
- CI cross-compiling and packaging the Windows and Linux apps as build artifacts.
- Saving to `<home>/Downloads`, and storing data in the OS user-data directory, on Windows and Linux.
- README updates: install steps, the SmartScreen note, Linux runtime libraries, and data paths for each OS.
- Updating `docs/requirements.md` (the new REQ-05, and rewording DSK-09).

## Out of scope

- Windows arm64 (the CLI isn't released for it either).
- Windows installers (NSIS), `.deb`, `.rpm` and Arch `PKGBUILD` packages.
- Package-manager distribution (winget, Scoop, Homebrew cask, apt/rpm repositories, Flatpak, Snap).
- Code signing of the Windows build, and macOS signing changes.
- Auto-updates.
- Running the Progen Desktop tests on Linux or Windows runners.
- Using the system's localized or relocated Downloads folder, or letting the user choose where to save.
- Changes to the generated Spring Boot projects.
