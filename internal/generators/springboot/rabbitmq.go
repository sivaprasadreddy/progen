package springboot

import "strings"

type RabbitMQConfig struct {
	pg projectGenerator
}

func NewRabbitMQConfig(pg projectGenerator) *RabbitMQConfig {
	return &RabbitMQConfig{pg: pg}
}

func (r RabbitMQConfig) generate(pc ProjectConfig) error {
	if !pc.RabbitMQSupport {
		return nil
	}
	if err := r.createSrcMainJava(pc); err != nil {
		return err
	}
	return nil
}

func (r RabbitMQConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"config/RabbitMQConfig.java.tmpl": "config/RabbitMQConfig.java",
	}

	for tmpl, filePath := range templateMap {
		err := r.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
