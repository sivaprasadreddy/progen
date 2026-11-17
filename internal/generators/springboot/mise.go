package springboot

import (
	"fmt"
	"strconv"
	"strings"
)

type MiseConfig struct {
	pg projectGenerator
}

func NewMiseConfig(pg projectGenerator) *MiseConfig {
	return &MiseConfig{pg: pg}
}

func (m MiseConfig) generate(pc ProjectConfig) error {
	if !pc.MiseSupport() {
		return nil
	}
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
	newContent := strings.NewReplacer(
		"{{ .ArtifactID }}", pc.ArtifactID,
		"{{ .JavaVersion }}", strconv.Itoa(pc.JavaVersion),
		"{{ .MavenVersion }}", pc.MavenVersion,
		"{{ .GradleVersion }}", pc.GradleVersion,
		"{{ .NodeMajorVersion }}", pc.NodeMajorVersion(),
	).Replace(string(content))

	return m.pg.writeTargetFile(pc, []byte(newContent), "mise.toml")
}
