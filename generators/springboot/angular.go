package springboot

import "strings"

type AngularConfig struct {
	pg projectGenerator
}

func NewAngularConfig(pg projectGenerator) *AngularConfig {
	return &AngularConfig{pg: pg}
}

func (b AngularConfig) generate(pc ProjectConfig) error {
	if pc.AppType != SpringBootAngularFullStack {
		return nil
	}
	if err := b.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := b.copyFrontend(pc); err != nil {
		return err
	}
	return nil
}

func (b AngularConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"ui/AngularForwardController.java.tmpl": "ui/AngularForwardController.java",
	}

	for tmpl, filePath := range templateMap {
		err := b.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b AngularConfig) copyFrontend(pc ProjectConfig) error {
	return b.pg.copyTemplateDir(pc, "angular-frontend", "frontend")
}
