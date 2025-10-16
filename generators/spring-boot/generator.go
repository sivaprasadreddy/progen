// Package spring_boot provides functionality for generating Spring Boot projects
// with various configurations and features.
package spring_boot

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"text/template"

	"github.com/sivaprasadreddy/progen/generators/helpers"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

// ProjectConfig holds all configuration options for generating a Spring Boot project.
type ProjectConfig struct {
	AppType               AppType
	AppName               string
	GroupID               string
	ArtifactID            string
	AppVersion            string
	BasePackage           string
	BuildTool             BuildTool
	DbType                DatabaseType
	DbMigrationTool       DbMigrationTool
	DockerComposeSupport  bool
	SpringModulithSupport bool
	SpringCloudAWSSupport bool
	ThymeleafSupport      bool
	HTMXSupport           bool
	SecuritySupport       bool
	JwtSecuritySupport    bool
}

func Run() error {
	projectConfig, err := getAnswers()
	helpers.FatalIfErr(err)
	err = GenerateProject(projectConfig)
	helpers.FatalIfErrOrMsg(err, "Project generated successfully")
	return err
}

func GenerateProject(pc ProjectConfig) error {
	pg := projectGenerator{tmplFS: tmplsFS}
	if err := pg.generate(pc); err != nil {
		return err
	}
	return writeConfigFile(pc, pc.AppName+"/.progen.json")
}

func GenerateInitConfig() error {
	pc := ProjectConfig{
		AppType:               RestApi,
		AppName:               "myapp",
		GroupID:               "com.mycompany",
		ArtifactID:            "myapp",
		AppVersion:            "1.0.0",
		BasePackage:           "com.mycompany.myapp",
		BuildTool:             BuildToolMaven,
		DbType:                PostgreSQL,
		DbMigrationTool:       Flyway,
		DockerComposeSupport:  true,
		SpringModulithSupport: false,
		SpringCloudAWSSupport: false,
		ThymeleafSupport:      false,
		HTMXSupport:           false,
		SecuritySupport:       false,
		JwtSecuritySupport:    false,
	}
	return writeConfigFile(pc, ".progen.json")
}

func writeConfigFile(pc ProjectConfig, filePath string) error {
	file, err := json.MarshalIndent(pc, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}
	if err = os.WriteFile(filePath, file, FilePermission); err != nil {
		return fmt.Errorf("failed to write .progen.json file: %w", err)
	}
	return nil
}

type configGenerator interface {
	generate(ProjectConfig) error
}

type projectGenerator struct {
	tmplFS embed.FS
}

func (pg projectGenerator) generate(pc ProjectConfig) error {
	if err := helpers.RecreateDir(pc.AppName); err != nil {
		return err
	}

	generators := []configGenerator{
		NewBuildToolConfig(pg),
		NewSdkmanConfig(pg),
		NewTaskfileConfig(pg),
		NewRenovateConfig(pg),
		NewGitIgnoreConfig(pg),
		NewDockerComposeConfig(pg),
		NewGhActionsConfig(pg),
		NewReadMeConfig(pg),
		NewAppCommonConfig(pg),
		NewThymeleafConfig(pg),
		NewDbMigrationsConfig(pg),
		NewSpringModulithConfig(pg),
		NewSecurityConfig(pg),
	}

	for _, gen := range generators {
		if err := gen.generate(pc); err != nil {
			return err
		}
	}

	return pg.formatCode(pc)
}

func (pg projectGenerator) copyTemplateDir(pc ProjectConfig, origin, dirName string) error {
	templateDirPath := fmt.Sprintf("%s/%s", templatesRootDir, origin)
	return helpers.CopyDir(pg.tmplFS, templateDirPath, pc.AppName, dirName)
}

func (pg projectGenerator) copyTemplateFile(pc ProjectConfig, sourceFilePath, targetFilePath string) error {
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, sourceFilePath)
	return helpers.CopyTemplateFile(pg.tmplFS, templateFilePath, pc.AppName, targetFilePath)
}

func (pg projectGenerator) executeTemplate(pc ProjectConfig, templatePath, targetFilePath string) error {
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, templatePath)
	tmplFileContent, err := pg.tmplFS.ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(templateFilePath).Parse(string(tmplFileContent))
	if err != nil {
		return err
	}
	f := helpers.CreateFile(path.Join(".", pc.AppName, targetFilePath))
	defer f.Close()
	err = tmpl.Execute(f, pc)
	if err != nil {
		return err
	}
	return nil
}

func (pg projectGenerator) formatCode(pc ProjectConfig) error {
	executable, formatCmd := pg.getBuildToolCommands(pc.BuildTool)
	appFormatCmd := pg.buildCommandString(pc.AppName, executable, formatCmd)
	cmd := pg.createOSCommand(appFormatCmd)
	_, err := cmd.CombinedOutput()
	return err
}

func (pg projectGenerator) getBuildToolCommands(buildTool BuildTool) (executable, formatCmd string) {
	isWindows := runtime.GOOS == "windows"

	if buildTool == BuildToolGradle {
		if isWindows {
			return "gradlew.bat", "spotlessApply"
		}
		return "./gradlew", "spotlessApply"
	}

	if isWindows {
		return "mvnw.cmd", "spotless:apply"
	}
	return "./mvnw", "spotless:apply"
}

func (pg projectGenerator) buildCommandString(dirName, executable, formatCmd string) string {
	return fmt.Sprintf("cd %s && %s %s", dirName, executable, formatCmd)
}

func (pg projectGenerator) createOSCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("/bin/sh", "-c", command)
}
