package springboot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateForTaskRunner(t *testing.T, taskRunner TaskRunner) string {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	pc := defaultProjectConfig
	pc.AppName = "taskrunnerapp"
	pc.TaskRunner = taskRunner
	pc.Versions = defaultVersions()
	pg := projectGenerator{tmplFS: tmplsFS}
	for _, gen := range []configGenerator{NewReadMeConfig(pg), NewSdkmanConfig(pg), NewTaskfileConfig(pg), NewMiseConfig(pg), NewGhActionsConfig(pg)} {
		require.NoError(t, gen.generate(pc))
	}
	return filepath.Join(dir, pc.AppName)
}

func TestTaskfileRunnerGeneratesTaskfileAndSdkmanrcOnly(t *testing.T) {
	appDir := generateForTaskRunner(t, TaskRunnerTaskfile)
	require.FileExists(t, filepath.Join(appDir, "Taskfile.yml"))
	require.FileExists(t, filepath.Join(appDir, ".sdkmanrc"))
	require.NoFileExists(t, filepath.Join(appDir, "mise.toml"))

	readme, err := os.ReadFile(filepath.Join(appDir, "README.md"))
	require.NoError(t, err)
	require.Contains(t, string(readme), "## Using Taskfile")
	require.Contains(t, string(readme), "sdk env install")
	require.NotContains(t, string(readme), "mise")

	ci, err := os.ReadFile(filepath.Join(appDir, ".github", "workflows", "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ci), "go-task/setup-task@v2")
	require.Contains(t, string(ci), "java-version-file: .sdkmanrc")
	require.Contains(t, string(ci), "run: task verify")
	require.Contains(t, string(ci), "run: task push_image")
	require.NotContains(t, string(ci), "mise")
}

func TestMiseRunnerGeneratesMiseTomlOnly(t *testing.T) {
	appDir := generateForTaskRunner(t, TaskRunnerMise)
	require.FileExists(t, filepath.Join(appDir, "mise.toml"))
	require.NoFileExists(t, filepath.Join(appDir, "Taskfile.yml"))
	require.NoFileExists(t, filepath.Join(appDir, ".sdkmanrc"))

	readme, err := os.ReadFile(filepath.Join(appDir, "README.md"))
	require.NoError(t, err)
	require.Contains(t, string(readme), "## Using mise")
	require.NotContains(t, string(readme), "Taskfile")
	require.NotContains(t, string(readme), "sdkman")

	ci, err := os.ReadFile(filepath.Join(appDir, ".github", "workflows", "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ci), "jdx/mise-action@v5")
	require.Contains(t, string(ci), "run: mise run verify")
	require.Contains(t, string(ci), "run: mise run push_image")
	require.NotContains(t, string(ci), "setup-java")
	require.NotContains(t, string(ci), "sdkmanrc")
	require.NotContains(t, string(ci), "Taskfile")
}

func TestValidateAndSanitizeDefaultsTaskRunner(t *testing.T) {
	pc := defaultProjectConfig
	pc.TaskRunner = ""
	require.Equal(t, TaskRunnerMise, sanitized(pc).TaskRunner)
	pc.TaskRunner = "bogus"
	require.Equal(t, TaskRunnerMise, sanitized(pc).TaskRunner)
	pc.TaskRunner = TaskRunnerTaskfile
	require.Equal(t, TaskRunnerTaskfile, sanitized(pc).TaskRunner)
}

func sanitized(pc ProjectConfig) ProjectConfig {
	pc, _ = SanitizeConfig(pc)
	return pc
}
