package springboot

import (
	"fmt"
	"strings"
)

type MiseConfig struct {
	pg projectGenerator
}

func NewMiseConfig(pg projectGenerator) *MiseConfig {
	return &MiseConfig{pg: pg}
}

func (m MiseConfig) generate(pc ProjectConfig) error {
	templateFile := "mise.gradle.toml.tmpl"
	if pc.BuildTool == Maven {
		templateFile = "mise.maven.toml.tmpl"
	}
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, templateFile)
	content, err := m.pg.tmplFS.ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	// mise.toml uses Tera templating ({{env.X}}), so substitute manually instead of using Go templates.
	newContent := strings.ReplaceAll(string(content), "{{ .ArtifactID }}", pc.ArtifactID)

	targetFilePath := fmt.Sprintf("%s/%s", pc.AppName, "mise.toml")
	return m.pg.writeTargetFile([]byte(newContent), targetFilePath)
}
