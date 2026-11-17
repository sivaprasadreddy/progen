# Tech Stack

## Progen

| Area                 | Technology                                                                                       |
|----------------------|--------------------------------------------------------------------------------------------------|
| Language             | [Go](https://go.dev/) 1.27 (`go.mod`)                                                             |
| CLI framework        | [Cobra](https://github.com/spf13/cobra): the `progen`, `init` and `version` commands, and the `-c/--config` flag |
| Interactive prompts  | [Huh](https://github.com/charmbracelet/huh) v2 (Charm: Bubble Tea, Bubbles, Lip Gloss)            |
| Templating           | Go `text/template`, with the templates embedded in the binary through `embed` (works offline)     |
| Config format        | JSON (`.progen.json`), read and written with `encoding/json`                                      |
| Desktop app          | [MyGo](https://mygo.run/) native UI (`ui` package, drawn on the GPU, no web view), run as `go tool mygo` |
| Desktop assets       | Inter font and logo, embedded into the app                                                        |
| Testing              | Go `testing`, [testify](https://github.com/stretchr/testify), MyGo `ui.Tester` for view tests     |
| Task runner          | [mise](https://mise.jdx.dev/) (`mise.toml`) and `make` (`Makefile`)                               |
| CI                   | GitHub Actions: `build.yml` runs on Ubuntu (Go and Java 25 Temurin) and on macOS for Progen Desktop |
| Release              | [GoReleaser](https://goreleaser.com/) builds the CLI binaries, archives and Homebrew formula; `release.yml` builds the universal macOS app |
| Dependency updates   | [Renovate](https://docs.renovatebot.com/) (`renovate.json`)                                       |

### Build targets

- **CLI**: a static binary (`CGO_ENABLED=0`) for macOS and Linux (amd64, arm64) and for Windows (amd64).
- **Progen Desktop**: a universal macOS app (Apple Silicon and Intel), macOS 12 or later. Its layered icon needs Xcode 26 to compile; without Xcode, the build uses the PNG icon instead.
