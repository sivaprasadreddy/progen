package spring_boot

type ThymeleafConfig struct {
	pg projectGenerator
}

func NewThymeleafConfig(pg projectGenerator) *ThymeleafConfig {
	return &ThymeleafConfig{pg: pg}
}

func (t ThymeleafConfig) generate(pc ProjectConfig) error {
	return t.createThymeleafTemplateFiles(pc)
}

func (t ThymeleafConfig) createThymeleafTemplateFiles(pc ProjectConfig) error {
	if !pc.ThymeleafSupport {
		return nil
	}
	templateMap := map[string]string{}

	if pc.ThymeleafSupport {
		templateMap["static/css/styles.css"] = "static/css/styles.css"
		templateMap["templates/index.html.tmpl"] = "templates/index.html"
		templateMap["templates/layout.html.tmpl"] = "templates/layout.html"
	}

	if pc.ThymeleafSupport && pc.SecuritySupport {
		templateMap["templates/login.html.tmpl"] = "templates/login.html"
		templateMap["templates/registration.html.tmpl"] = "templates/registration.html"
		templateMap["templates/registration-success.html.tmpl"] = "templates/registration-success.html"
	}

	for tmpl, filePath := range templateMap {
		err := t.pg.executeTemplate(pc, srcMainResourcesPath+tmpl, srcMainResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
