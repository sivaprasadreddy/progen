package spring_boot

import (
	"embed"
	"fmt"
	"github.com/sivaprasadreddy/progen/generators/helpers"
	"path"
	"strings"
	"text/template"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

type ProjectConfig struct {
	ApplicationName    string
	GroupID            string
	ArtifactID         string
	ApplicationVersion string
	BasePackage        string
	BuildTool          string
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
	if err := pg.createSrcMainResources(pc); err != nil {
		return err
	}
	if err := pg.createSrcTestJava(pc); err != nil {
		return err
	}
	if err := pg.createSrcTestResources(pc); err != nil {
		return err
	}
	return nil
}

/** Maven Functions **/

func (pg projectGenerator) createMavenBuildFiles(pc ProjectConfig) error {
	return pg.executeTemplate(pc, templatePath("pom.xml.tmpl"), "pom.xml")
}

func (pg projectGenerator) createGitIgnore(pc ProjectConfig) error {
	return pg.executeTemplate(pc, templatePath("gitignore.tmpl"), ".gitignore")
}

func (pg projectGenerator) createMavenWrapper(pc ProjectConfig) error {
	return helpers.CopyDir(pg.tmplFS, templatePath("maven-wrapper"), pc.ApplicationName, "")
}

/** Gradle Functions **/

func (pg projectGenerator) createGradleWrapper(pc ProjectConfig) error {
	return helpers.CopyDir(pg.tmplFS, templatePath("gradle-wrapper"), pc.ApplicationName, "")
}

func (pg projectGenerator) createGradleBuildFiles(pc ProjectConfig) error {
	templateMap := map[string]string{
		"build.gradle.tmpl":    "build.gradle",
		"settings.gradle.tmpl": "settings.gradle",
	}
	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(tmpl), filePath)
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
		"Application.java.tmpl": "Application.java",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(srcMainJavaPath+tmpl), srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createSrcMainResources(pc ProjectConfig) error {
	var srcMainResourcesPath = "src/main/resources/"
	templateMap := map[string]string{
		"application.properties.tmpl": "application.properties",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(srcMainResourcesPath+tmpl), srcMainResourcesPath+filePath)
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
		"ApplicationTests.java.tmpl":    "ApplicationTests.java",
		"ContainersConfig.java.tmpl":    "ContainersConfig.java",
		"BaseIntegrationTest.java.tmpl": "BaseIntegrationTest.java",
		"TestApplication.java.tmpl":     "TestApplication.java",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(srcTestJavaPath+tmpl), srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createSrcTestResources(pc ProjectConfig) error {
	var srcTestResourcesPath = "src/test/resources/"
	templateMap := map[string]string{
		"logback-test.xml.tmpl": "logback-test.xml",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, templatePath(srcTestResourcesPath+tmpl), srcTestResourcesPath+filePath)
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
