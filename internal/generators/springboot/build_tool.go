package springboot

type BuildToolConfig struct {
	pg projectGenerator
}

func NewBuildToolConfig(pg projectGenerator) *BuildToolConfig {
	return &BuildToolConfig{pg: pg}
}

func (b BuildToolConfig) generate(pc ProjectConfig) error {
	if pc.BuildTool == Maven {
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

func (b BuildToolConfig) createMavenWrapper(pc ProjectConfig) error {
	if err := b.pg.copyTemplateDir(pc, "maven-wrapper", ""); err != nil {
		return err
	}
	return b.pg.executeTemplate(pc, "maven-wrapper.properties.tmpl", ".mvn/wrapper/maven-wrapper.properties")
}

/** Gradle Functions **/

func (b BuildToolConfig) createGradleWrapper(pc ProjectConfig) error {
	if err := b.pg.copyTemplateDir(pc, "gradle-wrapper", ""); err != nil {
		return err
	}
	return b.pg.executeTemplate(pc, "gradle-wrapper.properties.tmpl", "gradle/wrapper/gradle-wrapper.properties")
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
