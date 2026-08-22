package springboot

import "strings"

type SecurityConfig struct {
	pg projectGenerator
}

func NewSecurityConfig(pg projectGenerator) *SecurityConfig {
	return &SecurityConfig{pg: pg}
}

func (s SecurityConfig) generate(pc ProjectConfig) error {
	if err := s.createSrcMainJava(pc); err != nil {
		return err
	}
	if err := s.createSrcTestJava(pc); err != nil {
		return err
	}
	return nil
}

func (s SecurityConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"domain/Role.java.tmpl":                         "domain/Role.java",
		"domain/User.java.tmpl":                         "domain/User.java",
		"domain/UserRepository.java.tmpl":               "domain/UserRepository.java",
		"domain/UserService.java.tmpl":                  "domain/UserService.java",
		"domain/SecurityUser.java.tmpl":                 "domain/SecurityUser.java",
		"config/SecurityConfig.java.tmpl":               "config/SecurityConfig.java",
		"security/SecurityUserDetailsService.java.tmpl": "security/SecurityUserDetailsService.java",
		"web/UserContextUtils.java.tmpl":                "web/UserContextUtils.java",
		"domain/CreateUserCmd.java.tmpl":                "domain/CreateUserCmd.java",
	}

	if pc.AppType == WebApp {
		templateMap["config/WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["web/WebAppExceptionHandler.java.tmpl"] = "web/GlobalExceptionHandler.java"
		templateMap["web/UserController.java.tmpl"] = "web/UserController.java"
	}

	if pc.AppType == RestApi {
		templateMap["config/JwtWebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["security/AuthToken.java.tmpl"] = "security/AuthToken.java"
		templateMap["security/TokenHelper.java.tmpl"] = "security/TokenHelper.java"
		templateMap["security/TokenAuthenticationFilter.java.tmpl"] = "security/TokenAuthenticationFilter.java"
		templateMap["web/RestApiExceptionHandler.java.tmpl"] = "web/GlobalExceptionHandler.java"
		templateMap["web/LoginRestController.java.tmpl"] = "web/LoginRestController.java"
		templateMap["web/UserRestController.java.tmpl"] = "web/UserRestController.java"
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SecurityConfig) createSrcTestJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{}

	if pc.AppType == WebApp {
		templateMap["web/UserControllerTests.java.tmpl"] = "web/UserControllerTests.java"
	}

	if pc.AppType == RestApi {
		templateMap["web/LoginRestControllerTests.java.tmpl"] = "web/LoginRestControllerTests.java"
		templateMap["web/UserRestControllerTests.java.tmpl"] = "web/UserRestControllerTests.java"
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
