package generator

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
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

// archiveFiles reads the archive at path into its entries' contents.
func archiveFiles(t *testing.T, path string) map[string]string {
	t.Helper()
	archive, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer archive.Close()
	entries := map[string]string{}
	for _, f := range archive.File {
		r, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		content, err := io.ReadAll(r)
		_ = r.Close()
		if err != nil {
			t.Fatalf("read %s: %v", f.Name, err)
		}
		entries[f.Name] = string(content)
	}
	return entries
}

// checkSavedAsPreviewed saves prepared and checks the archive holds exactly
// its files, byte for byte.
func checkSavedAsPreviewed(t *testing.T, prepared *Prepared, home string) {
	t.Helper()
	path, err := prepared.Save()
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if want := filepath.Join(home, "Downloads", prepared.Name()+".zip"); path != want {
		t.Errorf("saved at %s, want %s", path, want)
	}
	entries := archiveFiles(t, path)
	if len(entries) != len(prepared.Files()) {
		t.Errorf("archive has %d files, the preview %d", len(entries), len(prepared.Files()))
	}
	for _, f := range prepared.Files() {
		content, ok := entries[prepared.Name()+"/"+f.Path]
		switch {
		case !ok:
			t.Errorf("archive has no %s", f.Path)
		case int64(len(content)) != f.Size:
			t.Errorf("%s: archived %d bytes, previewed %d", f.Path, len(content), f.Size)
		case !f.Binary && content != f.Content:
			t.Errorf("%s: archived contents differ from the preview", f.Path)
		}
	}
}

func TestPreparedProjectIsSavedAsPreviewedAndDiscarded(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	c := validConfiguration()
	c.Project.BuildTool = config.BuildToolGradle

	prepared, err := (&Service{skipCodeFormat: true}).Prepare(c)
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	if prepared.Name() != "demo" {
		t.Errorf("name %q", prepared.Name())
	}
	byPath := map[string]File{}
	for _, f := range prepared.Files() {
		byPath[f.Path] = f
	}
	if f, ok := byPath["build.gradle"]; !ok || f.Binary || !strings.Contains(f.Content, "org.springframework.boot") {
		t.Errorf("build.gradle not previewed as text: %+v", f)
	}
	if f, ok := byPath["gradle/wrapper/gradle-wrapper.jar"]; !ok || !f.Binary || f.Content != "" || f.Size == 0 {
		t.Errorf("wrapper JAR not previewed as binary: %+v", f)
	}
	if _, ok := byPath[".gitignore"]; !ok {
		t.Error("hidden files are not previewed")
	}

	checkSavedAsPreviewed(t, prepared, home)

	directory := prepared.workingDirectory
	prepared.Discard()
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		t.Errorf("prepared directory kept: %v", err)
	}
}

func TestPrepareRejectsAnInvalidConfiguration(t *testing.T) {
	c := validConfiguration()
	c.Project.AppName = "../demo"
	if _, err := (&Service{skipCodeFormat: true}).Prepare(c); err == nil {
		t.Error("prepared an invalid configuration")
	}
}

func TestFormattedPreviewIsWhatIsSaved(t *testing.T) {
	if testing.Short() {
		t.Skip("generates and formats a full project")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	prepared, err := (&Service{}).Prepare(validConfiguration())
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	defer prepared.Discard()
	checkSavedAsPreviewed(t, prepared, home)
}

func TestIsBinary(t *testing.T) {
	for content, want := range map[string]bool{
		"":                      false,
		"plugins {\n}\n":        false,
		"héllo — “quoted”\n":    false,
		"PK\x03\x04\x00\x00":    true,
		"\xff\xfe\x00not utf-8": true,
	} {
		if got := isBinary([]byte(content)); got != want {
			t.Errorf("isBinary(%q) = %v, want %v", content, got, want)
		}
	}
}
