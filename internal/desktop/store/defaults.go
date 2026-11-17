package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/sivaprasadreddy/progen/internal/desktop/config"
)

const defaultsFileName = "defaults.json"

// mavenGroupIDPattern matches what Maven accepts in a group ID.
var mavenGroupIDPattern = regexp.MustCompile(`^[A-Za-z0-9_\-.]+$`)

// ProjectDefaults are the values the Project step starts a new application with.
type ProjectDefaults struct {
	AppVersion string            `json:"appVersion"`
	GroupID    string            `json:"groupId"`
	BuildTool  config.BuildTool  `json:"buildTool"`
	TaskRunner config.TaskRunner `json:"taskRunner"`
}

// BuiltInDefaults are the defaults until the user saves their own.
var BuiltInDefaults = ProjectDefaults{
	AppVersion: "1.0.0",
	GroupID:    "dev.sivalabs",
	BuildTool:  config.BuildToolMaven,
	TaskRunner: config.TaskRunnerMise,
}

// DefaultsService stores the user's project defaults in defaults.json in a
// directory of its own, the app's user data directory.
type DefaultsService struct {
	// dir returns the directory of defaults.json.
	dir func() (string, error)
	mu  sync.Mutex
}

// NewDefaultsService returns a service that keeps defaults.json in the
// directory dir returns.
func NewDefaultsService(dir func() (string, error)) *DefaultsService {
	return &DefaultsService{dir: dir}
}

// Get returns the saved project defaults, or the built-in ones until the user
// saves their own.
func (s *DefaultsService) Get() (ProjectDefaults, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	path, err := s.path()
	if err != nil {
		return ProjectDefaults{}, err
	}
	defaults := BuiltInDefaults
	if _, err := readJSONFile(path, &defaults); err != nil {
		return ProjectDefaults{}, fmt.Errorf("read defaults: %w", err)
	}
	return defaults, nil
}

// Save stores defaults and returns them as saved.
func (s *DefaultsService) Save(defaults ProjectDefaults) (ProjectDefaults, error) {
	defaults.AppVersion = strings.TrimSpace(defaults.AppVersion)
	defaults.GroupID = strings.TrimSpace(defaults.GroupID)
	if defaults.AppVersion == "" {
		return ProjectDefaults{}, errors.New("application version is required")
	}
	if defaults.GroupID == "" {
		return ProjectDefaults{}, errors.New("group ID is required")
	}
	if !mavenGroupIDPattern.MatchString(defaults.GroupID) {
		return ProjectDefaults{}, fmt.Errorf("invalid group ID %q: use letters, digits, '.', '-' and '_'", defaults.GroupID)
	}
	if !slices.Contains([]config.BuildTool{config.BuildToolMaven, config.BuildToolGradle}, defaults.BuildTool) {
		return ProjectDefaults{}, fmt.Errorf("unsupported build tool %q", defaults.BuildTool)
	}
	if !slices.Contains([]config.TaskRunner{config.TaskRunnerMise, config.TaskRunnerTaskfile}, defaults.TaskRunner) {
		return ProjectDefaults{}, fmt.Errorf("unsupported task runner %q", defaults.TaskRunner)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	path, err := s.path()
	if err != nil {
		return ProjectDefaults{}, err
	}
	if err := writeJSONFile(path, defaults); err != nil {
		return ProjectDefaults{}, fmt.Errorf("save defaults: %w", err)
	}
	return defaults, nil
}

func (s *DefaultsService) path() (string, error) {
	dir, err := s.dir()
	if err != nil {
		return "", fmt.Errorf("find defaults directory: %w", err)
	}
	return filepath.Join(dir, defaultsFileName), nil
}
