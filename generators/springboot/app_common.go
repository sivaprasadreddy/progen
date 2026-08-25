package springboot

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
		"Application.java.tmpl":                      "Application.java",
		"ApplicationProperties.java.tmpl":            "ApplicationProperties.java",
		"config/WebMvcConfig.java.tmpl":              "config/WebMvcConfig.java",
		"shared/PagedResult.java.tmpl":               "shared/PagedResult.java",
		"shared/BadRequestException.java.tmpl":       "shared/BadRequestException.java",
		"shared/ResourceNotFoundException.java.tmpl": "shared/ResourceNotFoundException.java",
	}
	if pc.PersistenceType == SpringDataJPA {
		templateMap["config/PersistenceConfig.java.tmpl"] = "config/PersistenceConfig.java"
		templateMap["shared/JpaBaseEntity.java.tmpl"] = "shared/BaseEntity.java"
	}
	if pc.PersistenceType == SpringJdbcClient {
		templateMap["shared/JdbcClientBaseEntity.java.tmpl"] = "shared/BaseEntity.java"
	}
	if pc.PersistenceType == SpringJOOQ {
		templateMap["shared/JooqBaseEntity.java.tmpl"] = "shared/BaseEntity.java"
		templateMap["jooq/package-info.java.tmpl"] = "jooq/package-info.java"
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
		"application.properties.tmpl":       "application.properties",
		"application-local.properties.tmpl": "application-local.properties",
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
		"BaseIT.java.tmpl":               "BaseIT.java",
		"ArchUnitTests.java.tmpl":        "ArchUnitTests.java",
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
