package springboot

type DockerComposeConfig struct {
	pg projectGenerator
}

func NewDockerComposeConfig(pg projectGenerator) *DockerComposeConfig {
	return &DockerComposeConfig{pg: pg}
}

func (d DockerComposeConfig) generate(pc ProjectConfig) error {
	if !pc.DockerComposeSupport {
		return nil
	}
	return d.pg.executeTemplate(pc, "compose.yml.tmpl", "compose.yml")
}
