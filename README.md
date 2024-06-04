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
