package spring_boot

type ReadMeConfig struct {
	pg projectGenerator
}

func NewReadMeConfig(pg projectGenerator) *ReadMeConfig {
	return &ReadMeConfig{pg: pg}
}

func (r ReadMeConfig) generate(pc ProjectConfig) error {
	return r.createReadMeFile(pc)
}

func (r ReadMeConfig) createReadMeFile(pc ProjectConfig) error {
	templateMap := map[string]string{
		"README.md.tmpl": "README.md",
	}
	for tmpl, filePath := range templateMap {
		err := r.pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
