package spring_boot

type TaskfileConfig struct {
	pg projectGenerator
}

func NewTaskfileConfig(pg projectGenerator) *TaskfileConfig {
	return &TaskfileConfig{pg: pg}
}

func (t TaskfileConfig) generate(pc ProjectConfig) error {
	if pc.BuildTool == Maven {
		return t.pg.copyTemplateFile(pc, "Taskfile.maven.yml.tmpl", "Taskfile.yml")
	}
	return t.pg.copyTemplateFile(pc, "Taskfile.gradle.yml.tmpl", "Taskfile.yml")
}
