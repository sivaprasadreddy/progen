package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	sb "github.com/sivaprasadreddy/progen/generators/springboot"
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

func TestGenerateSpringBootWithAllFeatures(t *testing.T) {
	var options = []struct {
		appType               sb.AppType
		buildTool             sb.BuildTool
		dbType                sb.DatabaseType
		migrationTool         sb.DbMigrationTool
		SpringCloudAWSSupport bool
		ThymeleafSupport      bool
		HTMXSupport           bool
	}{
		{sb.WebApp, sb.Maven, sb.MySQL, sb.Flyway, true, true, true},
		{sb.WebApp, sb.Gradle, sb.PostgreSQL, sb.Liquibase, true, true, true},
		{sb.RestApi, sb.Maven, sb.MariaDB, sb.Flyway, true, false, false},
		{sb.RestApi, sb.Gradle, sb.PostgreSQL, sb.Liquibase, true, false, false},
	}

	for _, tt := range options {
		t.Run(tt.appType.String()+"-"+tt.buildTool.String()+"-"+tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "springboot-demo-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			pc := sb.ProjectConfig{
				AppType:               tt.appType,
				AppName:               appName,
				GroupID:               "com.sivalabs",
				ArtifactID:            appName,
				AppVersion:            "1.0",
				BasePackage:           "com.sivalabs.myapp",
				BuildTool:             tt.buildTool,
				DbType:                tt.dbType,
				DbMigrationTool:       tt.migrationTool,
				SpringCloudAWSSupport: tt.SpringCloudAWSSupport,
				ThymeleafSupport:      tt.ThymeleafSupport,
				HTMXSupport:           tt.HTMXSupport,
			}
			err := sb.GenerateProject(pc)
			assert.Nil(t, err)
			if tt.buildTool == sb.Maven {
				err = testGeneratedProject(appName, mvnExec, "test")
			} else {
				err = testGeneratedProject(appName, gradleExec, "build")
			}
			assert.Nil(t, err)

			//cleanup
			//err = deleteDir(appName)
			//assert.Nil(t, err)
			if err == nil {
				err = deleteDir(appName)
			}
		})
	}
}

func TestGenerateSpringBootMavenRestApiWithPermutations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping all combination tests in short mode")
	}
	var options = []struct {
		dbType        sb.DatabaseType
		migrationTool sb.DbMigrationTool
	}{
		{sb.MySQL, sb.Flyway},
		{sb.PostgreSQL, sb.Flyway},
		{sb.MariaDB, sb.Flyway},

		{sb.MySQL, sb.Liquibase},
		{sb.PostgreSQL, sb.Liquibase},
		{sb.MariaDB, sb.Liquibase},
	}

	for _, tt := range options {
		t.Run(tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "my-springboot-mvn-api-" + strings.ToLower(tt.dbType.String()) + "-" + strings.ToLower(tt.migrationTool.String())

			pc := sb.ProjectConfig{
				AppType:         sb.RestApi,
				AppName:         appName,
				GroupID:         "com.sivalabs",
				ArtifactID:      appName,
				AppVersion:      "1.0",
				BasePackage:     "com.sivalabs.myapp",
				BuildTool:       sb.Maven,
				DbType:          tt.dbType,
				DbMigrationTool: tt.migrationTool,
			}
			err := sb.GenerateProject(pc)
			assert.Nil(t, err)
			err = testGeneratedProject(appName, mvnExec, "test")
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deleteDir(dirName string) error {
	cmd := exec.Command("/bin/sh", "-c", "rm -rf "+dirName)
	if hostOS == "windows" {
		cmd = exec.Command("cmd", "/C", "rd /s /q "+dirName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
