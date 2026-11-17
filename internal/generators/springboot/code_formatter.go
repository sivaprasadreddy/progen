package springboot

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

func (pg projectGenerator) formatCode(ctx context.Context, pc ProjectConfig) error {
	executable, formatCmd := pg.getCodeFormatCommand(pc.BuildTool)
	cmd := pg.createOSCommand(ctx, executable, formatCmd)
	cmd.Dir = pg.projectDir(pc)
	output := pg.output
	if output == nil {
		output = io.Discard
	}
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format project code: %w", err)
	}
	return nil
}

func (pg projectGenerator) getCodeFormatCommand(buildTool BuildTool) (executable, formatCmd string) {
	isWindows := runtime.GOOS == "windows"

	if buildTool == Gradle {
		if isWindows {
			return "gradlew.bat", "spotlessApply"
		}
		return "./gradlew", "spotlessApply"
	}

	if isWindows {
		return "mvnw.cmd", "spotless:apply"
	}
	return "./mvnw", "spotless:apply"
}

// createOSCommand builds the command to run executable in the project directory.
// The caller sets cmd.Dir; a relative executable path is resolved against it.
func (pg projectGenerator) createOSCommand(ctx context.Context, executable, arg string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", executable, arg)
	}
	return exec.CommandContext(ctx, executable, arg)
}
