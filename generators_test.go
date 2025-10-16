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

func TestGenerateSpringBootWithAllFeatures(t *testing.T) {
	var options = []struct {
		appType               string
		buildTool             string
		dbType                string
		migrationTool         string
		DockerComposeSupport  bool
		SpringModulithSupport bool
		SpringCloudAWSSupport bool
		ThymeleafSupport      bool
		HTMXSupport           bool
		SecuritySupport       bool
		JwtSecuritySupport    bool
	}{
		{springboot.AppTypeWebApp, springboot.BuildToolMaven, springboot.DbMySQL, springboot.DbMigrationToolFlyway, true, true, true, true, true, true, false},
		{springboot.AppTypeWebApp, springboot.BuildToolGradle, springboot.DbPostgreSQL, springboot.DbMigrationToolLiquibase, true, true, true, true, true, true, false},
		{springboot.AppTypeRestApi, springboot.BuildToolMaven, springboot.DbMariaDB, springboot.DbMigrationToolFlyway, true, true, true, false, false, false, true},
		{springboot.AppTypeRestApi, springboot.BuildToolGradle, springboot.DbPostgreSQL, springboot.DbMigrationToolLiquibase, true, true, true, false, false, false, true},
	}

	for _, tt := range options {
		t.Run(tt.appType+"-"+tt.buildTool+"-"+tt.dbType+"-"+tt.migrationTool, func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "spring-boot-demo-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			pc := springboot.ProjectConfig{
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
			err := springboot.GenerateProject(pc)
			assert.Nil(t, err)
			if tt.buildTool == springboot.BuildToolMaven {
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
		dbType             string
		migrationTool      string
		JwtSecuritySupport bool
	}{
		{springboot.DbMySQL, springboot.DbMigrationToolFlyway, false},
		{springboot.DbPostgreSQL, springboot.DbMigrationToolFlyway, false},
		{springboot.DbMariaDB, springboot.DbMigrationToolFlyway, false},

		{springboot.DbMySQL, springboot.DbMigrationToolLiquibase, false},
		{springboot.DbPostgreSQL, springboot.DbMigrationToolLiquibase, false},
		{springboot.DbMariaDB, springboot.DbMigrationToolLiquibase, false},

		{springboot.DbMySQL, springboot.DbMigrationToolFlyway, true},
		{springboot.DbMariaDB, springboot.DbMigrationToolFlyway, true},
		{springboot.DbPostgreSQL, springboot.DbMigrationToolFlyway, true},

		{springboot.DbMySQL, springboot.DbMigrationToolLiquibase, true},
		{springboot.DbMariaDB, springboot.DbMigrationToolLiquibase, true},
		{springboot.DbPostgreSQL, springboot.DbMigrationToolLiquibase, true},
	}

	for _, tt := range options {
		t.Run(tt.dbType+"-"+tt.migrationTool, func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "my-spring-boot-mvn-api-" + strings.ToLower(tt.dbType) + "-" + strings.ToLower(tt.migrationTool)

			pc := springboot.ProjectConfig{
				AppType:            springboot.AppTypeRestApi,
				AppName:            appName,
				GroupID:            "com.sivalabs",
				ArtifactID:         appName,
				AppVersion:         "1.0",
				BasePackage:        "com.sivalabs.myapp",
				BuildTool:          springboot.BuildToolMaven,
				DbType:             tt.dbType,
				DbMigrationTool:    tt.migrationTool,
				JwtSecuritySupport: tt.JwtSecuritySupport,
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
