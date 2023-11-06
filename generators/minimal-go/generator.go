package minimal_go

import (
	"embed"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/sivaprasadreddy/progen/generators/helpers"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

type ProjectConfig struct {
	ApplicationName string
	ModulePath      string
	RoutingLibrary  string
}

func Run() {
	projectConfig, err := getAnswers()
	helpers.FatalIfErr(err)
	err = GenerateProject(projectConfig)
	helpers.FatalIfErrOrMsg(err, "Project generated successfully")
}

func GenerateProject(pc ProjectConfig) error {
	pg := projectGenerator{tmplFS: tmplsFS}
	return pg.generate(pc)
}

type projectGenerator struct {
	tmplFS embed.FS
}

func (pg projectGenerator) generate(pc ProjectConfig) error {
	if err := helpers.RecreateDir(pc.ApplicationName); err != nil {
		return err
	}
	if err := pg.createGitIgnore(pc); err != nil {
		return err
	}
	if err := pg.createModuleFiles(pc); err != nil {
		return err
	}
	if err := pg.createMainSrc(pc); err != nil {
		return err
	}
	return nil
}

func (pg projectGenerator) createGitIgnore(pc ProjectConfig) error {
	return pg.executeTemplate(pc, templatePath("gitignore.tmpl"), ".gitignore")
}

func (pg projectGenerator) createModuleFiles(pc ProjectConfig) error {
	prefix := strings.ToLower(pc.RoutingLibrary)
	templateMap := map[string]string{
		"go.mod.tmpl": "go.mod",
		"go.sum.tmpl": "go.sum",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(prefix+"/"+tmpl), filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createMainSrc(pc ProjectConfig) error {
	prefix := strings.ToLower(pc.RoutingLibrary)
	templateMap := map[string]string{
		"config.go.tmpl": "config.go",
		"app.go.tmpl":    "app.go",
		"main.go.tmpl":   "main.go",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(prefix+"/"+tmpl), filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func templatePath(filepath string) string {
	return fmt.Sprintf("%s/%s", templatesRootDir, filepath)
}

func (pg projectGenerator) executeTemplate(pc ProjectConfig, templatePath, targetFilePath string) error {
	tmplFileContent, err := pg.tmplFS.ReadFile(templatePath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(templatePath).Parse(string(tmplFileContent))
	if err != nil {
		return err
	}
	f := helpers.CreateFile(path.Join(".", pc.ApplicationName, targetFilePath))
	err = tmpl.Execute(f, pc)
	if err != nil {
		return err
	}
	return nil
}
