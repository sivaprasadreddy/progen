# Project Generator CLI
A CLI tool for generating projects with sensible defaults.

## Motivation
There are few good project generator tools exist such as [JHipster](https://www.jhipster.tech/), [Bootify](https://bootify.io/), [generator-springboot](https://github.com/sivaprasadreddy/generator-springboot).

The reasons for creating this new **progen** CLI tool are:
* To generate the application with minimal code and the most commonly used features and configurations only.
* We should be able to use the tool offline to create the projects.
* We should be able to use the tool without requiring the installation of other languages SDKs.
  For example, for generating a Java application, we shouldn't need to install Node.js or Python, etc.

## Installation

### Using OS specific binary
Download the latest binary for your OS and Architecture(arm64, x86_64) 
from https://github.com/sivaprasadreddy/progen/releases.

For Linux and macOS, you can follow the below steps:

```shell
$ cd $HOME/Downloads
# For MacOS ARM64 (M1, M2 chips)
$ FILE="progen_Darwin_arm64.tar.gz"
# For MacOS AMD64
$ FILE="progen_Darwin_x86_64.tar.gz"
# For Linux ARM64
$ FILE="progen_Linux_arm64.tar.gz"
# For Linux AMD64
$ FILE="progen_Linux_x86_64.tar.gz"
$ curl --location --progress-bar "https://github.com/sivaprasadreddy/progen/releases/download/v0.0.2/${FILE}" > "progen.tar.gz"
$ tar -xf progen.tar.gz progen
$ ./progen --help
# move progen binary to somewhere in your PATH so that you can access progen from anywhere 
```

For Windows, download the latest Zip, extract `progen.exe` and use it as follows:

```shell
> progen.exe --help
```

### Using Go Install
If you have Go installed on your machine, you can use the following command to install `progen`:

```shell
$ go install github.com/sivaprasadreddy/progen@latest
$ progen --help
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


## License
The **progen** is an Open Source software
released under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html).
