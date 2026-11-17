package springboot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReactCopyFrontendReplacesAppTitleWithArtifactID(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "generated-app")
	pc := ProjectConfig{
		AppName:    projectDir,
		ArtifactID: "orders-service",
		Frontend:   FrontendReact,
	}

	reactConfig := NewReactConfig(projectGenerator{tmplFS: tmplsFS})
	if err := reactConfig.copyFrontend(pc); err != nil {
		t.Fatalf("copyFrontend() error = %v", err)
	}

	for _, file := range reactFrontend.titleFiles {
		filePath := filepath.Join(projectDir, "frontend", file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", filePath, err)
		}
		if strings.Contains(string(content), reactAppTitle) {
			t.Errorf("generated file %q still contains the default title", filePath)
		}
		if !strings.Contains(string(content), pc.ArtifactID) {
			t.Errorf("generated file %q does not contain ArtifactID %q", filePath, pc.ArtifactID)
		}
	}

	navbarPath := filepath.Join(projectDir, "frontend", "src", "components", "Navbar.tsx")
	navbar, err := os.ReadFile(navbarPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", navbarPath, err)
	}
	if !strings.Contains(string(navbar), "{user?.name}") {
		t.Error("JSX expressions were unexpectedly processed")
	}

	for _, file := range []string{"package.json", "vite.config.ts", ".gitignore", "public/favicon.ico", "src/main.tsx"} {
		if _, err := os.Stat(filepath.Join(projectDir, "frontend", file)); err != nil {
			t.Errorf("expected generated file %q: %v", file, err)
		}
	}
}

func TestReactGenerateCreatesSpaForwardController(t *testing.T) {
	// executeTemplate resolves the project directory relative to the working directory.
	t.Chdir(t.TempDir())
	projectDir := "generated-app"
	pc := ProjectConfig{
		AppName:     projectDir,
		ArtifactID:  "orders-service",
		BasePackage: "com.example.orders",
		Frontend:    FrontendReact,
	}

	if err := NewReactConfig(projectGenerator{tmplFS: tmplsFS}).generate(pc); err != nil {
		t.Fatalf("generate() error = %v", err)
	}

	controllerPath := filepath.Join(projectDir, "src", "main", "java", "com", "example", "orders", "ui", "SpaForwardController.java")
	content, err := os.ReadFile(controllerPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", controllerPath, err)
	}
	if !strings.Contains(string(content), "package com.example.orders.ui;") {
		t.Error("SpaForwardController has unexpected package declaration")
	}
}

func TestReactGenerateSkipsWhenNotSelected(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "generated-app")
	pc := ProjectConfig{AppName: projectDir, Frontend: FrontendAngular}

	if err := NewReactConfig(projectGenerator{tmplFS: tmplsFS}).generate(pc); err != nil {
		t.Fatalf("generate() error = %v", err)
	}
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		t.Errorf("expected no output for non-React frontend, stat err = %v", err)
	}
}

func TestFrontendFlags(t *testing.T) {
	tests := []struct {
		frontend Frontend
		spa      bool
		restAPI  bool
		distDir  string
	}{
		{FrontendReact, true, true, "dist"},
		{FrontendAngular, true, true, "dist/frontend/browser"},
		{FrontendNone, false, true, "dist/frontend/browser"},
		{FrontendThymeleaf, false, false, "dist/frontend/browser"},
	}
	for _, tt := range tests {
		pc := ProjectConfig{Frontend: tt.frontend}
		if got := pc.SPASupport(); got != tt.spa {
			t.Errorf("%s: SPASupport() = %v, want %v", tt.frontend, got, tt.spa)
		}
		if got := pc.RestApiEnabled(); got != tt.restAPI {
			t.Errorf("%s: RestApiEnabled() = %v, want %v", tt.frontend, got, tt.restAPI)
		}
		if tt.spa && pc.FrontendDistDir() != tt.distDir {
			t.Errorf("%s: FrontendDistDir() = %q, want %q", tt.frontend, pc.FrontendDistDir(), tt.distDir)
		}
	}
	if !FrontendReact.IsValid() {
		t.Error("FrontendReact should be valid")
	}
}
