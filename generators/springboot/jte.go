package springboot

import "strings"

const srcMainJtePath = "src/main/jte/"

type JteConfig struct {
	pg projectGenerator
}

func NewJteConfig(pg projectGenerator) *JteConfig {
	return &JteConfig{pg: pg}
}

func (j JteConfig) generate(pc ProjectConfig) error {
	if !pc.JTESupport() {
		return nil
	}
	if err := j.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := j.createSrcMainResources(pc); err != nil {
		return err
	}
	return j.createSrcMainJte(pc)
}

func (j JteConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")
	templateMap := map[string]string{
		"shared/web/ViewHelpers.java.tmpl": "shared/web/ViewHelpers.java",
	}
	for tmpl, filePath := range templateMap {
		err := j.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j JteConfig) createSrcMainResources(pc ProjectConfig) error {
	return j.pg.executeTemplate(pc, srcMainResourcesPath+"static/css/styles.css", srcMainResourcesPath+"static/css/styles.css")
}

func (j JteConfig) createSrcMainJte(pc ProjectConfig) error {
	templateMap := map[string]string{
		".jteroot":                ".jteroot",
		"layout.jte.tmpl":         "layout.jte",
		"index.jte.tmpl":          "index.jte",
		"login.jte.tmpl":          "login.jte",
		"registration.jte.tmpl":   "registration.jte",
		"error/403.jte.tmpl":      "error/403.jte",
		"error/404.jte.tmpl":      "error/404.jte",
		"error/500.jte.tmpl":      "error/500.jte",
		"tag/csrf.jte.tmpl":       "tag/csrf.jte",
		"tag/fieldError.jte.tmpl": "tag/fieldError.jte",
	}
	for tmpl, filePath := range templateMap {
		err := j.pg.executeTemplate(pc, srcMainJtePath+tmpl, srcMainJtePath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
