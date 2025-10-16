package springboot

type GhActionsConfig struct {
	pg projectGenerator
}

func NewGhActionsConfig(pg projectGenerator) *GhActionsConfig {
	return &GhActionsConfig{pg: pg}
}

func (g GhActionsConfig) generate(pc ProjectConfig) error {
	templateMap := map[string]string{
		"ci/github/workflows/ci.yml.tmpl": ".github/workflows/ci.yml",
	}

	for tmpl, filePath := range templateMap {
		err := g.pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
