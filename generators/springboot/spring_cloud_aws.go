package springboot

import "strings"

type SpringCloudAwsConfig struct {
	pg projectGenerator
}

func NewSpringCloudAwsConfig(pg projectGenerator) *SpringCloudAwsConfig {
	return &SpringCloudAwsConfig{pg: pg}
}

func (b SpringCloudAwsConfig) generate(pc ProjectConfig) error {
	if !pc.SpringCloudAWSSupport {
		return nil
	}
	if err := b.createSrcTestJava(pc); err != nil {
		return err
	}
	if err := b.createSrcTestResources(pc); err != nil {
		return err
	}
	return nil
}

func (b SpringCloudAwsConfig) createSrcTestJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"SpringCloudAwsConfigTest.java.tmpl": "SpringCloudAwsConfigTest.java",
	}

	for tmpl, filePath := range templateMap {
		err := b.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b SpringCloudAwsConfig) createSrcTestResources(pc ProjectConfig) error {
	return b.pg.copyTemplateDir(pc, "src/test/resources/floci-init", "src/test/resources/floci-init")
}
