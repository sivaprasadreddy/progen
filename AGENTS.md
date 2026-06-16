# AGENTS.md

This file provides guidance to AI Agents when working with code in this repository.

## What This Project Is

**Progen** is a CLI tool (written in Go) that generates complete, production-ready Spring Boot projects. 
It works offline, ships as a single binary, and uses interactive prompts or a JSON config file to scaffold a full Java project — build files, Java source, tests, Docker Compose, GitHub Actions, database migrations, security setup, and more.

## Commands

```bash
make run         # Run interactively (go run main.go)
make build       # Compile binaries for macOS, Linux, Windows → dist/
make test        # Run tests in short mode (fast, skips permutation tests)
make all_tests   # Run all tests including slow integration tests
make fmt         # Format Go source with goimports
make clean       # Remove dist/
```

Running a single test:
```bash
go test -v -run TestGenerateRestApiProjectWithMavenAndPostgresql ./...
```
