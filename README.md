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

```shell
$ progen
? Choose application type: Minimal Java
? What is the base name of your application? myapp
? What is your application groupId? com.mycompany
? What is your application artifactId? myapp
? What is your application version? 1.0.0-SNAPSHOT
? What is your application base package? com.mycompany.myapp
? Which build tool would you like to use? Maven
Project generated successfully
```

### Generate a Spring Boot Project
```shell
$ progen
? Choose application type: Spring Boot
? What is the base name of your application? myapp
? What is your application groupId? com.mycompany
? What is your application artifactId? myapp
? What is your application version? 1.0.0-SNAPSHOT
? What is your application base package? com.mycompany.myapp
? Which build tool would you like to use? Maven
? Which database would you like to use? postgresql
? Which database migration tool would you like to use? Flyway
Project generated successfully
```

### Generate a Go Project
```shell
$ progen
? Choose application type: Minimal Go
? What is ApplicationName? myapp
? What is Module Path? github.com/username/myapp
? Choose Routing Library: Gin
Project generated successfully
```
