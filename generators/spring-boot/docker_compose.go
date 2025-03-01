package spring_boot

type DockerComposeConfig struct {
	pg projectGenerator
}

func NewDockerComposeConfig(pg projectGenerator) *DockerComposeConfig {
	return &DockerComposeConfig{pg: pg}
}

func (d DockerComposeConfig) generate(pc ProjectConfig) error {
	return d.createComposeConfigFiles(pc)
}

func (d DockerComposeConfig) createComposeConfigFiles(pc ProjectConfig) error {
	if !pc.DockerComposeSupport {
		return nil
	}
	templateMap := map[string]string{
		"compose.yml.tmpl": "compose.yml",
	}
	for tmpl, filePath := range templateMap {
		err := d.pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
