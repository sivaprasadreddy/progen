package springboot

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/sivaprasadreddy/progen/internal/generators/helpers"
)

//go:embed all:templates/*
var tmplsFS embed.FS

const templatesRootDir = "templates"

// ProjectConfig holds all configuration options for generating a Spring Boot project.
type ProjectConfig struct {
	Versions              `json:"-"`
	Frontend              Frontend
	AppName               string
	GroupID               string
	ArtifactID            string
	AppVersion            string
	BasePackage           string
	BuildTool             BuildTool
	PersistenceType       PersistenceType
	DbType                DatabaseType
	DbMigrationTool       DbMigrationTool
	TaskRunner            TaskRunner
	SpringCloudAWSSupport bool
	HTMXSupport           bool
	EmailSupport          bool
	RabbitMQSupport       bool
	RedisCachingSupport   bool
	OpenTelemetrySupport  bool
	AxelixSupport         bool
	K8sSupport            bool
	SpringAIModels        []AIModelProvider
}

// TaskfileSupport reports whether Taskfile (with SDKMAN) is the selected task runner.
func (pc ProjectConfig) TaskfileSupport() bool {
	return pc.TaskRunner == TaskRunnerTaskfile
}

// MiseSupport reports whether mise is the selected task runner.
func (pc ProjectConfig) MiseSupport() bool {
	return pc.TaskRunner == TaskRunnerMise
}

// RestApiEnabled reports whether REST API-specific project features should be generated.
func (pc ProjectConfig) RestApiEnabled() bool {
	return pc.Frontend == FrontendNone || pc.SPASupport()
}

// ThymeleafSupport reports whether Thymeleaf server-rendered views should be generated.
func (pc ProjectConfig) ThymeleafSupport() bool {
	return pc.Frontend == FrontendThymeleaf
}

// JTESupport reports whether JTE server-rendered views should be generated.
func (pc ProjectConfig) JTESupport() bool {
	return pc.Frontend == FrontendJTE
}

// ServerRenderedViews reports whether the app uses a server-side view template engine.
func (pc ProjectConfig) ServerRenderedViews() bool {
	return pc.ThymeleafSupport() || pc.JTESupport()
}

// SpringAISupport reports whether any Spring AI chat model provider is selected.
func (pc ProjectConfig) SpringAISupport() bool {
	return len(pc.SpringAIModels) > 0
}

// OllamaSupport reports whether the Spring AI Ollama chat model is selected.
func (pc ProjectConfig) OllamaSupport() bool {
	return slices.Contains(pc.SpringAIModels, AIModelOllama)
}

// OpenAISupport reports whether the Spring AI OpenAI chat model is selected.
func (pc ProjectConfig) OpenAISupport() bool {
	return slices.Contains(pc.SpringAIModels, AIModelOpenAI)
}

// AnthropicSupport reports whether the Spring AI Anthropic chat model is selected.
func (pc ProjectConfig) AnthropicSupport() bool {
	return slices.Contains(pc.SpringAIModels, AIModelAnthropic)
}

// DefaultAIChatModel returns the provider activated by default via spring.ai.model.chat.
// It is the first selected provider in canonical order (Ollama, OpenAI, Anthropic).
func (pc ProjectConfig) DefaultAIChatModel() AIModelProvider {
	for _, p := range allAIModelProviders {
		if slices.Contains(pc.SpringAIModels, p) {
			return p
		}
	}
	return ""
}

// AngularSupport reports whether an Angular frontend should be generated.
func (pc ProjectConfig) AngularSupport() bool {
	return pc.Frontend == FrontendAngular
}

// ReactSupport reports whether a React frontend should be generated.
func (pc ProjectConfig) ReactSupport() bool {
	return pc.Frontend == FrontendReact
}

// SPASupport reports whether a single page application frontend (Angular or React)
// is built with npm and bundled into the Spring Boot application.
func (pc ProjectConfig) SPASupport() bool {
	return pc.AngularSupport() || pc.ReactSupport()
}

// FrontendDistDir returns the build output directory of the SPA frontend,
// relative to the frontend/ directory.
func (pc ProjectConfig) FrontendDistDir() string {
	if pc.ReactSupport() {
		return "dist"
	}
	return "dist/frontend/browser"
}

var defaultProjectConfig = ProjectConfig{
	Frontend:              FrontendNone,
	AppName:               "myapp",
	GroupID:               "com.mycompany",
	ArtifactID:            "myapp",
	AppVersion:            "0.0.1",
	BasePackage:           "com.mycompany.myapp",
	BuildTool:             Maven,
	PersistenceType:       SpringDataJPA,
	DbType:                PostgreSQL,
	DbMigrationTool:       Flyway,
	TaskRunner:            TaskRunnerMise,
	SpringCloudAWSSupport: false,
	HTMXSupport:           false,
	EmailSupport:          false,
	RabbitMQSupport:       false,
	RedisCachingSupport:   false,
	OpenTelemetrySupport:  false,
	K8sSupport:            false,
	SpringAIModels:        []AIModelProvider{},
}

// DefaultProjectConfig returns the configuration used when no explicit values are provided.
func DefaultProjectConfig() ProjectConfig {
	pc := defaultProjectConfig
	pc.SpringAIModels = slices.Clone(defaultProjectConfig.SpringAIModels)
	return pc
}

// ErrProjectDirExists is returned by Generate when the target project directory
// already exists and GenerateOptions.Overwrite is false.
var ErrProjectDirExists = errors.New("project directory already exists")

// GenerateOptions controls where and how a project is generated.
type GenerateOptions struct {
	// OutputDir is the parent directory in which the project directory (named after
	// ProjectConfig.AppName) is created. Defaults to the current working directory.
	OutputDir string
	// Overwrite deletes an existing project directory before generating.
	// When false, Generate fails with ErrProjectDirExists if the directory exists.
	Overwrite bool
	// SkipCodeFormat skips running the build tool's code formatter after generation.
	SkipCodeFormat bool
	// Output receives the output of external commands run during generation
	// (e.g. the code formatter). Defaults to io.Discard.
	Output io.Writer
}

// GenerateResult describes a successfully generated project.
type GenerateResult struct {
	// ProjectDir is the absolute path of the generated project directory.
	ProjectDir string
	// Warnings describe steps that failed without failing generation, as the
	// code formatter when the required Java version is not installed.
	Warnings []string
}

// Generate validates the given configuration and generates a Spring Boot project
// in <OutputDir>/<AppName>. It never prompts, prints or terminates the process,
// so it can be used by the CLI as well as by other frontends (e.g. a desktop app).
// Cancelling ctx stops generation between steps and kills the code formatter.
// A failing code formatter does not fail generation: the project is kept
// unformatted and the failure is reported in GenerateResult.Warnings.
func Generate(ctx context.Context, pc ProjectConfig, opts GenerateOptions) (GenerateResult, error) {
	if pc.TaskRunner == "" {
		pc.TaskRunner = defaultProjectConfig.TaskRunner
	}
	if err := pc.Validate(); err != nil {
		return GenerateResult{}, err
	}
	pc.Versions = defaultVersions()

	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = "."
	}
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("resolve output directory: %w", err)
	}
	pg := projectGenerator{tmplFS: tmplsFS, outputDir: absOutputDir, output: opts.Output}
	projectDir := pg.projectDir(pc)

	if _, err := os.Stat(projectDir); err == nil && !opts.Overwrite {
		return GenerateResult{}, fmt.Errorf("%w: %s", ErrProjectDirExists, projectDir)
	}
	if err := pg.generate(ctx, pc); err != nil {
		return GenerateResult{}, err
	}
	var warnings []string
	if !opts.SkipCodeFormat {
		if err := formatProjectCode(pg, ctx, pc); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return GenerateResult{}, ctxErr
			}
			warnings = append(warnings, codeFormatWarning(pc, err))
		}
	}
	if err := WriteConfigFile(pc, filepath.Join(projectDir, ProjectConfigFile)); err != nil {
		return GenerateResult{}, err
	}
	return GenerateResult{ProjectDir: projectDir, Warnings: warnings}, nil
}

type configGenerator interface {
	generate(ProjectConfig) error
}

type projectGenerator struct {
	tmplFS embed.FS
	// outputDir is the parent directory of the generated project.
	// When empty, the project is generated relative to the working directory.
	outputDir string
	// output receives the output of external commands. nil discards it.
	output io.Writer
}

// projectDir returns the root directory of the generated project.
func (pg projectGenerator) projectDir(pc ProjectConfig) string {
	return filepath.Join(pg.outputDir, pc.AppName)
}

func (pg projectGenerator) generate(ctx context.Context, pc ProjectConfig) error {
	if err := helpers.RecreateDir(pg.projectDir(pc)); err != nil {
		return err
	}

	generators := []configGenerator{
		NewGitIgnoreConfig(pg),
		NewBuildToolConfig(pg),
		NewAppCommonConfig(pg),
		NewThymeleafConfig(pg),
		NewJteConfig(pg),
		NewDbMigrationsConfig(pg),
		NewSpringModulithConfig(pg),
		NewSecurityConfig(pg),
		NewSpringCloudAwsConfig(pg),
		NewEmailConfig(pg),
		NewRabbitMQConfig(pg),
		NewRedisCachingConfig(pg),
		NewSpringAIConfig(pg),
		NewDockerComposeConfig(pg),
		NewOpenTelemetryConfig(pg),
		NewK8sConfig(pg),
		NewGhActionsConfig(pg),
		NewReadMeConfig(pg),
		NewSdkmanConfig(pg),
		NewTaskfileConfig(pg),
		NewMiseConfig(pg),
		NewRenovateConfig(pg),
		NewAngularConfig(pg),
		NewReactConfig(pg),
		NewAIConfig(pg),
	}

	for _, gen := range generators {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := gen.generate(pc); err != nil {
			return err
		}
	}
	return nil
}

func (pg projectGenerator) copyTemplateDir(pc ProjectConfig, origin, dirName string) error {
	templateDirPath := fmt.Sprintf("%s/%s", templatesRootDir, origin)
	return helpers.CopyDir(pg.tmplFS, templateDirPath, pg.projectDir(pc), dirName)
}

func (pg projectGenerator) copyTemplateFile(pc ProjectConfig, sourceFilePath, targetFilePath string) error {
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, sourceFilePath)
	return helpers.CopyTemplateFile(pg.tmplFS, templateFilePath, pg.projectDir(pc), targetFilePath)
}

// writeTargetFile writes content to targetFilePath, relative to the project directory.
func (pg projectGenerator) writeTargetFile(pc ProjectConfig, content []byte, targetFilePath string) error {
	return helpers.WriteTargetFile(content, filepath.Join(pg.projectDir(pc), targetFilePath))
}

func (pg projectGenerator) executeTemplate(pc ProjectConfig, templatePath, targetFilePath string) error {
	templateFilePath := fmt.Sprintf("%s/%s", templatesRootDir, templatePath)
	tmplFileContent, err := pg.tmplFS.ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(templateFilePath).Parse(string(tmplFileContent))
	if err != nil {
		return err
	}
	f, err := helpers.CreateFile(filepath.Join(pg.projectDir(pc), targetFilePath))
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, pc); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
