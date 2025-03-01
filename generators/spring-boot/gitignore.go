package spring_boot

type GitIgnoreConfig struct {
	pg projectGenerator
}

func NewGitIgnoreConfig(pg projectGenerator) *GitIgnoreConfig {
	return &GitIgnoreConfig{pg: pg}
}

func (g GitIgnoreConfig) generate(pc ProjectConfig) error {
	return g.createGitIgnore(pc)
}

func (g GitIgnoreConfig) createGitIgnore(pc ProjectConfig) error {
	return g.pg.executeTemplate(pc, "gitignore.tmpl", ".gitignore")
}
