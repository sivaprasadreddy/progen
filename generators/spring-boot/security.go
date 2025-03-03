package spring_boot

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
	if pc.SecuritySupport == false && pc.JwtSecuritySupport == false {
		return nil
	}
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{}

	if pc.SecuritySupport || pc.JwtSecuritySupport {
		templateMap["domain/model/Role.java.tmpl"] = "domain/model/Role.java"
		templateMap["domain/entities/User.java.tmpl"] = "domain/entities/User.java"
		templateMap["domain/repositories/UserRepository.java.tmpl"] = "domain/repositories/UserRepository.java"
		templateMap["domain/services/UserService.java.tmpl"] = "domain/services/UserService.java"
		templateMap["domain/model/SecurityUser.java.tmpl"] = "domain/model/SecurityUser.java"
		templateMap["security/SecurityConfig.java.tmpl"] = "security/SecurityConfig.java"
		templateMap["security/SecurityUserDetailsService.java.tmpl"] = "security/SecurityUserDetailsService.java"
		templateMap["web/utils/UserContextUtils.java.tmpl"] = "web/utils/UserContextUtils.java"
		templateMap["domain/model/CreateUserCmd.java.tmpl"] = "domain/model/CreateUserCmd.java"
	}

	if pc.SecuritySupport {
		templateMap["config/WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["web/advices/WebAppExceptionHandler.java.tmpl"] = "web/advices/GlobalExceptionHandler.java"
		templateMap["web/controllers/UserController.java.tmpl"] = "web/controllers/UserController.java"
	}

	if pc.JwtSecuritySupport {
		templateMap["config/JwtWebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["security/AuthToken.java.tmpl"] = "security/AuthToken.java"
		templateMap["security/TokenHelper.java.tmpl"] = "security/TokenHelper.java"
		templateMap["security/TokenAuthenticationFilter.java.tmpl"] = "security/TokenAuthenticationFilter.java"
		templateMap["web/advices/RestApiExceptionHandler.java.tmpl"] = "web/advices/GlobalExceptionHandler.java"
		templateMap["web/controllers/LoginRestController.java.tmpl"] = "web/controllers/LoginRestController.java"
		templateMap["web/controllers/UserRestController.java.tmpl"] = "web/controllers/UserRestController.java"
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

	if pc.EnabledSecuritySupport() {
		templateMap["web/controllers/UserControllerTests.java.tmpl"] = "web/controllers/UserControllerTests.java"
	}

	if pc.EnabledJwtSecuritySupport() {
		templateMap["web/controllers/LoginRestControllerTests.java.tmpl"] = "web/controllers/LoginRestControllerTests.java"
		templateMap["web/controllers/UserRestControllerTests.java.tmpl"] = "web/controllers/UserRestControllerTests.java"
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
