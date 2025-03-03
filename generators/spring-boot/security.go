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
		templateMap["Role.java.tmpl"] = "domain/model/Role.java"
		templateMap["User.java.tmpl"] = "domain/entities/User.java"
		templateMap["UserRepository.java.tmpl"] = "domain/repositories/UserRepository.java"
		templateMap["UserService.java.tmpl"] = "domain/services/UserService.java"
		templateMap["SecurityUser.java.tmpl"] = "domain/model/SecurityUser.java"
		templateMap["SecurityConfig.java.tmpl"] = "config/SecurityConfig.java"
		templateMap["SecurityUserDetailsService.java.tmpl"] = "security/SecurityUserDetailsService.java"
		templateMap["UserContextUtils.java.tmpl"] = "web/utils/UserContextUtils.java"
		templateMap["CreateUserCmd.java.tmpl"] = "domain/model/CreateUserCmd.java"
	}

	if pc.SecuritySupport {
		templateMap["WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["WebAppExceptionHandler.java.tmpl"] = "web/advices/GlobalExceptionHandler.java"
		templateMap["UserController.java.tmpl"] = "web/controllers/UserController.java"
	}

	if pc.JwtSecuritySupport {
		templateMap["JwtWebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["AuthToken.java.tmpl"] = "security/AuthToken.java"
		templateMap["TokenHelper.java.tmpl"] = "security/TokenHelper.java"
		templateMap["TokenAuthenticationFilter.java.tmpl"] = "security/TokenAuthenticationFilter.java"
		templateMap["RestApiExceptionHandler.java.tmpl"] = "web/advices/GlobalExceptionHandler.java"
		templateMap["LoginRestController.java.tmpl"] = "web/controllers/LoginRestController.java"
		templateMap["UserRestController.java.tmpl"] = "web/controllers/UserRestController.java"
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
		templateMap["UserControllerTests.java.tmpl"] = "web/controllers/UserControllerTests.java"
	}

	if pc.EnabledJwtSecuritySupport() {
		templateMap["LoginRestControllerTests.java.tmpl"] = "web/controllers/LoginRestControllerTests.java"
		templateMap["UserRestControllerTests.java.tmpl"] = "web/controllers/UserRestControllerTests.java"
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
