// Package generator generates Spring Boot projects with progen and archives
// them as ZIP files.
package generator

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sivaprasadreddy/progen/internal/generators/springboot"
	"github.com/sivaprasadreddy/progen/internal/sbgen/config"
)

// Service generates Spring Boot projects.
type Service struct{}

// GenerateProject generates a Spring Boot project for configuration, saves it
// as a ZIP archive in the user's Downloads directory and returns the
// archive's path.
func (s *Service) GenerateProject(configuration config.GeneratorConfiguration) (string, error) {
	if err := config.Validate(configuration); err != nil {
		return "", err
	}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	downloadsDirectory := filepath.Join(homeDirectory, "Downloads")
	if err := os.MkdirAll(downloadsDirectory, 0o755); err != nil {
		return "", fmt.Errorf("create downloads directory: %w", err)
	}

	projectConfig := toProjectConfig(configuration)

	workingDirectory, err := os.MkdirTemp("", "sbgen-")
	if err != nil {
		return "", fmt.Errorf("create generation directory: %w", err)
	}
	defer os.RemoveAll(workingDirectory)

	result, err := springboot.Generate(context.Background(), projectConfig, springboot.GenerateOptions{
		OutputDir: workingDirectory,
	})
	if err != nil {
		return "", fmt.Errorf("generate project: %w", err)
	}

	zipFile, zipPath, err := createUniqueZipFile(downloadsDirectory, projectConfig.AppName)
	if err != nil {
		return "", fmt.Errorf("create project archive: %w", err)
	}
	if err := zipDirectory(zipFile, result.ProjectDir); err != nil {
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
