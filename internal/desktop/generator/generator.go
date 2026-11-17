// Package generator generates Spring Boot projects with progen, to preview
// and archive as ZIP files.
package generator

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
)

// Service generates Spring Boot projects.
type Service struct {
	// skipCodeFormat skips the build tool's code formatter, for tests.
	skipCodeFormat bool
}

// GenerateProject generates a Spring Boot project for configuration, saves it
// as a ZIP archive in the user's Downloads directory and returns the
// archive's path.
func (s *Service) GenerateProject(configuration config.GeneratorConfiguration) (string, error) {
	prepared, err := s.Prepare(configuration)
	if err != nil {
		return "", err
	}
	defer prepared.Discard()
	return prepared.Save()
}

// Prepare generates a Spring Boot project for configuration into a
// temporary directory, formatted as an archive would have it, and reads its
// files. Save archives it, and Discard removes it.
func (s *Service) Prepare(configuration config.GeneratorConfiguration) (*Prepared, error) {
	if err := config.Validate(configuration); err != nil {
		return nil, err
	}
	projectConfig := toProjectConfig(configuration)

	workingDirectory, err := os.MkdirTemp("", "progen-desktop-")
	if err != nil {
		return nil, fmt.Errorf("create generation directory: %w", err)
	}
	result, err := springboot.Generate(context.Background(), projectConfig, springboot.GenerateOptions{
		OutputDir:      workingDirectory,
		SkipCodeFormat: s.skipCodeFormat,
	})
	if err != nil {
		_ = os.RemoveAll(workingDirectory)
		return nil, fmt.Errorf("generate project: %w", err)
	}
	files, err := readFiles(result.ProjectDir)
	if err != nil {
		_ = os.RemoveAll(workingDirectory)
		return nil, fmt.Errorf("read generated project: %w", err)
	}
	return &Prepared{workingDirectory: workingDirectory, projectDirectory: result.ProjectDir, name: projectConfig.AppName, files: files}, nil
}

// Prepared is a generated project waiting in a temporary directory to be
// previewed and saved.
type Prepared struct {
	workingDirectory string
	projectDirectory string
	name             string
	files            []File
}

// File is a file of a prepared project.
type File struct {
	// Path is the file's path in the project directory, with slashes.
	Path string
	Size int64
	// Binary is whether the file is not UTF-8 text, as a JAR; Content is
	// empty then.
	Binary  bool
	Content string
}

// Name is the project directory's name, the application's.
func (p *Prepared) Name() string { return filepath.Base(p.projectDirectory) }

// Files are the project's files, in the order a directory walk finds them.
func (p *Prepared) Files() []File { return p.files }

// Save archives the project as a ZIP in the user's Downloads directory,
// named after the application with a numeric suffix if the name is taken,
// and returns the archive's path.
func (p *Prepared) Save() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	downloadsDirectory := filepath.Join(homeDirectory, "Downloads")
	if err := os.MkdirAll(downloadsDirectory, 0o755); err != nil {
		return "", fmt.Errorf("create downloads directory: %w", err)
	}

	zipFile, zipPath, err := createUniqueZipFile(downloadsDirectory, p.name)
	if err != nil {
		return "", fmt.Errorf("create project archive: %w", err)
	}
	if err := zipDirectory(zipFile, p.projectDirectory); err != nil {
		_ = zipFile.Close()
		_ = os.Remove(zipPath)
		return "", fmt.Errorf("archive project: %w", err)
	}
	if err := zipFile.Close(); err != nil {
		_ = os.Remove(zipPath)
		return "", fmt.Errorf("close project archive: %w", err)
	}
	return zipPath, nil
}

// Discard removes the prepared project's directory.
func (p *Prepared) Discard() {
	_ = os.RemoveAll(p.workingDirectory)
}

// readFiles reads every file under directory.
func readFiles(directory string) ([]File, error) {
	files := []File{}
	err := filepath.WalkDir(directory, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		relativePath, err := filepath.Rel(directory, filePath)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		file := File{Path: filepath.ToSlash(relativePath), Size: int64(len(content))}
		if isBinary(content) {
			file.Binary = true
		} else {
			file.Content = string(content)
		}
		files = append(files, file)
		return nil
	})
	return files, err
}

// isBinary reports whether content is not text: not UTF-8, or holding a
// NUL byte.
func isBinary(content []byte) bool {
	return !utf8.Valid(content) || bytes.IndexByte(content, 0) >= 0
}

func createUniqueZipFile(directory, name string) (*os.File, string, error) {
	for suffix := 0; ; suffix++ {
		fileName := name + ".zip"
		if suffix > 0 {
			fileName = fmt.Sprintf("%s-%d.zip", name, suffix)
		}
		filePath := filepath.Join(directory, fileName)
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			return file, filePath, nil
		}
		if !os.IsExist(err) {
			return nil, "", err
		}
	}
}

func toProjectConfig(configuration config.GeneratorConfiguration) springboot.ProjectConfig {
	buildTool := springboot.Maven
	if configuration.Project.BuildTool == config.BuildToolGradle {
		buildTool = springboot.Gradle
	}

	frontends := map[config.FrontendType]springboot.Frontend{
		config.FrontendNone:      springboot.FrontendNone,
		config.FrontendThymeleaf: springboot.FrontendThymeleaf,
		config.FrontendJTE:       springboot.FrontendJTE,
		config.FrontendAngular:   springboot.FrontendAngular,
		config.FrontendReact:     springboot.FrontendReact,
	}
	taskRunners := map[config.TaskRunner]springboot.TaskRunner{
		config.TaskRunnerMise:     springboot.TaskRunnerMise,
		config.TaskRunnerTaskfile: springboot.TaskRunnerTaskfile,
	}
	databaseTypes := map[config.DatabaseEngine]springboot.DatabaseType{
		config.DatabasePostgreSQL: springboot.PostgreSQL,
		config.DatabaseMySQL:      springboot.MySQL,
		config.DatabaseMariaDB:    springboot.MariaDB,
	}
	migrationTools := map[config.MigrationTool]springboot.DbMigrationTool{
		config.MigrationFlyway:    springboot.Flyway,
		config.MigrationLiquibase: springboot.Liquibase,
	}
	persistenceTypes := map[config.PersistenceFramework]springboot.PersistenceType{
		config.PersistenceSpringDataJPA:    springboot.SpringDataJPA,
		config.PersistenceSpringJdbcClient: springboot.SpringJdbcClient,
		config.PersistenceJOOQ:             springboot.SpringJOOQ,
	}

	aiModels := []springboot.AIModelProvider{}
	for _, model := range configuration.AIModels {
		aiModels = append(aiModels, springboot.AIModelProvider(model))
	}

	features := configuration.Features
	return springboot.ProjectConfig{
		Frontend:              frontends[configuration.Frontend.Type],
		AppName:               strings.TrimSpace(configuration.Project.AppName),
		GroupID:               strings.TrimSpace(configuration.Project.GroupID),
		ArtifactID:            strings.TrimSpace(configuration.Project.ArtifactID),
		AppVersion:            strings.TrimSpace(configuration.Project.AppVersion),
		BasePackage:           strings.TrimSpace(configuration.Project.PackageName),
		BuildTool:             buildTool,
		TaskRunner:            taskRunners[configuration.Project.TaskRunner],
		DbType:                databaseTypes[configuration.Database.Engine],
		DbMigrationTool:       migrationTools[configuration.Database.MigrationTool],
		PersistenceType:       persistenceTypes[configuration.Database.PersistenceFramework],
		SpringCloudAWSSupport: slices.Contains(features, config.FeatureSpringCloudAWS),
		EmailSupport:          slices.Contains(features, config.FeatureEmail),
		RabbitMQSupport:       slices.Contains(features, config.FeatureRabbitMQ),
		RedisCachingSupport:   slices.Contains(features, config.FeatureRedisCaching),
		OpenTelemetrySupport:  slices.Contains(features, config.FeatureOpenTelemetry),
		K8sSupport:            slices.Contains(features, config.FeatureKubernetes),
		AxelixSupport:         slices.Contains(features, config.FeatureAxelix),
		HTMXSupport:           configuration.Frontend.HTMX && config.IsServerRendered(configuration.Frontend.Type),
		SpringAIModels:        aiModels,
	}
}

func zipDirectory(destination io.Writer, sourceDirectory string) error {
	archive := zip.NewWriter(destination)
	err := filepath.Walk(sourceDirectory, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(filepath.Dir(sourceDirectory), filePath)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relativePath)
		header.Method = zip.Deflate
		entry, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(entry, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
	if err != nil {
		_ = archive.Close()
		return err
	}
	return archive.Close()
}
