package springboot

const reactAppTitle = "Spring Boot + React App"

var reactFrontend = spaFrontend{
	templateDir:  "react-frontend",
	defaultTitle: reactAppTitle,
	titleFiles: []string{
		"index.html",
		"src/config/app.ts",
	},
}

// ReactConfig generates a React + Vite + Tailwind CSS single page application
// in the frontend/ directory, served by Spring Boot.
type ReactConfig struct {
	pg projectGenerator
}

func NewReactConfig(pg projectGenerator) *ReactConfig {
	return &ReactConfig{pg: pg}
}

func (r ReactConfig) generate(pc ProjectConfig) error {
	if !pc.ReactSupport() {
		return nil
	}
	if err := generateSpaForwardController(r.pg, pc); err != nil {
		return err
	}
	return r.copyFrontend(pc)
}

func (r ReactConfig) copyFrontend(pc ProjectConfig) error {
	return copySpaFrontend(r.pg, pc, reactFrontend)
}
