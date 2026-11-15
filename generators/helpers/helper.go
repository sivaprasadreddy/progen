package helpers

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const FilePermission = 0755
const DirPermission = 0700

func CopyDir(tmplFS embed.FS, origin, projectName, dirName string) error {
	err := fs.WalkDir(tmplFS, origin, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return copyFile(tmplFS, origin, path, projectName, dirName)
		}
		return nil
	})
	return err
}

func copyFile(tmplFS fs.FS, origin, filePath, projectName, dirName string) error {
	f, err := tmplFS.Open(filePath)
	if err != nil {
		return err
	}
	fileContent, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	targetFile, _ := strings.CutPrefix(filePath, origin+"/")
	targetFilePath := fmt.Sprintf("%s/%s/%s", projectName, dirName, targetFile)

	if err := ensureDir(targetFilePath); err != nil {
		return err
	}

	if err := os.WriteFile(targetFilePath, fileContent, FilePermission); err != nil {
		return err
	}

	return nil
}

func CopyTemplateFile(tmplFS fs.FS, filePath, projectName, destFileName string) error {
	f, err := tmplFS.Open(filePath)
	if err != nil {
		return err
	}
	fileContent, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	targetFilePath := fmt.Sprintf("%s/%s", projectName, destFileName)

	if err := ensureDir(targetFilePath); err != nil {
		return err
	}

	if err := os.WriteFile(targetFilePath, fileContent, FilePermission); err != nil {
		return err
	}

	return nil
}

func WriteTargetFile(content []byte, targetFilePath string) error {
	if err := ensureDir(targetFilePath); err != nil {
		return err
	}
	if err := os.WriteFile(targetFilePath, content, FilePermission); err != nil {
		return err
	}
	return nil
}

func ensureDir(fileName string) error {
	return os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
}

func RecreateDir(dirName string) error {
	if err := os.RemoveAll(path.Join(".", dirName)); err != nil {
		return err
	}
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(filePath string) (*os.File, error) {
	parent := filepath.Dir(filePath)
	_ = os.MkdirAll(parent, DirPermission)
	return os.Create(filePath)
}

func CreateFileWithData(filePath string, data []byte) (*os.File, error) {
	parent := filepath.Dir(filePath)
	_ = os.MkdirAll(parent, DirPermission)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filePath, data, FilePermission)
	if err != nil {
		return nil, err
	}
	return f, err
}

func ValidateApplicationName(val interface{}) error {
	str, ok := val.(string)
	if !ok || strings.Trim(str, " ") == "" {
		return errors.New("application name is required")
	}
	appName := strings.Trim(str, " ")
	if stat, err := os.Stat(appName); err == nil && stat.IsDir() {
		return fmt.Errorf("a directory with name '%s' already exists", appName)
	}
	return nil
}
