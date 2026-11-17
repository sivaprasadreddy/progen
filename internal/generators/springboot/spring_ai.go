package springboot

import "strings"

type SpringAIConfig struct {
	pg projectGenerator
}

func NewSpringAIConfig(pg projectGenerator) *SpringAIConfig {
	return &SpringAIConfig{pg: pg}
}

func (s SpringAIConfig) generate(pc ProjectConfig) error {
	if !pc.SpringAISupport() {
		return nil
	}
	return s.createSrcMainJava(pc)
}

func (s SpringAIConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"shared/ai/AiChatService.java.tmpl": "shared/ai/AiChatService.java",
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
