package helpers

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

	ensureDir(targetFilePath)

	if err := os.WriteFile(targetFilePath, fileContent, 0666); err != nil {
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

	ensureDir(targetFilePath)

	if err := os.WriteFile(targetFilePath, fileContent, 0666); err != nil {
		return err
	}

	return nil
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
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

func CreateFile(filePath string) *os.File {
	parent := filepath.Dir(filePath)
	_ = os.MkdirAll(parent, 0700)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return f
}

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FatalIfErrOrMsg(err error, msg string) {
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(msg)
	}
}
