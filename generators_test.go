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

func TestGenerateSpringBootMavenApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppName:         "my-spring-boot-mvn-app",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-mvn-app",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Maven",
		DbType:          "PostgreSQL",
		DbMigrationTool: "flyway",
		Features:        []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "Thymeleaf", "HTMX", "Security", "JWT Security"},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-mvn-app", mvnExec, "test")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-mvn-app")
	assert.Nil(t, err)
}

func TestGenerateSpringBootGradleApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppName:         "my-spring-boot-gradle-app",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-gradle-app",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Gradle",
		DbType:          "PostgreSQL",
		DbMigrationTool: "flyway",
		Features:        []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "Thymeleaf", "HTMX", "Security", "JWT Security"},
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-gradle-app", gradleExec, "build")
	fmt.Println("Error:", err)
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-gradle-app")
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
