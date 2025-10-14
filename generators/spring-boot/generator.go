package spring_boot

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"slices"
	"text/template"

	"github.com/sivaprasadreddy/progen/generators/helpers"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

type ProjectConfig struct {
	AppType               string
	AppName               string
	GroupID               string
	ArtifactID            string
	AppVersion            string
	BasePackage           string
	BuildTool             string
	DbType                string
	DbMigrationTool       string
	Features              []string
	DockerComposeSupport  bool
	SpringModulithSupport bool
	SpringCloudAWSSupport bool
	ThymeleafSupport      bool
	HTMXSupport           bool
	SecuritySupport       bool
	JwtSecuritySupport    bool
}

func Run() {
	projectConfig, err := getAnswers()
	helpers.FatalIfErr(err)
	err = GenerateProject(projectConfig)
	helpers.FatalIfErrOrMsg(err, "Project generated successfully")
}

func GenerateProject(pc ProjectConfig) error {
	updateFeatureFlags(&pc)
	pg := projectGenerator{tmplFS: tmplsFS}
	err := pg.generate(pc)
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(pc, "", " ")
	if err != nil {
		fmt.Println("failed to marshall projectConfig")
	} else {
		if err = os.WriteFile(pc.AppName+"/.progen.json", file, 0644); err != nil {
			fmt.Println("failed to write .progrn.json file")
		}
	}
	return nil
}

func updateFeatureFlags(pc *ProjectConfig) {
	pc.DockerComposeSupport = pc.EnabledDockerComposeSupport()
	pc.SpringModulithSupport = pc.EnabledSpringModulithSupport()
	pc.SpringCloudAWSSupport = pc.EnabledSpringCloudAWSSupport()
	pc.ThymeleafSupport = pc.EnabledThymeleafSupport()
	pc.HTMXSupport = pc.EnabledHTMXSupport()
	pc.SecuritySupport = pc.EnabledSecuritySupport()
	pc.JwtSecuritySupport = pc.EnabledJwtSecuritySupport()
}

func (p ProjectConfig) EnabledSecuritySupport() bool {
	return p.Enabled(FeatureSecuritySupport)
}

func (p ProjectConfig) EnabledJwtSecuritySupport() bool {
	return p.Enabled(FeatureJwtSecuritySupport)
}

func (p ProjectConfig) EnabledSpringModulithSupport() bool {
	return p.Enabled(FeatureSpringModulithSupport)
}

func (p ProjectConfig) EnabledSpringCloudAWSSupport() bool {
	return p.Enabled(FeatureSpringCloudAWSSupport)
}

func (p ProjectConfig) EnabledThymeleafSupport() bool {
	return p.Enabled(FeatureThymeleafSupport)
}

func (p ProjectConfig) EnabledHTMXSupport() bool {
	return p.Enabled(FeatureHTMXSupport)
}

func (p ProjectConfig) EnabledDockerComposeSupport() bool {
	return p.Enabled(FeatureDockerComposeSupport)
}

func (p ProjectConfig) Enabled(feature string) bool {
	return p.Features != nil && slices.Contains(p.Features, feature)
}

type projectGenerator struct {
	tmplFS embed.FS
}

func (pg projectGenerator) generate(pc ProjectConfig) error {
	if err := helpers.RecreateDir(pc.AppName); err != nil {
		return err
	}

	if err := NewBuildToolConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewTaskfileConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewRenovateConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewGitIgnoreConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewDockerComposeConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewGhActionsConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewReadMeConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewAppCommonConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewThymeleafConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewDbMigrationsConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewSpringModulithConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := NewSecurityConfig(pg).generate(pc); err != nil {
		return err
	}
	if err := pg.formatCode(pc); err != nil {
		return err
	}
	return nil
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
	err = tmpl.Execute(f, pc)
	if err != nil {
		return err
	}
	return nil
}

func (pg projectGenerator) formatCode(pc ProjectConfig) error {
	var hostOS = runtime.GOOS
	dirName := pc.AppName
	executable := "./mvnw"
	formatCmd := "spotless:apply"
	if pc.BuildTool == BuildToolGradle {
		executable = "./gradlew"
		formatCmd = "spotlessApply"
	}
	appFormatCmd := fmt.Sprintf("cd %s; %s %s;", dirName, executable, formatCmd)
	cmd := exec.Command("/bin/sh", "-c", appFormatCmd)
	if hostOS == "windows" {
		appFormatCmd = fmt.Sprintf("cd %s && %s %s", dirName, executable, formatCmd)
		cmd = exec.Command("cmd", "/C", appFormatCmd)
	}
	//fmt.Println("appTestCmd: ", appTestCmd)
	_, err := cmd.CombinedOutput()
	//fmt.Println("Error:", err)
	//fmt.Println("Output:", string(out))
	return err
}
