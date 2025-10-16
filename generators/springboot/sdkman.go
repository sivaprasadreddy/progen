package springboot

type SdkmanConfig struct {
	pg projectGenerator
}

func NewSdkmanConfig(pg projectGenerator) *SdkmanConfig {
	return &SdkmanConfig{pg: pg}
}

func (t SdkmanConfig) generate(pc ProjectConfig) error {
	return t.pg.executeTemplate(pc, "sdkmanrc.tmpl", ".sdkmanrc")
}
