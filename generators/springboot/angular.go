package springboot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const angularAppTitle = "Spring Boot + Angular App"

type AngularConfig struct {
	pg projectGenerator
}

func NewAngularConfig(pg projectGenerator) *AngularConfig {
	return &AngularConfig{pg: pg}
}

func (b AngularConfig) generate(pc ProjectConfig) error {
	if pc.AppType != SpringBootAngularFullStack {
		return nil
	}
	if err := b.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := b.copyFrontend(pc); err != nil {
		return err
	}
	return nil
}

func (b AngularConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"ui/AngularForwardController.java.tmpl": "ui/AngularForwardController.java",
	}

	for tmpl, filePath := range templateMap {
		err := b.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b AngularConfig) copyFrontend(pc ProjectConfig) error {
	if err := b.pg.copyTemplateDir(pc, "angular-frontend", "frontend"); err != nil {
		return err
	}

	titleFiles := []string{
		filepath.Join(pc.AppName, "frontend", "src", "index.html"),
		filepath.Join(pc.AppName, "frontend", "src", "app", "components", "navbar", "navbar.html"),
	}
	for _, filePath := range titleFiles {
		if err := replaceAngularAppTitle(filePath, pc.ArtifactID); err != nil {
			return err
		}
	}
	return nil
}

func replaceAngularAppTitle(filePath, artifactID string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read Angular frontend file %s: %w", filePath, err)
	}
	content = []byte(strings.ReplaceAll(string(content), angularAppTitle, artifactID))
	if err := os.WriteFile(filePath, content, 0755); err != nil {
		return fmt.Errorf("write Angular frontend file %s: %w", filePath, err)
	}
	return nil
}
