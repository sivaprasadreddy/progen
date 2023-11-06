package minimal_java

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
	AppName     string
	GroupID     string
	ArtifactID  string
	AppVersion  string
	BasePackage string
	BuildTool   string
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
	if err := helpers.RecreateDir(pc.AppName); err != nil {
		return err
	}
	if pc.BuildTool == "Maven" {
		if err := pg.createMavenWrapper(pc); err != nil {
			return err
		}
		if err := pg.createMavenBuildFiles(pc); err != nil {
			return err
		}
	} else {
		if err := pg.createGradleWrapper(pc); err != nil {
			return err
		}
		if err := pg.createGradleBuildFiles(pc); err != nil {
			return err
		}
	}
	if err := pg.createGitIgnore(pc); err != nil {
		return err
	}

	if err := pg.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := pg.createSrcTestJava(pc); err != nil {
		return err
	}
	return nil
}

/** Maven Functions **/

func (pg projectGenerator) createMavenBuildFiles(pc ProjectConfig) error {
	return pg.executeTemplate(pc, "pom.xml.tmpl", "pom.xml")
}

func (pg projectGenerator) createGitIgnore(pc ProjectConfig) error {
	return pg.executeTemplate(pc, "gitignore.tmpl", ".gitignore")
}

func (pg projectGenerator) createMavenWrapper(pc ProjectConfig) error {
	return pg.copyTemplateDir(pc, "maven-wrapper", "")
}

/** Gradle Functions **/

func (pg projectGenerator) createGradleWrapper(pc ProjectConfig) error {
	return pg.copyTemplateDir(pc, "gradle-wrapper", "")
}

func (pg projectGenerator) createGradleBuildFiles(pc ProjectConfig) error {
	templateMap := map[string]string{
		"build.gradle.tmpl":    "build.gradle",
		"settings.gradle.tmpl": "settings.gradle",
	}
	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

/** Common Functions **/

func (pg projectGenerator) createSrcMainJava(pc ProjectConfig) error {
	var srcMainJavaPath = "src/main/java/"
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"App.java.tmpl":  "App.java",
		"Main.java.tmpl": "Main.java",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createSrcTestJava(pc ProjectConfig) error {
	var srcTestJavaPath = "src/test/java/"
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"AppTest.java.tmpl": "AppTest.java",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
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
