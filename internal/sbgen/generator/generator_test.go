package generator

import (
	"archive/zip"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
	"github.com/sivaprasadreddy/progen/internal/sbgen/config"
)

func validConfiguration() config.GeneratorConfiguration {
	return config.GeneratorConfiguration{
		Project: config.ProjectSettings{
			AppName:     "demo",
			AppVersion:  "1.0.0",
			GroupID:     "dev.sivalabs",
			ArtifactID:  "demo",
			PackageName: "dev.sivalabs.demo",
			BuildTool:   config.BuildToolMaven,
			TaskRunner:  config.TaskRunnerMise,
		},
		Frontend: config.FrontendSettings{Type: config.FrontendNone},
		Database: config.DatabaseSettings{
			Engine:               config.DatabasePostgreSQL,
			MigrationTool:        config.MigrationFlyway,
			PersistenceFramework: config.PersistenceSpringDataJPA,
		},
	}
}

func TestToProjectConfigMapsJTEWithHTMXAndAIModels(t *testing.T) {
	c := validConfiguration()
	c.Project.BuildTool = config.BuildToolGradle
	c.Frontend.Type = config.FrontendJTE
	c.Frontend.HTMX = true
	c.Features = []config.Feature{config.FeatureEmail, config.FeatureKubernetes, config.FeatureAxelix}
	c.AIModels = []config.AIModel{config.AIModelOpenAI, config.AIModelOllama}

	if err := config.Validate(c); err != nil {
		t.Fatalf("expected valid configuration, got %v", err)
	}
	pc := toProjectConfig(c)

	if pc.Frontend != springboot.FrontendJTE || !pc.JTESupport() || !pc.HTMXSupport {
		t.Errorf("unexpected frontend mapping: %+v", pc)
	}
	if pc.BuildTool != springboot.Gradle {
		t.Errorf("expected Gradle, got %s", pc.BuildTool)
	}
	if !pc.EmailSupport || !pc.K8sSupport || !pc.AxelixSupport || pc.RabbitMQSupport {
		t.Errorf("unexpected feature flags: %+v", pc)
	}
	if !slices.Contains(pc.SpringAIModels, springboot.AIModelOpenAI) || !slices.Contains(pc.SpringAIModels, springboot.AIModelOllama) {
		t.Errorf("unexpected AI models: %v", pc.SpringAIModels)
	}
	if pc.AppName != "demo" || pc.AppVersion != "1.0.0" {
		t.Errorf("unexpected app identity: %s %s", pc.AppName, pc.AppVersion)
	}
}

func TestToProjectConfigMapsReactWithTaskfile(t *testing.T) {
	c := validConfiguration()
	c.Project.TaskRunner = config.TaskRunnerTaskfile
	c.Frontend.Type = config.FrontendReact

	if err := config.Validate(c); err != nil {
		t.Fatalf("expected valid configuration, got %v", err)
	}
	pc := toProjectConfig(c)

	if pc.Frontend != springboot.FrontendReact || !pc.ReactSupport() || !pc.SPASupport() || pc.HTMXSupport {
		t.Errorf("unexpected frontend mapping: %+v", pc)
	}
	if pc.TaskRunner != springboot.TaskRunnerTaskfile || !pc.TaskfileSupport() || pc.MiseSupport() {
		t.Errorf("expected Taskfile task runner, got %q", pc.TaskRunner)
	}
}

func TestToProjectConfigMapsMiseTaskRunner(t *testing.T) {
	pc := toProjectConfig(validConfiguration())
	if pc.TaskRunner != springboot.TaskRunnerMise || !pc.MiseSupport() {
		t.Errorf("expected Mise task runner, got %q", pc.TaskRunner)
	}
}

func TestGenerateProjectWritesZipToDownloads(t *testing.T) {
	if testing.Short() {
		t.Skip("generates a full project")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	path, err := (&Service{}).GenerateProject(validConfiguration())
	if err != nil {
		t.Fatalf("generate project: %v", err)
	}
	if want := filepath.Join(home, "Downloads", "demo.zip"); path != want {
		t.Errorf("expected archive at %s, got %s", want, path)
	}
	archive, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer archive.Close()
	if !slices.ContainsFunc(archive.File, func(f *zip.File) bool { return f.Name == "demo/pom.xml" }) {
		t.Errorf("archive has no demo/pom.xml")
	}

	// A second generation must not overwrite the first archive.
	second, err := (&Service{}).GenerateProject(validConfiguration())
	if err != nil {
		t.Fatalf("generate second project: %v", err)
	}
	if filepath.Base(second) != "demo-1.zip" {
		t.Errorf("expected demo-1.zip, got %s", second)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("first archive is gone: %v", err)
	}
}
