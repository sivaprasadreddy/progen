package springboot

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateWritesIntoOutputDirWithoutChangingWorkingDir(t *testing.T) {
	cwd := t.TempDir()
	t.Chdir(cwd)
	outputDir := filepath.Join(t.TempDir(), "with space")
	pc := DefaultProjectConfig()
	pc.Frontend = FrontendReact

	result, err := Generate(context.Background(), pc, GenerateOptions{OutputDir: outputDir, SkipCodeFormat: true})
	require.NoError(t, err)
	require.Equal(t, filepath.Join(outputDir, pc.AppName), result.ProjectDir)
	for _, f := range []string{"pom.xml", "mise.toml", ProjectConfigFile, "frontend/package.json",
		"src/main/java/com/mycompany/myapp/ui/SpaForwardController.java"} {
		require.FileExists(t, filepath.Join(result.ProjectDir, f))
	}
	entries, err := os.ReadDir(cwd)
	require.NoError(t, err)
	require.Empty(t, entries, "nothing should be written to the working directory")

	loaded, warnings, err := LoadConfigFile(filepath.Join(result.ProjectDir, ProjectConfigFile))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Equal(t, pc.Frontend, loaded.Frontend)
}

func TestGenerateRefusesExistingProjectDirUnlessOverwrite(t *testing.T) {
	outputDir := t.TempDir()
	pc := DefaultProjectConfig()
	marker := filepath.Join(outputDir, pc.AppName, "keep.txt")
	require.NoError(t, os.MkdirAll(filepath.Dir(marker), 0o755))
	require.NoError(t, os.WriteFile(marker, []byte("x"), 0o644))

	_, err := Generate(context.Background(), pc, GenerateOptions{OutputDir: outputDir, SkipCodeFormat: true})
	require.ErrorIs(t, err, ErrProjectDirExists)
	require.FileExists(t, marker)

	_, err = Generate(context.Background(), pc, GenerateOptions{OutputDir: outputDir, SkipCodeFormat: true, Overwrite: true})
	require.NoError(t, err)
	require.NoFileExists(t, marker)
}

func TestGenerateStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Generate(ctx, DefaultProjectConfig(), GenerateOptions{OutputDir: t.TempDir(), SkipCodeFormat: true})
	require.ErrorIs(t, err, context.Canceled)
}

func TestGenerateRejectsInvalidConfig(t *testing.T) {
	pc := DefaultProjectConfig()
	pc.AppName = "../escape"
	pc.GroupID = ""
	pc.BuildTool = "Ant"
	_, err := Generate(context.Background(), pc, GenerateOptions{OutputDir: t.TempDir(), SkipCodeFormat: true})

	var verrs ValidationErrors
	require.True(t, errors.As(err, &verrs))
	fields := []string{}
	for _, e := range verrs {
		fields = append(fields, e.Field)
	}
	require.Equal(t, []string{"AppName", "GroupID", "BuildTool"}, fields)
}

func TestValidateAppName(t *testing.T) {
	require.NoError(t, ValidateAppName("boot-demo"))
	for _, name := range []string{"", "  ", " app", ".", "..", "a/b", `a\b`} {
		require.Error(t, ValidateAppName(name), name)
	}
}

func TestParseConfigReturnsWarningsForReplacedValues(t *testing.T) {
	data, err := MarshalConfig(DefaultProjectConfig())
	require.NoError(t, err)
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	raw["AppName"] = "demo"
	raw["BuildTool"] = "Ant"
	raw["SpringAIModels"] = []string{"bogus", "openai"}
	data, err = json.Marshal(raw)
	require.NoError(t, err)

	pc, warnings, err := ParseConfig(data)
	require.NoError(t, err)
	require.Equal(t, "demo", pc.AppName)
	require.Equal(t, Maven, pc.BuildTool)
	require.Equal(t, []AIModelProvider{AIModelOpenAI}, pc.SpringAIModels)
	require.Len(t, warnings, 2)
	require.Contains(t, warnings[0], "BuildTool 'Ant' is invalid")
	require.Contains(t, warnings[1], "SpringAIModels entry 'bogus' is invalid")
	require.NoError(t, pc.Validate())
}

func TestFeaturesRoundTrip(t *testing.T) {
	require.NotContains(t, AvailableFeatures(FrontendReact), FeatureHTMXSupport)
	require.Contains(t, AvailableFeatures(FrontendJTE), FeatureHTMXSupport)

	pc := ProjectConfig{Frontend: FrontendThymeleaf}
	pc.SetFeatures([]string{FeatureHTMXSupport, FeatureEmailSupport})
	require.Equal(t, []string{FeatureEmailSupport, FeatureHTMXSupport}, pc.EnabledFeatures())

	pc.Frontend = FrontendNone
	pc.SetFeatures([]string{FeatureHTMXSupport})
	require.False(t, pc.HTMXSupport)
}

func TestChoicesAreValidAndDefaultFirst(t *testing.T) {
	c := Choices()
	def := DefaultProjectConfig()
	require.Equal(t, def.BuildTool.String(), c.BuildTools[0].Value)
	require.Equal(t, def.Frontend.String(), c.Frontends[0].Value)
	require.Equal(t, def.TaskRunner.String(), c.TaskRunners[0].Value)
	for _, ch := range c.Frontends {
		require.True(t, Frontend(ch.Value).IsValid())
	}
	for _, ch := range c.SpringAIModels {
		require.True(t, AIModelProvider(ch.Value).IsValid())
	}
}
