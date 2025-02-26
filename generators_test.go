package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestGenerateSpringBootMavenRestApiWithPermutations(t *testing.T) {
	var options = []struct {
		dbType        string
		migrationTool string
		features      []string
	}{
		{"MySQL", "Flyway", []string{}},
		{"PostgreSQL", "Flyway", []string{}},
		{"MariaDB", "Flyway", []string{}},

		{"MySQL", "Liquibase", []string{}},
		{"PostgreSQL", "Liquibase", []string{}},
		{"MariaDB", "Liquibase", []string{}},

		{"MySQL", "Flyway", []string{"JWT Security"}},
		{"MariaDB", "Flyway", []string{"JWT Security"}},
		{"PostgreSQL", "Flyway", []string{"JWT Security"}},

		{"MySQL", "Liquibase", []string{"JWT Security"}},
		{"MariaDB", "Liquibase", []string{"JWT Security"}},
		{"PostgreSQL", "Liquibase", []string{"JWT Security"}},
	}

	for _, tt := range options {
		t.Run(tt.dbType+"-"+tt.migrationTool, func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "my-spring-boot-mvn-api-" + strings.ToLower(tt.dbType) + "-" + strings.ToLower(tt.migrationTool)

			pc := springboot.ProjectConfig{
				AppType:         "REST API",
				AppName:         appName,
				GroupID:         "com.sivalabs",
				ArtifactID:      appName,
				AppVersion:      "1.0",
				BasePackage:     "com.sivalabs.myapp",
				BuildTool:       "Maven",
				DbType:          tt.dbType,
				DbMigrationTool: tt.migrationTool,
				Features:        tt.features,
			}
			err := springboot.GenerateProject(pc)
			assert.Nil(t, err)
			err = testGeneratedProject(appName, mvnExec, "test")
			assert.Nil(t, err)

			//cleanup
			err = deleteDir(appName)
			assert.Nil(t, err)
		})
	}
}

func TestGenerateSpringBootWithAllFeatures(t *testing.T) {
	var options = []struct {
		appType       string
		buildTool     string
		dbType        string
		migrationTool string
		features      []string
	}{
		{"Web App", "Maven", "MySQL", "Flyway", []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "Thymeleaf", "HTMX", "Security"}},
		{"Web App", "Gradle", "PostgreSQL", "Liquibase", []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "Thymeleaf", "HTMX", "Security"}},

		{"REST API", "Maven", "MariaDB", "Flyway", []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "JWT Security"}},
		{"REST API", "Gradle", "PostgreSQL", "Liquibase", []string{"Docker Compose", "Spring Modulith", "Spring Cloud AWS", "JWT Security"}},
	}

	for _, tt := range options {
		t.Run(tt.appType+"-"+tt.buildTool+"-"+tt.dbType+"-"+tt.migrationTool, func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "spring-boot-demo-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			pc := springboot.ProjectConfig{
				AppType:         tt.appType,
				AppName:         appName,
				GroupID:         "com.sivalabs",
				ArtifactID:      appName,
				AppVersion:      "1.0",
				BasePackage:     "com.sivalabs.myapp",
				BuildTool:       tt.buildTool,
				DbType:          tt.dbType,
				DbMigrationTool: tt.migrationTool,
				Features:        tt.features,
			}
			err := springboot.GenerateProject(pc)
			assert.Nil(t, err)
			if tt.buildTool == "Maven" {
				err = testGeneratedProject(appName, mvnExec, "test")
			} else {
				err = testGeneratedProject(appName, gradleExec, "build")
			}
			assert.Nil(t, err)

			//cleanup
			err = deleteDir(appName)
			assert.Nil(t, err)
		})
	}
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
