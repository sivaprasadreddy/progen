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
	if !pc.SecuritySupport && !pc.JwtSecuritySupport {
		return nil
	}
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{}

	if pc.SecuritySupport || pc.JwtSecuritySupport {
		templateMap["domain/Role.java.tmpl"] = "domain/Role.java"
		templateMap["domain/User.java.tmpl"] = "domain/User.java"
		templateMap["domain/UserRepository.java.tmpl"] = "domain/UserRepository.java"
		templateMap["domain/UserService.java.tmpl"] = "domain/UserService.java"
		templateMap["domain/SecurityUser.java.tmpl"] = "domain/SecurityUser.java"
		templateMap["config/SecurityConfig.java.tmpl"] = "config/SecurityConfig.java"
		templateMap["security/SecurityUserDetailsService.java.tmpl"] = "security/SecurityUserDetailsService.java"
		templateMap["web/UserContextUtils.java.tmpl"] = "web/UserContextUtils.java"
		templateMap["domain/CreateUserCmd.java.tmpl"] = "domain/CreateUserCmd.java"
	}

	if pc.SecuritySupport {
		templateMap["config/WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["web/WebAppExceptionHandler.java.tmpl"] = "web/GlobalExceptionHandler.java"
		templateMap["web/UserController.java.tmpl"] = "web/UserController.java"
	}

	if pc.JwtSecuritySupport {
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

	if pc.SecuritySupport {
		templateMap["web/UserControllerTests.java.tmpl"] = "web/UserControllerTests.java"
	}

	if pc.JwtSecuritySupport {
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
