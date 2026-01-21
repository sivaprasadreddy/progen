package springboot

import (
	"fmt"
	"strings"
)

type TaskfileConfig struct {
	pg projectGenerator
}

func NewTaskfileConfig(pg projectGenerator) *TaskfileConfig {
	return &TaskfileConfig{pg: pg}
}

func (t TaskfileConfig) generate(pc ProjectConfig) error {
	templateFile := "Taskfile.gradle.yml.tmpl"
	if pc.BuildTool == Maven {
		templateFile = "Taskfile.maven.yml.tmpl"
	}
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, templateFile)
	content, err := t.pg.tmplFS.ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	newContent := strings.ReplaceAll(string(content), "{{ .ArtifactID }}", pc.ArtifactID)

	targetFilePath := fmt.Sprintf("%s/%s", pc.AppName, "Taskfile.yml")
	return t.pg.writeTargetFile([]byte(newContent), targetFilePath)
}
