package springboot

type SpringCloudAwsConfig struct {
	pg projectGenerator
}

func NewSpringCloudAwsConfig(pg projectGenerator) *SpringCloudAwsConfig {
	return &SpringCloudAwsConfig{pg: pg}
}

func (b SpringCloudAwsConfig) generate(pc ProjectConfig) error {
	return b.copyFlociInit(pc)
}

func (b SpringCloudAwsConfig) copyFlociInit(pc ProjectConfig) error {
	return b.pg.copyTemplateDir(pc, "src/test/resources/floci-init", "src/test/resources/floci-init")
}
