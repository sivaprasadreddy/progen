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
		"users/Role.java.tmpl":                       "users/Role.java",
		"users/UserDto.java.tmpl":                    "users/UserDto.java",
		"users/UserService.java.tmpl":                "users/UserService.java",
		"users/SecurityUser.java.tmpl":               "users/SecurityUser.java",
		"config/SecurityConfig.java.tmpl":            "config/SecurityConfig.java",
		"users/SecurityUserDetailsService.java.tmpl": "users/SecurityUserDetailsService.java",
		"users/CreateUserCmd.java.tmpl":              "users/CreateUserCmd.java",
		"users/AuthUtils.java.tmpl":                  "users/AuthUtils.java",
		"users/UserNotFoundException.java.tmpl":      "users/UserNotFoundException.java",
		"users/LoginRequest.java.tmpl":               "users/LoginRequest.java",
		"users/LoginResponse.java.tmpl":              "users/LoginResponse.java",
		"users/RefreshTokenResponse.java.tmpl":       "users/RefreshTokenResponse.java",
	}

	if pc.PersistenceType == SpringDataJPA {
		templateMap["users/JpaUserEntity.java.tmpl"] = "users/UserEntity.java"
		templateMap["users/JpaUserRepository.java.tmpl"] = "users/UserRepository.java"
	}
	if pc.PersistenceType == SpringJdbcClient {
		templateMap["users/JdbcClientUserEntity.java.tmpl"] = "users/UserEntity.java"
		templateMap["users/JdbcClientUserRepository.java.tmpl"] = "users/UserRepository.java"
	}

	if pc.AppType == WebApp {
		templateMap["config/WebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["config/WebAppExceptionHandler.java.tmpl"] = "config/GlobalExceptionHandler.java"
		templateMap["users/UserController.java.tmpl"] = "users/UserController.java"
	}

	if pc.RestApiEnabled() {
		templateMap["config/JwtWebSecurityConfig.java.tmpl"] = "config/WebSecurityConfig.java"
		templateMap["users/AuthToken.java.tmpl"] = "users/AuthToken.java"
		templateMap["users/TokenHelper.java.tmpl"] = "users/TokenHelper.java"
		templateMap["users/JwtFilter.java.tmpl"] = "users/JwtFilter.java"
		templateMap["config/RestApiExceptionHandler.java.tmpl"] = "config/GlobalExceptionHandler.java"
		templateMap["users/AuthController.java.tmpl"] = "users/AuthController.java"
		templateMap["users/UserRestController.java.tmpl"] = "users/UserRestController.java"
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
		templateMap["users/UserControllerTests.java.tmpl"] = "users/UserControllerTests.java"
	}

	if pc.RestApiEnabled() {
		templateMap["users/AuthControllerTests.java.tmpl"] = "users/AuthControllerTests.java"
		templateMap["users/UserRestControllerTests.java.tmpl"] = "users/UserRestControllerTests.java"
	}

	for tmpl, filePath := range templateMap {
		err := s.pg.executeTemplate(pc, srcTestJavaPath+tmpl, srcTestJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
