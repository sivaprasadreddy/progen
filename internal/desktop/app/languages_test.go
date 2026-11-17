package app

import (
	"testing"

	"github.com/egoist/mygo/ui/editor"
)

func TestPreviewHighlightsTheGeneratedFileTypes(t *testing.T) {
	for name, want := range map[string]string{
		"Application.java":       "Java",
		"pom.xml":                "XML",
		"build.gradle":           "Gradle",
		"settings.gradle":        "Gradle",
		"application.properties": "Properties",
		"application.yml":        "YAML",
		"compose.yaml":           "YAML",
		"V1__create_tables.sql":  "SQL",
		"README.md":              "Markdown",
		"renovate.json":          "JSON",
		"mise.toml":              "TOML",
		"run.sh":                 "Shell",
		"index.html":             "HTML",
		"home.jte":               "JTE",
		"main.ts":                "TypeScript",
		"App.tsx":                "TypeScript",
		"vite.config.js":         "JavaScript",
		"styles.css":             "CSS",
	} {
		l := languageFor(name, "")
		if l == nil || l.Name != want {
			t.Errorf("%s highlighted as %v, want %s", name, l, want)
		}
	}
	if l := languageFor("mvnw", "#!/bin/sh\necho\n"); l == nil || l.Name != "Shell" {
		t.Errorf("mvnw highlighted as %v, want Shell", l)
	}
	if l := languageFor(".gitignore", "target/\n"); l != nil {
		t.Errorf(".gitignore highlighted as %s, want plain text", l.Name)
	}
}

// runs is line split into its tokens' text and kinds.
func runs(l *editor.Language, line string) map[string]editor.TokenKind {
	tokens, _ := l.Tokenizer.Tokenize(line, 0, nil)
	got := map[string]editor.TokenKind{}
	start := 0
	for _, tk := range tokens {
		got[line[start:tk.End]] = tk.Kind
		start = tk.End
	}
	return got
}

func TestPropertiesTokens(t *testing.T) {
	got := runs(propertiesLanguage, "spring.datasource.url=jdbc:postgresql://${DB_HOST:localhost}/app")
	want := map[string]editor.TokenKind{
		"spring.datasource.url": editor.TokenProperty,
		"=":                     editor.TokenOperator,
		"jdbc:postgresql://":    editor.TokenString,
		"${DB_HOST:localhost}":  editor.TokenEscape,
		"/app":                  editor.TokenString,
	}
	for text, kind := range want {
		if got[text] != kind {
			t.Errorf("%q is %v, want %v (tokens %v)", text, got[text], kind, got)
		}
	}
	if got := runs(propertiesLanguage, "  # a comment"); got["  # a comment"] != editor.TokenComment {
		t.Errorf("comment: %v", got)
	}
}

func TestJTETokens(t *testing.T) {
	got := runs(jteLanguage, `@if(user != null)<p>${user.name()} <%-- note --%></p>@endif`)
	want := map[string]editor.TokenKind{
		"@if":            editor.TokenKeyword,
		"${":             editor.TokenEscape,
		"user.name()":    editor.TokenVariable,
		"<%-- note --%>": editor.TokenComment,
		"@endif":         editor.TokenKeyword,
	}
	for text, kind := range want {
		if got[text] != kind {
			t.Errorf("%q is %v, want %v (tokens %v)", text, got[text], kind, got)
		}
	}
	if got := runs(jteLanguage, `<a href="mailto:me@example.com">`); got["@example"] == editor.TokenKeyword {
		t.Error("an e-mail address highlighted as a directive")
	}
}
