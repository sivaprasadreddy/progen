package springboot

type AIConfig struct {
	pg projectGenerator
}

func NewAIConfig(pg projectGenerator) *AIConfig {
	return &AIConfig{pg: pg}
}

func (a AIConfig) generate(pc ProjectConfig) error {
	err := a.pg.copyTemplateFile(pc, "CLAUDE.md.tmpl", "CLAUDE.md")
	if err != nil {
		return err
	}
	return a.pg.executeTemplate(pc, "AGENTS.md.tmpl", "AGENTS.md")
}
