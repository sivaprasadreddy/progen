package springboot

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (pg projectGenerator) formatCode(pc ProjectConfig) error {
	executable, formatCmd := pg.getCodeFormatCommand(pc.BuildTool)
	appFormatCmd := pg.buildCommandString(pc.AppName, executable, formatCmd)
	cmd := pg.createOSCommand(appFormatCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed to format project code\n%w", output)
	}
	return err
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

func (pg projectGenerator) buildCommandString(dirName, executable, formatCmd string) string {
	return fmt.Sprintf("cd %s && %s %s", dirName, executable, formatCmd)
}

func (pg projectGenerator) createOSCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("/bin/sh", "-c", command)
}
