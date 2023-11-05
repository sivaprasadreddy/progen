# Project Generator CLI
A CLI tool for generating projects with sensible defaults.

## Installation

### Using Go Install

```shell
$ go install github.com/sivaprasadreddy/progen@latest
$ progen --help
```

### Linux or MacOS

`$ curl -s "https://raw.githubusercontent.com/sivaprasadreddy/progen/main/install.sh" | bash`

Create Symlink to the current version:

```
$ sudo ln -sf $HOME/.progen/current/progen /usr/local/bin/progen
$ progen --help
```

## Usage:

### Generate a Minimal Java Project

```shell
$ progen
? Choose application type: Minimal Java
? What is ApplicationName? myapp
? What is GroupID? com.mycompany
? What is ArtifactID? myapp
? What is Application Version? 1.0.0-SNAPSHOT
? What is base package? com.mycompany.myapp
? Choose Build Tool: Gradle
Project generated successfully
```

### Generate a Spring Boot Project
```shell
$ progen
? Choose application type: Spring Boot
? What is ApplicationName? myapp
? What is GroupID? com.mycompany
? What is ArtifactID? myapp
? What is Application Version? 1.0.0-SNAPSHOT
? What is base package? com.mycompany.myapp
? Choose Build Tool: Maven
Project generated successfully
```
