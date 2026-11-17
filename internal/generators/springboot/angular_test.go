package springboot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFrontendReplacesAppTitleWithArtifactID(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "generated-app")
	pc := ProjectConfig{
		AppName:    projectDir,
		ArtifactID: "orders-service",
	}

	angularConfig := NewAngularConfig(projectGenerator{tmplFS: tmplsFS})
	if err := angularConfig.copyFrontend(pc); err != nil {
		t.Fatalf("copyFrontend() error = %v", err)
	}

	titleFiles := []string{
		filepath.Join(projectDir, "frontend", "src", "index.html"),
		filepath.Join(projectDir, "frontend", "src", "app", "components", "navbar", "navbar.html"),
	}
	for _, filePath := range titleFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", filePath, err)
		}
		if strings.Contains(string(content), angularAppTitle) {
			t.Errorf("generated file %q still contains the default title", filePath)
		}
		if !strings.Contains(string(content), pc.ArtifactID) {
			t.Errorf("generated file %q does not contain ArtifactID %q", filePath, pc.ArtifactID)
		}
	}

	navbarPath := filepath.Join(projectDir, "frontend", "src", "app", "components", "navbar", "navbar.html")
	navbar, err := os.ReadFile(navbarPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", navbarPath, err)
	}
	if !strings.Contains(string(navbar), "{{ authService.user()?.name }}") {
		t.Error("Angular interpolation was unexpectedly processed")
	}
}
