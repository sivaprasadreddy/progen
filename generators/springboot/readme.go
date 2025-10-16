package springboot

type ReadMeConfig struct {
	pg projectGenerator
}

func NewReadMeConfig(pg projectGenerator) *ReadMeConfig {
	return &ReadMeConfig{pg: pg}
}

func (r ReadMeConfig) generate(pc ProjectConfig) error {
	return r.pg.executeTemplate(pc, "README.md.tmpl", "README.md")
}
