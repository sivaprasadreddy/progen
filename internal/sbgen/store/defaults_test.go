package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sivaprasadreddy/progen/internal/sbgen/config"
)

func newTestDefaultsService(t *testing.T) (*DefaultsService, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "SBGen")
	return NewDefaultsService(func() (string, error) { return dir, nil }), dir
}

func TestDefaultsServiceReturnsBuiltInDefaultsUntilSaved(t *testing.T) {
	s, dir := newTestDefaultsService(t)
	defaults, err := s.Get()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if defaults != BuiltInDefaults {
		t.Errorf("expected the built-in defaults, got %+v", defaults)
	}
	if _, err := os.Stat(filepath.Join(dir, defaultsFileName)); !os.IsNotExist(err) {
		t.Errorf("getting wrote defaults.json: %v", err)
	}
}

func TestDefaultsServiceSavesDefaults(t *testing.T) {
	s, _ := newTestDefaultsService(t)
	saved, err := s.Save(ProjectDefaults{AppVersion: " 0.1.0-SNAPSHOT ", GroupID: " com.example ", BuildTool: config.BuildToolGradle, TaskRunner: config.TaskRunnerTaskfile})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := ProjectDefaults{AppVersion: "0.1.0-SNAPSHOT", GroupID: "com.example", BuildTool: config.BuildToolGradle, TaskRunner: config.TaskRunnerTaskfile}
	if saved != want {
		t.Errorf("saved %+v, want %+v", saved, want)
	}
	// A fresh service reads what the first one saved.
	got, err := NewDefaultsService(s.dir).Get()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestDefaultsServiceRejectsInvalidDefaults(t *testing.T) {
	cases := map[string]func(*ProjectDefaults){
		"missing version":     func(d *ProjectDefaults) { d.AppVersion = " " },
		"missing group ID":    func(d *ProjectDefaults) { d.GroupID = "" },
		"invalid group ID":    func(d *ProjectDefaults) { d.GroupID = "com example" },
		"unknown build tool":  func(d *ProjectDefaults) { d.BuildTool = "ant" },
		"unknown task runner": func(d *ProjectDefaults) { d.TaskRunner = "make" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			s, dir := newTestDefaultsService(t)
			defaults := BuiltInDefaults
			mutate(&defaults)
			if _, err := s.Save(defaults); err == nil {
				t.Errorf("expected an error")
			}
			if _, err := os.Stat(filepath.Join(dir, defaultsFileName)); !os.IsNotExist(err) {
				t.Errorf("invalid defaults were written: %v", err)
			}
		})
	}
}
