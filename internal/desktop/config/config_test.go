package config

import "testing"

func validConfiguration() GeneratorConfiguration {
	return GeneratorConfiguration{
		Project: ProjectSettings{
			AppName:     "demo",
			AppVersion:  "1.0.0",
			GroupID:     "dev.sivalabs",
			ArtifactID:  "demo",
			PackageName: "dev.sivalabs.demo",
			BuildTool:   BuildToolMaven,
			TaskRunner:  TaskRunnerMise,
		},
		Frontend: FrontendSettings{Type: FrontendNone},
		Database: DatabaseSettings{
			Engine:               DatabasePostgreSQL,
			MigrationTool:        MigrationFlyway,
			PersistenceFramework: PersistenceSpringDataJPA,
		},
	}
}

func TestValidateRejectsInvalidCombinations(t *testing.T) {
	cases := map[string]func(*GeneratorConfiguration){
		"missing app name":    func(c *GeneratorConfiguration) { c.Project.AppName = " " },
		"path in app name":    func(c *GeneratorConfiguration) { c.Project.AppName = "../demo" },
		"missing version":     func(c *GeneratorConfiguration) { c.Project.AppVersion = "" },
		"invalid package":     func(c *GeneratorConfiguration) { c.Project.PackageName = "1dev.sivalabs" },
		"unknown build tool":  func(c *GeneratorConfiguration) { c.Project.BuildTool = "gradle-groovy" },
		"missing task runner": func(c *GeneratorConfiguration) { c.Project.TaskRunner = "" },
		"unknown task runner": func(c *GeneratorConfiguration) { c.Project.TaskRunner = "make" },
		"htmx with react":     func(c *GeneratorConfiguration) { c.Frontend.Type = FrontendReact; c.Frontend.HTMX = true },
		"htmx without views":  func(c *GeneratorConfiguration) { c.Frontend.HTMX = true },
		"unknown ai provider": func(c *GeneratorConfiguration) { c.AIModels = []AIModel{"gemini"} },
		"unknown feature":     func(c *GeneratorConfiguration) { c.Features = []Feature{"htmx"} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			c := validConfiguration()
			mutate(&c)
			if err := Validate(c); err == nil {
				t.Errorf("expected validation error")
			}
		})
	}
}
