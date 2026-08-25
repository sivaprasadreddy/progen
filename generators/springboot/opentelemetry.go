package springboot

import "strings"

type OpenTelemetryConfig struct {
	pg projectGenerator
}

func NewOpenTelemetryConfig(pg projectGenerator) *OpenTelemetryConfig {
	return &OpenTelemetryConfig{pg: pg}
}

func (a OpenTelemetryConfig) generate(pc ProjectConfig) error {
	if !pc.OpenTelemetrySupport {
		return nil
	}
	if err := a.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := a.createSrcMainResources(pc); err != nil {
		return err
	}
	return nil
}

func (a OpenTelemetryConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"config/InstallOpenTelemetryAppender.java.tmpl": "config/InstallOpenTelemetryAppender.java",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a OpenTelemetryConfig) createSrcMainResources(pc ProjectConfig) error {
	templateMap := map[string]string{
		"logback-spring.xml.tmpl": "logback-spring.xml",
	}

	for tmpl, filePath := range templateMap {
		err := a.pg.executeTemplate(pc, srcMainResourcesPath+tmpl, srcMainResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
