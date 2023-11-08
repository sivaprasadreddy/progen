# Project Generator CLI
A CLI tool for generating projects with sensible defaults.

## Installation

### Using Go Install

```shell
$ go install github.com/sivaprasadreddy/progen@latest
$ progen --help
```

### Linux or MacOS
Download the latest binary for your OS from https://github.com/sivaprasadreddy/progen/releases.

```shell
$ cd $HOME/Downloads
$ curl --location --progress-bar "https://github.com/sivaprasadreddy/progen/releases/download/v0.0.2/progen_Darwin_arm64.tar.gz" > "progen.tar.gz"
$ tar -xf progen.tar.gz progen
$ ./progen --help
# move progen binary to somewhere in your PATH so that you can access progen from anywhere 
```

## Usage:

```shell
$ progen
? Choose application type:  [Use arrows to move, type to filter]
> Minimal Java
  Spring Boot
  Minimal Go
```

### Generate a Minimal Java Project

![minimal-java.gif](docs%2Fminimal-java.gif)

### Generate a Spring Boot Project

![spring-boot.gif](docs%2Fspring-boot.gif)

### Generate a Go Project

![minimal-go.gif](docs%2Fminimal-go.gif)