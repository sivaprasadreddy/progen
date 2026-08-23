package springboot

import "strings"

type EmailConfig struct {
	pg projectGenerator
}

func NewEmailConfig(pg projectGenerator) *EmailConfig {
	return &EmailConfig{pg: pg}
}

func (s EmailConfig) generate(pc ProjectConfig) error {
	if !pc.EmailSupport {
		return nil
	}
	if err := s.createSrcMainJava(pc); err != nil {
		return err
	}
	return nil
}

func (s EmailConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"shared/email/EmailService.java.tmpl": "shared/email/EmailService.java",
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
