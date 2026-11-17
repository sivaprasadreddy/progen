package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/sivaprasadreddy/progen/internal/desktop/config"
)

const templatesFileName = "templates.json"

// Template is a named, saved selection of the wizard's options.
type Template struct {
	ID            string                        `json:"id"`
	Name          string                        `json:"name"`
	Configuration config.GeneratorConfiguration `json:"configuration"`
	UpdatedAt     time.Time                     `json:"updatedAt"`
}

// TemplateService stores the user's templates in templates.json in a
// directory of its own, the app's user data directory.
type TemplateService struct {
	// dir returns the directory of templates.json.
	dir func() (string, error)
	mu  sync.Mutex
}

// NewTemplateService returns a service that keeps templates.json in the
// directory dir returns.
func NewTemplateService(dir func() (string, error)) *TemplateService {
	return &TemplateService{dir: dir}
}

// List returns the saved templates in the order they were saved: the starter
// templates until the user first saves or deletes one. Updating a template
// keeps its place; new templates go at the end.
func (s *TemplateService) List() ([]Template, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load()
}

// Save stores template and returns it as saved. A template without an ID is
// added as a new template; one with an ID replaces the saved template of
// that ID. Names are unique, ignoring case.
func (s *TemplateService) Save(template Template) (Template, error) {
	template.Name = strings.TrimSpace(template.Name)
	if template.Name == "" {
		return Template{}, errors.New("template name is required")
	}
	if err := config.Validate(template.Configuration); err != nil {
		return Template{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	templates, err := s.load()
	if err != nil {
		return Template{}, err
	}
	for _, existing := range templates {
		if existing.ID != template.ID && strings.EqualFold(existing.Name, template.Name) {
			return Template{}, fmt.Errorf("a template named %q already exists", existing.Name)
		}
	}

	template.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	if template.ID == "" {
		if template.ID, err = newTemplateID(); err != nil {
			return Template{}, err
		}
		templates = append(templates, template)
	} else {
		index := slices.IndexFunc(templates, func(t Template) bool { return t.ID == template.ID })
		if index < 0 {
			return Template{}, fmt.Errorf("template %q no longer exists", template.Name)
		}
		templates[index] = template
	}
	if err := s.store(templates); err != nil {
		return Template{}, err
	}
	return template, nil
}

// Delete removes the template with id. Deleting a template that no longer
// exists is not an error.
func (s *TemplateService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	templates, err := s.load()
	if err != nil {
		return err
	}
	remaining := slices.DeleteFunc(templates, func(t Template) bool { return t.ID == id })
	return s.store(remaining)
}

func (s *TemplateService) path() (string, error) {
	dir, err := s.dir()
	if err != nil {
		return "", fmt.Errorf("find templates directory: %w", err)
	}
	return filepath.Join(dir, templatesFileName), nil
}

func (s *TemplateService) load() ([]Template, error) {
	path, err := s.path()
	if err != nil {
		return nil, err
	}
	var templates []Template
	found, err := readJSONFile(path, &templates)
	if err != nil {
		return nil, fmt.Errorf("read templates: %w", err)
	}
	if !found {
		// Before the user's first change, the starter templates. The first
		// save or delete writes them to templates.json with that change, so
		// from then on they are the user's to update or delete.
		return starterTemplates(), nil
	}
	if templates == nil {
		templates = []Template{}
	}
	return templates, nil
}

func (s *TemplateService) store(templates []Template) error {
	path, err := s.path()
	if err != nil {
		return err
	}
	if err := writeJSONFile(path, templates); err != nil {
		return fmt.Errorf("save templates: %w", err)
	}
	return nil
}

// starterTemplates returns the templates the app comes with, in the order
// they are listed: Maven, Mise, PostgreSQL, Flyway and Spring Data JPA with each kind of
// frontend people start from most.
func starterTemplates() []Template {
	starter := func(id, name string, frontend config.FrontendSettings) Template {
		return Template{
			ID:   id,
			Name: name,
			Configuration: config.GeneratorConfiguration{
				Project: config.ProjectSettings{
					AppName:     "demo",
					AppVersion:  "1.0.0",
					GroupID:     "dev.sivalabs",
					ArtifactID:  "demo",
					PackageName: "dev.sivalabs.demo",
					BuildTool:   config.BuildToolMaven,
					TaskRunner:  config.TaskRunnerMise,
				},
				Frontend: frontend,
				Database: config.DatabaseSettings{
					Engine:               config.DatabasePostgreSQL,
					MigrationTool:        config.MigrationFlyway,
					PersistenceFramework: config.PersistenceSpringDataJPA,
				},
				Features: []config.Feature{},
				AIModels: []config.AIModel{},
			},
		}
	}
	return []Template{
		starter("starter-rest-api", "REST API", config.FrontendSettings{Type: config.FrontendNone}),
		starter("starter-thymeleaf-htmx", "Thymeleaf + HTMX", config.FrontendSettings{Type: config.FrontendThymeleaf, HTMX: true}),
		starter("starter-angular", "Angular", config.FrontendSettings{Type: config.FrontendAngular}),
	}
}

func newTemplateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("create template ID: %w", err)
	}
	return hex.EncodeToString(b), nil
}
