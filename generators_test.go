package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	sb "github.com/sivaprasadreddy/progen/generators/spring-boot"
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
		DockerComposeSupport  bool
		SpringModulithSupport bool
		SpringCloudAWSSupport bool
		ThymeleafSupport      bool
		HTMXSupport           bool
		SecuritySupport       bool
		JwtSecuritySupport    bool
	}{
		{sb.WebApp, sb.BuildToolMaven, sb.MySQL, sb.Flyway, true, true, true, true, true, true, false},
		{sb.WebApp, sb.BuildToolGradle, sb.PostgreSQL, sb.Liquibase, true, true, true, true, true, true, false},
		{sb.RestApi, sb.BuildToolMaven, sb.MariaDB, sb.Flyway, true, true, true, false, false, false, true},
		{sb.RestApi, sb.BuildToolGradle, sb.PostgreSQL, sb.Liquibase, true, true, true, false, false, false, true},
	}

	for _, tt := range options {
		t.Run(tt.appType.String()+"-"+tt.buildTool.String()+"-"+tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "spring-boot-demo-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

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
				DockerComposeSupport:  tt.DockerComposeSupport,
				SpringModulithSupport: tt.SpringModulithSupport,
				SpringCloudAWSSupport: tt.SpringCloudAWSSupport,
				ThymeleafSupport:      tt.ThymeleafSupport,
				HTMXSupport:           tt.HTMXSupport,
				SecuritySupport:       tt.SecuritySupport,
				JwtSecuritySupport:    tt.JwtSecuritySupport,
			}
			err := sb.GenerateProject(pc)
			assert.Nil(t, err)
			if tt.buildTool == sb.BuildToolMaven {
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

func TestGenerateSpringBootMavenRestApiWithPermutations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping all combination tests in short mode")
	}
	var options = []struct {
		dbType             sb.DatabaseType
		migrationTool      sb.DbMigrationTool
		JwtSecuritySupport bool
	}{
		{sb.MySQL, sb.Flyway, false},
		{sb.PostgreSQL, sb.Flyway, false},
		{sb.MariaDB, sb.Flyway, false},

		{sb.MySQL, sb.Liquibase, false},
		{sb.PostgreSQL, sb.Liquibase, false},
		{sb.MariaDB, sb.Liquibase, false},

		{sb.MySQL, sb.Flyway, true},
		{sb.MariaDB, sb.Flyway, true},
		{sb.PostgreSQL, sb.Flyway, true},

		{sb.MySQL, sb.Liquibase, true},
		{sb.MariaDB, sb.Liquibase, true},
		{sb.PostgreSQL, sb.Liquibase, true},
	}

	for _, tt := range options {
		t.Run(tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "my-spring-boot-mvn-api-" + strings.ToLower(tt.dbType.String()) + "-" + strings.ToLower(tt.migrationTool.String())

			pc := sb.ProjectConfig{
				AppType:            sb.RestApi,
				AppName:            appName,
				GroupID:            "com.sivalabs",
				ArtifactID:         appName,
				AppVersion:         "1.0",
				BasePackage:        "com.sivalabs.myapp",
				BuildTool:          sb.BuildToolMaven,
				DbType:             tt.dbType,
				DbMigrationTool:    tt.migrationTool,
				JwtSecuritySupport: tt.JwtSecuritySupport,
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
