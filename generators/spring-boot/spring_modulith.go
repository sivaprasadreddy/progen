package spring_boot

import "strings"

type SpringModulithConfig struct {
	pg projectGenerator
}

func NewSpringModulithConfig(pg projectGenerator) *SpringModulithConfig {
	return &SpringModulithConfig{pg: pg}
}

func (s SpringModulithConfig) generate(pc ProjectConfig) error {
	if !pc.SpringModulithSupport {
		return nil
	}
	return s.createSrcTestJava(pc)
}

func (s SpringModulithConfig) createSrcTestJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"ModularityTests.java.tmpl": "ModularityTests.java",
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
