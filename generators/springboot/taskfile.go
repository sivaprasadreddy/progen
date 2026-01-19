package springboot

type TaskfileConfig struct {
	pg projectGenerator
}

func NewTaskfileConfig(pg projectGenerator) *TaskfileConfig {
	return &TaskfileConfig{pg: pg}
}

func (t TaskfileConfig) generate(pc ProjectConfig) error {
	//TODO; Taskfile templates has {{ .ArtifactID }} to replace
	// But it also has many other placeholders using Go templates syntax that shouldn't be processed.
	// Need to fix it to only replace {{ .ArtifactID }} and ignore other placeholders.
	if pc.BuildTool == Maven {
		return t.pg.copyTemplateFile(pc, "Taskfile.maven.yml.tmpl", "Taskfile.yml")
	}
	return t.pg.copyTemplateFile(pc, "Taskfile.gradle.yml.tmpl", "Taskfile.yml")
}
