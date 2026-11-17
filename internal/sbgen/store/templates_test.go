package store

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sivaprasadreddy/progen/internal/sbgen/config"
)

func templateNames(templates []Template) []string {
	names := make([]string, len(templates))
	for i, template := range templates {
		names[i] = template.Name
	}
	return names
}

// newTestTemplateService returns a service whose templates.json is in a
// temporary directory. It starts with no templates unless withStarters.
func newTestTemplateService(t *testing.T, withStarters bool) (*TemplateService, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "SBGen")
	if !withStarters {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, templatesFileName), []byte("[]"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return NewTemplateService(func() (string, error) { return dir, nil }), dir
}

func TestTemplateServiceStartsWithStarterTemplates(t *testing.T) {
	s, dir := newTestTemplateService(t, true)
	templates, err := s.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	want := map[string]config.FrontendSettings{
		"Angular":          {Type: config.FrontendAngular},
		"REST API":         {Type: config.FrontendNone},
		"Thymeleaf + HTMX": {Type: config.FrontendThymeleaf, HTMX: true},
	}
	if got := templateNames(templates); !slices.Equal(got, []string{"REST API", "Thymeleaf + HTMX", "Angular"}) {
		t.Fatalf("unexpected starter templates or order: %v", got)
	}
	for _, template := range templates {
		frontend, ok := want[template.Name]
		if !ok {
			t.Errorf("unexpected starter template %q", template.Name)
			continue
		}
		c := template.Configuration
		if c.Frontend != frontend {
			t.Errorf("%s: frontend %+v, want %+v", template.Name, c.Frontend, frontend)
		}
		if c.Project.BuildTool != config.BuildToolMaven || c.Project.TaskRunner != config.TaskRunnerMise ||
			c.Database != (config.DatabaseSettings{Engine: config.DatabasePostgreSQL, MigrationTool: config.MigrationFlyway, PersistenceFramework: config.PersistenceSpringDataJPA}) ||
			len(c.Features) != 0 || len(c.AIModels) != 0 {
			t.Errorf("%s: unexpected configuration %+v", template.Name, c)
		}
		if err := config.Validate(c); err != nil {
			t.Errorf("%s is not a valid configuration: %v", template.Name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, templatesFileName)); !os.IsNotExist(err) {
		t.Errorf("listing wrote templates.json: %v", err)
	}

	// The first save keeps the starters, in their order, and adds the new template at the end.
	if _, err := s.Save(Template{Name: "Mine", Configuration: validConfiguration()}); err != nil {
		t.Fatalf("save: %v", err)
	}
	templates, _ = s.List()
	if got := templateNames(templates); !slices.Equal(got, []string{"REST API", "Thymeleaf + HTMX", "Angular", "Mine"}) {
		t.Errorf("expected the starters, then the new template, got %v", got)
	}

	// Deleted starters stay deleted.
	for _, template := range templates {
		if err := s.Delete(template.ID); err != nil {
			t.Fatalf("delete %s: %v", template.Name, err)
		}
	}
	if templates, _ = s.List(); len(templates) != 0 {
		t.Errorf("deleted templates came back: %+v", templates)
	}
}

func TestTemplateServiceSavesUpdatesAndDeletes(t *testing.T) {
	s, dir := newTestTemplateService(t, false)

	jte := validConfiguration()
	jte.Frontend = config.FrontendSettings{Type: config.FrontendJTE, HTMX: true}
	saved, err := s.Save(Template{Name: "  JTE + HTMX  ", Configuration: jte})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID == "" || saved.Name != "JTE + HTMX" || saved.UpdatedAt.IsZero() {
		t.Errorf("unexpected saved template: %+v", saved)
	}
	if _, err := os.Stat(filepath.Join(dir, templatesFileName)); err != nil {
		t.Errorf("templates file not written: %v", err)
	}
	if _, err := s.Save(Template{Name: "api", Configuration: validConfiguration()}); err != nil {
		t.Fatalf("save second: %v", err)
	}

	// A fresh service reads what the first one saved, in the order it was saved.
	templates, err := NewTemplateService(s.dir).List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(templates) != 2 || templates[0].Name != "JTE + HTMX" || templates[1].Name != "api" {
		t.Fatalf("unexpected templates: %+v", templates)
	}
	if !templates[0].Configuration.Frontend.HTMX {
		t.Errorf("configuration not kept: %+v", templates[0].Configuration)
	}

	saved.Name = "JTE"
	saved.Configuration.Database.Engine = config.DatabaseMySQL
	updated, err := s.Save(saved)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.ID != saved.ID {
		t.Errorf("update changed the ID: %s -> %s", saved.ID, updated.ID)
	}
	templates, _ = s.List()
	// The updated template keeps its place in the list.
	if len(templates) != 2 || templates[0].Name != "JTE" || templates[0].Configuration.Database.Engine != config.DatabaseMySQL {
		t.Fatalf("update not saved in place: %+v", templates)
	}

	if err := s.Delete(saved.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := s.Delete(saved.ID); err != nil {
		t.Errorf("deleting again: %v", err)
	}
	templates, _ = s.List()
	if len(templates) != 1 || templates[0].Name != "api" {
		t.Errorf("unexpected templates after delete: %+v", templates)
	}
}

func TestTemplateServiceRejectsInvalidTemplates(t *testing.T) {
	s, _ := newTestTemplateService(t, false)
	existing, err := s.Save(Template{Name: "Default", Configuration: validConfiguration()})
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	invalid := validConfiguration()
	invalid.Frontend.HTMX = true
	cases := map[string]Template{
		"missing name":        {Name: "  ", Configuration: validConfiguration()},
		"duplicate name":      {Name: "default", Configuration: validConfiguration()},
		"invalid selection":   {Name: "Broken", Configuration: invalid},
		"unknown template ID": {ID: "missing", Name: "Gone", Configuration: validConfiguration()},
	}
	for name, template := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := s.Save(template); err == nil {
				t.Errorf("expected an error")
			}
		})
	}

	// Saving a template under its own name is not a duplicate.
	if _, err := s.Save(existing); err != nil {
		t.Errorf("re-saving a template: %v", err)
	}
}

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
