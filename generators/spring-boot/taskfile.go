package spring_boot

type TaskfileConfig struct {
	pg projectGenerator
}

func NewTaskfileConfig(pg projectGenerator) *TaskfileConfig {
	return &TaskfileConfig{pg: pg}
}

func (t TaskfileConfig) generate(pc ProjectConfig) error {
	if pc.BuildTool == BuildToolMaven {
		if err := t.createMavenTaskFile(pc); err != nil {
			return err
		}
	} else {
		if err := t.createGradleTaskFile(pc); err != nil {
			return err
		}
	}
	return nil
}

func (t TaskfileConfig) createMavenTaskFile(pc ProjectConfig) error {
	return t.pg.copyTemplateFile(pc, "Taskfile.maven.yml.tmpl", "Taskfile.yml")
}

func (t TaskfileConfig) createGradleTaskFile(pc ProjectConfig) error {
	return t.pg.copyTemplateFile(pc, "Taskfile.gradle.yml.tmpl", "Taskfile.yml")
}
