package springboot

import (
	"fmt"
	"os"
	"strings"
)

// spaFrontend describes a single page application frontend that is built with npm
// and served by Spring Boot as static resources.
type spaFrontend struct {
	// templateDir is the directory under templates/ holding the frontend sources.
	templateDir string
	// defaultTitle is the placeholder app title used in the template sources.
	defaultTitle string
	// titleFiles are paths, relative to the generated frontend directory, containing defaultTitle.
	titleFiles []string
}

const frontendDir = "frontend"

func generateSpaForwardController(pg projectGenerator, pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")
	return pg.executeTemplate(pc,
		srcMainJavaPath+"ui/SpaForwardController.java.tmpl",
		srcMainJavaPath+basePackagePath+"/ui/SpaForwardController.java")
}

// copySpaFrontend copies the frontend sources verbatim (they are not Go templates)
// and replaces the placeholder app title with the project's ArtifactID.
func copySpaFrontend(pg projectGenerator, pc ProjectConfig, spa spaFrontend) error {
	if err := pg.copyTemplateDir(pc, spa.templateDir, frontendDir); err != nil {
		return err
	}
	for _, file := range spa.titleFiles {
		filePath := pc.AppName + "/" + frontendDir + "/" + file
		if err := replaceInFile(filePath, spa.defaultTitle, pc.ArtifactID); err != nil {
			return err
		}
	}
	return nil
}

func replaceInFile(filePath, old, new string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read frontend file %s: %w", filePath, err)
	}
	content = []byte(strings.ReplaceAll(string(content), old, new))
	if err := os.WriteFile(filePath, content, 0755); err != nil {
		return fmt.Errorf("write frontend file %s: %w", filePath, err)
	}
	return nil
}
