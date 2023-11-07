package main

import (
	"fmt"
	minimalgo "github.com/sivaprasadreddy/progen/generators/minimal-go"
	minimaljava "github.com/sivaprasadreddy/progen/generators/minimal-java"
	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"runtime"
	"testing"
)

var hostOS = runtime.GOOS
var mvnExec = "./mvnw"
var gradleExec = "./gradlew"

func init() {
	fmt.Println("Host OS:", hostOS)
	if hostOS == "windows" {
		mvnExec = "./mvnw.cmd"
		gradleExec = "./gradlew.bat"
	}
}
func TestGenerateMinimalJavaMavenApp(t *testing.T) {
	pc := minimaljava.ProjectConfig{
		AppName:     "my-minimal-java-mvn-app",
		GroupID:     "com.sivalabs",
		ArtifactID:  "my-minimal-java-mvn-app",
		AppVersion:  "1.0",
		BasePackage: "com.sivalabs.myapp",
		BuildTool:   "Maven",
	}
	err := minimaljava.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-minimal-java-mvn-app", mvnExec, "test")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-minimal-java-mvn-app")
	assert.Nil(t, err)
}

func TestGenerateMinimalJavaGradleApp(t *testing.T) {
	pc := minimaljava.ProjectConfig{
		AppName:     "my-minimal-java-gradle-app",
		GroupID:     "com.sivalabs",
		ArtifactID:  "my-minimal-java-gradle-app",
		AppVersion:  "1.0",
		BasePackage: "com.sivalabs.myapp",
		BuildTool:   "Gradle",
	}
	err := minimaljava.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-minimal-java-gradle-app", gradleExec, "build")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-minimal-java-gradle-app")
	assert.Nil(t, err)
}

func TestGenerateSpringBootMavenApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppName:         "my-spring-boot-mvn-app",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-mvn-app",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Maven",
		DbType:          "postgresql",
		DbMigrationTool: "flyway",
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
		DbType:          "postgresql",
		DbMigrationTool: "flyway",
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-spring-boot-gradle-app", gradleExec, "build")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-spring-boot-gradle-app")
	assert.Nil(t, err)
}

func TestGenerateMinimalGoGinApp(t *testing.T) {
	pc := minimalgo.ProjectConfig{
		AppName:        "my-minimal-go-gin-app",
		ModulePath:     "github.com/sivalabs/myginapp",
		RoutingLibrary: "Gin",
	}
	err := minimalgo.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-minimal-go-gin-app", "go", "build")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-minimal-go-gin-app")
	assert.Nil(t, err)
}

func TestGenerateMinimalGoChiApp(t *testing.T) {
	pc := minimalgo.ProjectConfig{
		AppName:        "my-minimal-go-chi-app",
		ModulePath:     "github.com/sivalabs/mychiapp",
		RoutingLibrary: "Chi",
	}
	err := minimalgo.GenerateProject(pc)
	assert.Nil(t, err)
	err = testGeneratedProject("my-minimal-go-chi-app", "go", "build")
	assert.Nil(t, err)

	//cleanup
	err = deleteDir("my-minimal-go-chi-app")
	assert.Nil(t, err)
}

func deleteDir(dirName string) error {
	cmd := exec.Command("/bin/sh", "-c", "rm -rf "+dirName)
	if hostOS == "windows" {
		cmd = exec.Command("cmd", "/C", "rd /s /q "+dirName)
	}
	err := cmd.Run()
	fmt.Println("Error:", err)
	return nil
}

func testGeneratedProject(dirName, executable, testCmd string) error {
	appTestCmd := fmt.Sprintf("cd %s; %s %s;", dirName, executable, testCmd)
	cmd := exec.Command("/bin/sh", "-c", appTestCmd)
	if hostOS == "windows" {
		cmd = exec.Command("cmd", "/C", appTestCmd)
	}
	fmt.Println("appTestCmd: ", appTestCmd)
	err := cmd.Run()
	fmt.Println("Error:", err)
	return nil
}
