package spring_boot

import "strings"

type SpringModulithConfig struct {
	pg projectGenerator
}

func NewSpringModulithConfig(pg projectGenerator) *SpringModulithConfig {
	return &SpringModulithConfig{pg: pg}
}

func (s SpringModulithConfig) generate(pc ProjectConfig) error {
	if err := s.createSrcTestJava(pc); err != nil {
		return err
	}
	return nil
}

func (s SpringModulithConfig) createSrcTestJava(pc ProjectConfig) error {
	if !pc.EnabledSpringModulithSupport() {
		return nil
	}
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{}

	templateMap["ModularityTests.java.tmpl"] = "ModularityTests.java"

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
