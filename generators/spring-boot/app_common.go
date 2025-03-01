package spring_boot

import "strings"

type AppCommonConfig struct {
	pg projectGenerator
}

func NewAppCommonConfig(pg projectGenerator) *AppCommonConfig {
	return &AppCommonConfig{pg: pg}
}

func (a AppCommonConfig) generate(pc ProjectConfig) error {
	if err := a.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := a.createSrcMainResources(pc); err != nil {
		return err
	}
	if err := a.createSrcTestJava(pc); err != nil {
		return err
	}
	if err := a.createSrcTestResources(pc); err != nil {
		return err
	}
	return nil
}

func (a AppCommonConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"Application.java.tmpl":               "Application.java",
		"ApplicationProperties.java.tmpl":     "ApplicationProperties.java",
		"WebMvcConfig.java.tmpl":              "config/WebMvcConfig.java",
		"BaseEntity.java.tmpl":                "domain/BaseEntity.java",
		"BadRequestException.java.tmpl":       "domain/BadRequestException.java",
		"ResourceNotFoundException.java.tmpl": "domain/ResourceNotFoundException.java",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AppCommonConfig) createSrcMainResources(pc ProjectConfig) error {
	templateMap := map[string]string{
		"application.properties.tmpl": "application.properties",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcMainResourcesPath+tmpl, srcMainResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AppCommonConfig) createSrcTestJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"ApplicationTests.java.tmpl":     "ApplicationTests.java",
		"TestcontainersConfig.java.tmpl": "TestcontainersConfig.java",
		"BaseIntegrationTest.java.tmpl":  "BaseIntegrationTest.java",
		"TestApplication.java.tmpl":      "TestApplication.java",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AppCommonConfig) createSrcTestResources(pc ProjectConfig) error {
	templateMap := map[string]string{
		"logback-test.xml.tmpl": "logback-test.xml",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcTestResourcesPath+tmpl, srcTestResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
