package springboot

const angularAppTitle = "Spring Boot + Angular App"

var angularFrontend = spaFrontend{
	templateDir:  "angular-frontend",
	defaultTitle: angularAppTitle,
	titleFiles: []string{
		"src/index.html",
		"src/app/components/navbar/navbar.html",
		"src/app/layout/dashboard-layout/dashboard-layout.html",
	},
}

type AngularConfig struct {
	pg projectGenerator
}

func NewAngularConfig(pg projectGenerator) *AngularConfig {
	return &AngularConfig{pg: pg}
}

func (b AngularConfig) generate(pc ProjectConfig) error {
	if !pc.AngularSupport() {
		return nil
	}
	if err := generateSpaForwardController(b.pg, pc); err != nil {
		return err
	}
	return b.copyFrontend(pc)
}

func (b AngularConfig) copyFrontend(pc ProjectConfig) error {
	return copySpaFrontend(b.pg, pc, angularFrontend)
}
