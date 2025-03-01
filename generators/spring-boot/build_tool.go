package spring_boot

type BuildToolConfig struct {
	pg projectGenerator
}

func NewBuildToolConfig(pg projectGenerator) *BuildToolConfig {
	return &BuildToolConfig{pg: pg}
}

func (b BuildToolConfig) generate(pc ProjectConfig) error {
	if pc.BuildTool == BuildToolMaven {
		if err := b.createMavenWrapper(pc); err != nil {
			return err
		}
		if err := b.createMavenBuildFiles(pc); err != nil {
			return err
		}
	} else {
		if err := b.createGradleWrapper(pc); err != nil {
			return err
		}
		if err := b.createGradleBuildFiles(pc); err != nil {
			return err
		}
	}
	return nil
}

func (b BuildToolConfig) createMavenBuildFiles(pc ProjectConfig) error {
	return b.pg.executeTemplate(pc, "pom.xml.tmpl", "pom.xml")
}

func (b BuildToolConfig) createGitIgnore(pc ProjectConfig) error {
	return b.pg.executeTemplate(pc, "gitignore.tmpl", ".gitignore")
}

func (b BuildToolConfig) createMavenWrapper(pc ProjectConfig) error {
	return b.pg.copyTemplateDir(pc, "maven-wrapper", "")
}

/** Gradle Functions **/

func (b BuildToolConfig) createGradleWrapper(pc ProjectConfig) error {
	return b.pg.copyTemplateDir(pc, "gradle-wrapper", "")
}

func (b BuildToolConfig) createGradleBuildFiles(pc ProjectConfig) error {
	templateMap := map[string]string{
		"build.gradle.tmpl":    "build.gradle",
		"settings.gradle.tmpl": "settings.gradle",
	}
	for tmpl, filePath := range templateMap {
		err := b.pg.executeTemplate(pc, tmpl, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
