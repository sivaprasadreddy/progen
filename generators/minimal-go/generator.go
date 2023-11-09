package minimal_go

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/sivaprasadreddy/progen/generators/helpers"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

type ProjectConfig struct {
	AppName        string
	ModulePath     string
	RoutingLibrary string
}

func Run() {
	projectConfig, err := getAnswers()
	helpers.FatalIfErr(err)
	err = GenerateProject(projectConfig)
	if err == nil {
		file, err := json.MarshalIndent(projectConfig, "", " ")
		if err != nil {
			fmt.Println("failed to marshall projectConfig")
		} else {
			if err = os.WriteFile(projectConfig.AppName+"/.progen.json", file, 0644); err != nil {
				fmt.Println("failed to write .progrn.json file")
			}
		}
	}
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
	if err := helpers.RecreateDir(pc.AppName); err != nil {
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
	return pg.executeTemplate(pc, "gitignore.tmpl", ".gitignore")
}

func (pg projectGenerator) createModuleFiles(pc ProjectConfig) error {
	prefix := strings.ToLower(pc.RoutingLibrary)
	templateMap := map[string]string{
		"go.mod.tmpl": "go.mod",
		"go.sum.tmpl": "go.sum",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, prefix+"/"+tmpl, filePath)
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
		err := pg.executeTemplate(pc, prefix+"/"+tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) copyTemplateDir(pc ProjectConfig, origin, dirName string) error {
	templateDirPath := fmt.Sprintf("%s/%s", templatesRootDir, origin)
	return helpers.CopyDir(pg.tmplFS, templateDirPath, pc.AppName, dirName)
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
