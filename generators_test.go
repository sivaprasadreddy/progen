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
	"github.com/stretchr/testify/require"
)

var hostOS = runtime.GOOS
var mvnExec = "./mvnw"
var gradleExec = "./gradlew"

// Running test/build command runs all tests by spinning up a lot of containers.
// This consumes a lot of free GitHub Actions minutes.
// For now, I am manually running all tests locally and only running compile commands in the CI pipeline.
// var mvnTestCmd = "verify"
// var gradleTestCmd = "build"
var mvnTestCmd = "test-compile"
var gradleTestCmd = "compileTestJava"

func init() {
	fmt.Println("Host OS:", hostOS)
	if hostOS == "windows" {
		mvnExec = "mvnw.cmd"
		gradleExec = "gradlew.bat"
	}
}

func TestGenerateSpringBootWithAllFeatures(t *testing.T) {
	var options = []struct {
		frontend              sb.Frontend
		buildTool             sb.BuildTool
		persistenceType       sb.PersistenceType
		dbType                sb.DatabaseType
		migrationTool         sb.DbMigrationTool
		SpringCloudAWSSupport bool
		HTMXSupport           bool
		EmailSupport          bool
		RabbitMQSupport       bool
		RedisCachingSupport   bool
		OpenTelemetrySupport  bool
		K8sSupport            bool
		SpringAIModels        []sb.AIModelProvider
	}{
		{sb.FrontendThymeleaf, sb.Maven, sb.SpringDataJPA, sb.MySQL, sb.Flyway, true, true, true, true, true, true, true, []sb.AIModelProvider{sb.AIModelOllama, sb.AIModelOpenAI, sb.AIModelAnthropic}},
		{sb.FrontendThymeleaf, sb.Gradle, sb.SpringJdbcClient, sb.PostgreSQL, sb.Liquibase, true, true, false, false, false, true, true, []sb.AIModelProvider{sb.AIModelOpenAI}},
		{sb.FrontendJTE, sb.Maven, sb.SpringDataJPA, sb.PostgreSQL, sb.Flyway, false, true, false, false, false, false, false, []sb.AIModelProvider{sb.AIModelAnthropic}},
		{sb.FrontendJTE, sb.Gradle, sb.SpringJOOQ, sb.MySQL, sb.Liquibase, false, false, false, false, false, false, false, []sb.AIModelProvider{sb.AIModelOllama}},
		{sb.FrontendNone, sb.Maven, sb.SpringJdbcClient, sb.MariaDB, sb.Flyway, true, false, false, false, false, false, false, nil},
		{sb.FrontendNone, sb.Gradle, sb.SpringDataJPA, sb.PostgreSQL, sb.Liquibase, true, false, true, false, false, false, false, nil},
		{sb.FrontendAngular, sb.Maven, sb.SpringJdbcClient, sb.PostgreSQL, sb.Flyway, false, false, false, false, false, false, false, nil},
		{sb.FrontendAngular, sb.Gradle, sb.SpringDataJPA, sb.PostgreSQL, sb.Flyway, false, false, false, false, false, false, false, nil},
		{sb.FrontendNone, sb.Maven, sb.SpringJOOQ, sb.PostgreSQL, sb.Flyway, false, false, false, false, false, false, false, nil},
		{sb.FrontendNone, sb.Gradle, sb.SpringJOOQ, sb.MySQL, sb.Liquibase, false, false, false, false, false, false, false, nil},
	}

	for _, tt := range options {
		passed := t.Run(tt.frontend.String()+"-"+tt.buildTool.String()+"-"+tt.persistenceType.String()+"-"+tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "springboot-demo-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			pc := sb.ProjectConfig{
				Frontend:              tt.frontend,
				AppName:               appName,
				GroupID:               "com.sivalabs",
				ArtifactID:            appName,
				AppVersion:            "1.0",
				BasePackage:           "com.sivalabs.myapp",
				BuildTool:             tt.buildTool,
				PersistenceType:       tt.persistenceType,
				DbType:                tt.dbType,
				DbMigrationTool:       tt.migrationTool,
				SpringCloudAWSSupport: tt.SpringCloudAWSSupport,
				HTMXSupport:           tt.HTMXSupport,
				EmailSupport:          tt.EmailSupport,
				RabbitMQSupport:       tt.RabbitMQSupport,
				RedisCachingSupport:   tt.RedisCachingSupport,
				OpenTelemetrySupport:  tt.OpenTelemetrySupport,
				K8sSupport:            tt.K8sSupport,
				SpringAIModels:        tt.SpringAIModels,
			}
			err := sb.GenerateProject(pc)
			require.NoError(t, err)
			if tt.buildTool == sb.Maven {
				err = testGeneratedProject(appName, mvnExec, mvnTestCmd)
			} else {
				err = testGeneratedProject(appName, gradleExec, gradleTestCmd)
			}
			require.NoError(t, err)

			if err == nil {
				err = deleteDir(appName)
			}
		})
		if !passed {
			t.FailNow()
		}
	}
}

func TestGenerateSpringBootMavenRestApiWithPermutations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping all combination tests in short mode")
	}
	var options = []struct {
		persistenceType sb.PersistenceType
		dbType          sb.DatabaseType
		migrationTool   sb.DbMigrationTool
	}{
		{sb.SpringDataJPA, sb.MySQL, sb.Flyway},
		{sb.SpringDataJPA, sb.PostgreSQL, sb.Flyway},
		{sb.SpringDataJPA, sb.MariaDB, sb.Flyway},

		{sb.SpringJdbcClient, sb.MySQL, sb.Flyway},
		{sb.SpringJdbcClient, sb.PostgreSQL, sb.Flyway},
		{sb.SpringJdbcClient, sb.MariaDB, sb.Flyway},

		{sb.SpringDataJPA, sb.MySQL, sb.Liquibase},
		{sb.SpringDataJPA, sb.PostgreSQL, sb.Liquibase},
		{sb.SpringJdbcClient, sb.MariaDB, sb.Liquibase},

		{sb.SpringJOOQ, sb.MySQL, sb.Flyway},
		{sb.SpringJOOQ, sb.PostgreSQL, sb.Flyway},
		{sb.SpringJOOQ, sb.MariaDB, sb.Flyway},
		{sb.SpringJOOQ, sb.PostgreSQL, sb.Liquibase},
	}

	for _, tt := range options {
		passed := t.Run(tt.dbType.String()+"-"+tt.migrationTool.String(), func(t *testing.T) {
			t.Log("Generating App with Options: ", tt)
			appName := "my-springboot-mvn-api-" + strings.ToLower(tt.dbType.String()) + "-" + strings.ToLower(tt.migrationTool.String())

			pc := sb.ProjectConfig{
				Frontend:        sb.FrontendNone,
				AppName:         appName,
				GroupID:         "com.sivalabs",
				ArtifactID:      appName,
				AppVersion:      "1.0",
				BasePackage:     "com.sivalabs.myapp",
				BuildTool:       sb.Maven,
				PersistenceType: tt.persistenceType,
				DbType:          tt.dbType,
				DbMigrationTool: tt.migrationTool,
				EmailSupport:    true,
				K8sSupport:      true,
			}
			err := sb.GenerateProject(pc)
			require.NoError(t, err)
			err = testGeneratedProject(appName, mvnExec, "test")
			require.NoError(t, err)

			//cleanup
			//err = deleteDir(appName)
			//assert.Nil(t, err)
			if err == nil {
				err = deleteDir(appName)
			}
		})
		if !passed {
			t.FailNow()
		}
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
