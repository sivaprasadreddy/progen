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

### Windows
1. Go to https://github.com/sivaprasadreddy/progen/releases.
2. Download the latest file for Windows.
3. Extract the downloaded file.
4. Open the extracted folder in File Explorer.
5. Hold down the Shift key, right-click within the folder, and select "Open command window here" or "Open PowerShell window here."
6. In the command prompt or PowerShell window that appears, type progen.exe and press Enter.
```
C:\Users\appUser\Documents\workspace-sts\git\progen_Windows_x86_64>progen.exe
```
These steps should help you download, extract, and run 'progen.exe' from the command prompt in a Windows environment.

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

