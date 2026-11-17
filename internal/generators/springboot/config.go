package springboot

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/sivaprasadreddy/progen/internal/generators/helpers"
)

// FieldError describes an invalid ProjectConfig field.
type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidationErrors is returned by ProjectConfig.Validate when one or more fields are invalid.
type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	msgs := make([]string, len(v))
	for i, e := range v {
		msgs[i] = e.Error()
	}
	return "invalid project config: " + strings.Join(msgs, "; ")
}

// Validate reports every invalid field of the configuration.
// It returns nil or a ValidationErrors value.
func (pc ProjectConfig) Validate() error {
	var errs ValidationErrors
	if err := ValidateAppName(pc.AppName); err != nil {
		errs = append(errs, FieldError{"AppName", err.Error()})
	}
	required := []struct{ field, value string }{
		{"GroupID", pc.GroupID},
		{"ArtifactID", pc.ArtifactID},
		{"AppVersion", pc.AppVersion},
		{"BasePackage", pc.BasePackage},
	}
	for _, r := range required {
		if strings.TrimSpace(r.value) == "" {
			errs = append(errs, FieldError{r.field, "is required"})
		}
	}
	enums := []struct {
		field string
		value fmt.Stringer
		valid bool
	}{
		{"Frontend", pc.Frontend, pc.Frontend.IsValid()},
		{"BuildTool", pc.BuildTool, pc.BuildTool.IsValid()},
		{"PersistenceType", pc.PersistenceType, pc.PersistenceType.IsValid()},
		{"DbType", pc.DbType, pc.DbType.IsValid()},
		{"DbMigrationTool", pc.DbMigrationTool, pc.DbMigrationTool.IsValid()},
		{"TaskRunner", pc.TaskRunner, pc.TaskRunner.IsValid()},
	}
	for _, e := range enums {
		if !e.valid {
			errs = append(errs, FieldError{e.field, fmt.Sprintf("invalid value '%s'", e.value)})
		}
	}
	for _, m := range pc.SpringAIModels {
		if !m.IsValid() {
			errs = append(errs, FieldError{"SpringAIModels", fmt.Sprintf("invalid value '%s'", m)})
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateAppName checks that name can be used as the generated project's directory name.
func ValidateAppName(name string) error {
	trimmed := strings.TrimSpace(name)
	switch {
	case trimmed == "":
		return fmt.Errorf("application name is required")
	case trimmed != name:
		return fmt.Errorf("application name must not have leading or trailing spaces")
	case name == "." || name == ".." || strings.ContainsAny(name, `/\`):
		return fmt.Errorf("application name '%s' must be a plain directory name", name)
	}
	return nil
}

// DeriveBasePackage suggests a Java base package from the group and artifact IDs,
// e.g. ("com.mycompany", "boot-demo") -> "com.mycompany.bootdemo".
func DeriveBasePackage(groupID, artifactID string) string {
	groupID = strings.TrimSpace(groupID)
	segment := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(artifactID), "-", ""))
	if groupID == "" {
		return segment
	}
	if segment == "" {
		return groupID
	}
	return groupID + "." + segment
}

// ParseConfig decodes a JSON project config (the .progen.json format) and sanitizes it.
// Invalid values are replaced with defaults; a warning is returned for each replacement.
func ParseConfig(data []byte) (ProjectConfig, []string, error) {
	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return config, nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	config, warnings := SanitizeConfig(config)
	return config, warnings, nil
}

// LoadConfigFile reads and sanitizes a JSON project config file. See ParseConfig.
func LoadConfigFile(configFile string) (ProjectConfig, []string, error) {
	if strings.TrimSpace(configFile) == "" {
		return ProjectConfig{}, nil, fmt.Errorf("configFile name is empty")
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return ProjectConfig{}, nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return ParseConfig(data)
}

// SanitizeConfig replaces invalid values with defaults and returns a warning for each replacement.
func SanitizeConfig(pc ProjectConfig) (ProjectConfig, []string) {
	var warnings []string
	warn := func(format string, args ...any) {
		warnings = append(warnings, fmt.Sprintf(format, args...))
	}

	// Validate string fields
	if strings.TrimSpace(pc.AppName) == "" {
		warn("AppName '%s' is invalid. Using default: %s", pc.AppName, defaultProjectConfig.AppName)
		pc.AppName = defaultProjectConfig.AppName
	}
	if strings.TrimSpace(pc.GroupID) == "" {
		warn("GroupID '%s' is invalid. Using default: %s", pc.GroupID, defaultProjectConfig.GroupID)
		pc.GroupID = defaultProjectConfig.GroupID
	}
	if strings.TrimSpace(pc.ArtifactID) == "" {
		warn("ArtifactID '%s' is invalid. Using default: %s", pc.ArtifactID, defaultProjectConfig.ArtifactID)
		pc.ArtifactID = defaultProjectConfig.ArtifactID
	}
	if strings.TrimSpace(pc.AppVersion) == "" {
		warn("AppVersion '%s' is invalid. Using default: %s", pc.AppVersion, defaultProjectConfig.AppVersion)
		pc.AppVersion = defaultProjectConfig.AppVersion
	}
	if strings.TrimSpace(pc.BasePackage) == "" {
		warn("BasePackage '%s' is invalid. Using default: %s", pc.BasePackage, defaultProjectConfig.BasePackage)
		pc.BasePackage = defaultProjectConfig.BasePackage
	}
	if !pc.Frontend.IsValid() {
		warn("Frontend '%s' is invalid. Using default: %s", pc.Frontend, defaultProjectConfig.Frontend)
		pc.Frontend = defaultProjectConfig.Frontend
	}
	if !pc.BuildTool.IsValid() {
		warn("BuildTool '%s' is invalid. Using default: %s", pc.BuildTool, defaultProjectConfig.BuildTool)
		pc.BuildTool = defaultProjectConfig.BuildTool
	}
	if !pc.PersistenceType.IsValid() {
		warn("PersistenceType '%s' is invalid. Using default: %s", pc.PersistenceType, defaultProjectConfig.PersistenceType)
		pc.PersistenceType = defaultProjectConfig.PersistenceType
	}
	if !pc.DbType.IsValid() {
		warn("DbType '%s' is invalid. Using default: %s", pc.DbType, defaultProjectConfig.DbType)
		pc.DbType = defaultProjectConfig.DbType
	}
	if !pc.DbMigrationTool.IsValid() {
		warn("DbMigrationTool '%s' is invalid. Using default: %s", pc.DbMigrationTool, defaultProjectConfig.DbMigrationTool)
		pc.DbMigrationTool = defaultProjectConfig.DbMigrationTool
	}
	// Configs created before the TaskRunner option existed have no value.
	if !pc.TaskRunner.IsValid() {
		if pc.TaskRunner != "" {
			warn("TaskRunner '%s' is invalid. Using default: %s", pc.TaskRunner, defaultProjectConfig.TaskRunner)
		}
		pc.TaskRunner = defaultProjectConfig.TaskRunner
	}
	for _, m := range pc.SpringAIModels {
		if !m.IsValid() {
			warn("SpringAIModels entry '%s' is invalid. Valid values: %v. Ignoring it.", m, allAIModelProviders)
		}
	}
	pc.SpringAIModels = NormalizeAIModels(pc.SpringAIModels)

	return pc, warnings
}

// NormalizeAIModels drops invalid or duplicate providers and returns them in canonical order.
func NormalizeAIModels(models []AIModelProvider) []AIModelProvider {
	result := []AIModelProvider{}
	for _, p := range allAIModelProviders {
		if slices.Contains(models, p) {
			result = append(result, p)
		}
	}
	return result
}

// MarshalConfig encodes the configuration in the .progen.json format.
func MarshalConfig(pc ProjectConfig) ([]byte, error) {
	data, err := json.MarshalIndent(pc, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal project config: %w", err)
	}
	return data, nil
}

// WriteConfigFile writes the configuration to filePath in the .progen.json format.
func WriteConfigFile(pc ProjectConfig, filePath string) error {
	data, err := MarshalConfig(pc)
	if err != nil {
		return err
	}
	if err := helpers.WriteTargetFile(data, filePath); err != nil {
		return fmt.Errorf("failed to write %s file: %w", ProjectConfigFile, err)
	}
	return nil
}
