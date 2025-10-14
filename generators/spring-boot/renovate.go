package spring_boot

type RenovateConfig struct {
	pg projectGenerator
}

func NewRenovateConfig(pg projectGenerator) *RenovateConfig {
	return &RenovateConfig{pg: pg}
}

func (t RenovateConfig) generate(pc ProjectConfig) error {
	return t.pg.copyTemplateFile(pc, "renovate.json.tmpl", "renovate.json")
}
