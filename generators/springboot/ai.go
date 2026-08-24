package springboot

type AIConfig struct {
	pg projectGenerator
}

func NewAIConfig(pg projectGenerator) *AIConfig {
	return &AIConfig{pg: pg}
}

func (a AIConfig) generate(pc ProjectConfig) error {
	return a.pg.executeTemplate(pc, "AGENTS.md.tmpl", "AGENTS.md")
}
