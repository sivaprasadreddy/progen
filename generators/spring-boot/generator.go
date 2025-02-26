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
	"strings"
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

func (p ProjectConfig) Enabled(feature string) bool {
	return p.Features != nil && slices.Contains(p.Features, feature)
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
	pc.DockerComposeSupport = pc.Enabled("Docker Compose")
	pc.SpringModulithSupport = pc.Enabled("Spring Modulith")
	pc.SpringCloudAWSSupport = pc.Enabled("Spring Cloud AWS")
	pc.ThymeleafSupport = pc.Enabled("Thymeleaf")
	pc.HTMXSupport = pc.Enabled("HTMX")
	pc.SecuritySupport = pc.Enabled("Security")
	pc.JwtSecuritySupport = pc.Enabled("JWT Security")
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
	if err := pg.createSrcMainResources(pc); err != nil {
		return err
	}
	if err := pg.createSrcTestJava(pc); err != nil {
		return err
	}
	if err := pg.createSrcTestResources(pc); err != nil {
		return err
	}
	if err := pg.createComposeConfigFiles(pc); err != nil {
		return err
	}
	if err := pg.createCIConfigFiles(pc); err != nil {
		return err
	}
	if err := pg.createReadMeFile(pc); err != nil {
		return err
	}
	if err := pg.formatCode(pc); err != nil {
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
		"Application.java.tmpl":           "Application.java",
		"ApplicationProperties.java.tmpl": "ApplicationProperties.java",
		"WebMvcConfig.java.tmpl":          "config/WebMvcConfig.java",
	}

	if pc.SecuritySupport || pc.JwtSecuritySupport {
		templateMap["SecurityConfig.java.tmpl"] = "config/SecurityConfig.java"
		templateMap["SecurityUserDetailsService.java.tmpl"] = "security/SecurityUserDetailsService.java"
	}

	if pc.SecuritySupport {
		templateMap["WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
	}

	if pc.JwtSecuritySupport {
		templateMap["JwtWebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["AuthToken.java.tmpl"] = "security/AuthToken.java"
		templateMap["TokenHelper.java.tmpl"] = "security/TokenHelper.java"
		templateMap["TokenAuthenticationFilter.java.tmpl"] = "security/TokenAuthenticationFilter.java"
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
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

	if pc.ThymeleafSupport {
		templateMap["static/css/styles.css"] = "static/css/styles.css"
		templateMap["templates/index.html.tmpl"] = "templates/index.html"
		templateMap["templates/layout.html.tmpl"] = "templates/layout.html"
	}

	if pc.ThymeleafSupport && (pc.SecuritySupport || pc.JwtSecuritySupport) {
		templateMap["templates/login.html.tmpl"] = "templates/login.html"
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, srcMainResourcesPath+tmpl, srcMainResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	if pc.DbMigrationTool == "Flyway" {
		err := pg.copyTemplateDir(pc, "src/main/resources/db/migration/flyway", "src/main/resources/db/migration")
		if err != nil {
			return err
		}
	} else {
		err := pg.copyTemplateDir(pc, "src/main/resources/db/migration/liquibase", "src/main/resources/db/migration")
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
		"ApplicationTests.java.tmpl":     "ApplicationTests.java",
		"TestcontainersConfig.java.tmpl": "TestcontainersConfig.java",
		"BaseIntegrationTest.java.tmpl":  "BaseIntegrationTest.java",
		"TestApplication.java.tmpl":      "TestApplication.java",
	}

	if pc.Enabled("Spring Modulith") {
		templateMap["ModularityTests.java.tmpl"] = "ModularityTests.java"
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
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
		err := pg.executeTemplate(pc, srcTestResourcesPath+tmpl, srcTestResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createCIConfigFiles(pc ProjectConfig) error {
	templateMap := map[string]string{
		"ci/github/workflows/ci.yml.tmpl": ".github/workflows/ci.yml",
	}

	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createComposeConfigFiles(pc ProjectConfig) error {
	if !pc.DockerComposeSupport {
		return nil
	}
	templateMap := map[string]string{
		"compose.yml.tmpl": "compose.yml",
	}
	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) createReadMeFile(pc ProjectConfig) error {
	templateMap := map[string]string{
		"README.md.tmpl": "README.md",
	}
	for tmpl, filePath := range templateMap {
		err := pg.executeTemplate(pc, tmpl, filePath)
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

func (pg projectGenerator) formatCode(pc ProjectConfig) error {
	var hostOS = runtime.GOOS
	dirName := pc.AppName
	executable := "./mvnw"
	formatCmd := "spotless:apply"
	if pc.BuildTool == "Gradle" {
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
