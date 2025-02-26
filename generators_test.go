package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
	"github.com/stretchr/testify/assert"
)

var hostOS = runtime.GOOS
var mvnExec = "./mvnw"
var gradleExec = "./gradlew"

func init() {
	fmt.Println("Host OS:", hostOS)
	if hostOS == "windows" {
		mvnExec = "mvnw.cmd"
		gradleExec = "gradlew.bat"
	}
}

func TestGenerateSpringBootMavenWebApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppType:         "Web App",
		AppName:         "my-spring-boot-mvn-webapp",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-mvn-webapp",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Maven",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Flyway",
		Features: []string{
			"Docker Compose",
			"Spring Modulith",
			"Spring Cloud AWS",
			"Thymeleaf",
			"HTMX",
			"Security",
		},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-mvn-webapp", mvnExec, "test")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-mvn-webapp")
	assert.Nil(t, err)
}

func TestGenerateSpringBootMavenRestApi(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppType:         "REST API",
		AppName:         "my-spring-boot-mvn-api",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-mvn-api",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Maven",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Liquibase",
		Features: []string{
			"Docker Compose",
			"Spring Modulith",
			"Spring Cloud AWS",
			"JWT Security",
		},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-mvn-api", mvnExec, "test")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-mvn-api")
	assert.Nil(t, err)
}

func TestGenerateSpringBootGradleWebApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppType:         "Web App",
		AppName:         "my-spring-boot-gradle-webapp",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-gradle-webapp",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Gradle",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Flyway",
		Features: []string{
			"Docker Compose",
			"Spring Modulith",
			"Spring Cloud AWS",
			"Thymeleaf",
			"HTMX",
			"Security",
		},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-gradle-webapp", gradleExec, "build")
	fmt.Println("Error:", err)
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-gradle-webapp")
	assert.Nil(t, err)
}

func TestGenerateSpringBootGradleRestApi(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppType:         "REST API",
		AppName:         "my-spring-boot-gradle-api",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-gradle-api",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Gradle",
		DbType:          "PostgreSQL",
		DbMigrationTool: "Liquibase",
		Features: []string{
			"Docker Compose",
			"Spring Modulith",
			"Spring Cloud AWS",
			"JWT Security",
		},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-gradle-api", gradleExec, "build")
	fmt.Println("Error:", err)
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-gradle-api")
	assert.Nil(t, err)
}

func testGeneratedProject(dirName, executable, testCmd string) error {
	appTestCmd := fmt.Sprintf("cd %s; %s %s;", dirName, executable, testCmd)
	cmd := exec.Command("/bin/sh", "-c", appTestCmd)
	if hostOS == "windows" {
		appTestCmd = fmt.Sprintf("cd %s && %s %s", dirName, executable, testCmd)
		cmd = exec.Command("cmd", "/C", appTestCmd)
	}
	//fmt.Println("appTestCmd: ", appTestCmd)
	out, err := cmd.CombinedOutput()
	fmt.Println("Error:", err)
	fmt.Println("Output:", string(out))
	return err
}

func deleteDir(dirName string) error {
	cmd := exec.Command("/bin/sh", "-c", "rm -rf "+dirName)
	if hostOS == "windows" {
		cmd = exec.Command("cmd", "/C", "rd /s /q "+dirName)
	}
	_, err := cmd.CombinedOutput()
	//fmt.Println("Error:", err)
	//fmt.Println("Output:", string(out))
	return err
}
