package app

import (
	"path"
	"strings"

	"github.com/egoist/mygo/ui/editor"
)

// The languages the preview highlights beyond MyGo's own: Gradle's Groovy
// DSL, Java's .properties files and JTE templates.

// gradleLanguage highlights build.gradle and settings.gradle with Kotlin's
// tokenizer, whose strings, interpolations, comments and calls are
// Groovy's too.
var gradleLanguage = func() *editor.Language {
	l := *editor.LanguageNamed("Kotlin")
	l.Name = "Gradle"
	l.Extensions = []string{".gradle"}
	l.Filenames = nil
	l.Interpreters = nil
	return &l
}()

var propertiesLanguage = &editor.Language{
	Name:        "Properties",
	Extensions:  []string{".properties"},
	LineComment: "#",
	Tokenizer:   editor.TokenizerFunc(tokenizeProperties),
}

// jteLanguage is HTML with JTE's directives (@if, @param, @template.layout),
// its expressions (${…}, $unsafe{…}, !{…}) and its comments (<%-- … --%>).
var jteLanguage = func() *editor.Language {
	html := editor.LanguageNamed("HTML")
	l := *html
	l.Name = "JTE"
	l.Extensions = []string{".jte"}
	l.Filenames = nil
	l.Tokenizer = editor.TokenizerFunc(func(line string, state uint64, tokens []editor.Token) ([]editor.Token, uint64) {
		start := len(tokens)
		tokens, next := html.Tokenizer.Tokenize(line, state, tokens)
		kinds := tokenKinds(tokens[start:], len(line))
		markJTE(line, kinds)
		return appendTokens(tokens[:start], kinds), next
	})
	return &l
}()

// languageFor is the language a preview highlights a file in: by its name,
// else by its #! line, as mvnw's; nil for plain text.
func languageFor(name, text string) *editor.Language {
	switch path.Ext(name) {
	case ".gradle":
		return gradleLanguage
	case ".properties":
		return propertiesLanguage
	case ".jte":
		return jteLanguage
	}
	return editor.LanguageForText(name, text)
}

// tokenizeProperties splits a line of a .properties file into a key, its
// separator and its value, with ${…} placeholders in the value; lines
// starting with # or ! are comments.
func tokenizeProperties(line string, _ uint64, tokens []editor.Token) ([]editor.Token, uint64) {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return tokens, 0
	}
	if trimmed[0] == '#' || trimmed[0] == '!' {
		return append(tokens, editor.Token{End: len(line), Kind: editor.TokenComment}), 0
	}
	kinds := make([]editor.TokenKind, len(line))
	separator := strings.IndexAny(line, "=:")
	keyEnd := len(line)
	if separator >= 0 {
		keyEnd = separator
	}
	for i := len(line) - len(trimmed); i < keyEnd; i++ {
		kinds[i] = editor.TokenProperty
	}
	if separator >= 0 {
		kinds[separator] = editor.TokenOperator
		for i := separator + 1; i < len(line); i++ {
			kinds[i] = editor.TokenString
		}
		value := line[separator+1:]
		for offset := 0; ; {
			open := strings.Index(value[offset:], "${")
			if open < 0 {
				break
			}
			open += offset
			end := strings.IndexByte(value[open:], '}')
			if end < 0 {
				end = len(value)
			} else {
				end += open + 1
			}
			for i := open; i < end; i++ {
				kinds[separator+1+i] = editor.TokenEscape
			}
			offset = end
		}
	}
	return appendTokens(tokens, kinds), 0
}

// markJTE overwrites the kinds of line's bytes that are JTE's: comments,
// directives and expressions.
func markJTE(line string, kinds []editor.TokenKind) {
	for i := 0; i < len(line); {
		switch {
		case strings.HasPrefix(line[i:], "<%--"):
			end := strings.Index(line[i:], "--%>")
			if end < 0 {
				end = len(line)
			} else {
				end += i + len("--%>")
			}
			fill(kinds, i, end, editor.TokenComment)
			i = end
		case strings.HasPrefix(line[i:], "${"), strings.HasPrefix(line[i:], "!{"), strings.HasPrefix(line[i:], "$unsafe{"):
			open := i + strings.IndexByte(line[i:], '{')
			end := matchingBrace(line, open)
			fill(kinds, i, open+1, editor.TokenEscape)
			fill(kinds, open+1, end, editor.TokenVariable)
			if end < len(line) {
				kinds[end] = editor.TokenEscape
				end++
			}
			i = end
		case line[i] == '@' && (i == 0 || !isWordByte(line[i-1])) && i+1 < len(line) && isLetter(line[i+1]):
			end := i + 1
			for end < len(line) && (isWordByte(line[end]) || line[end] == '.') {
				end++
			}
			fill(kinds, i, end, editor.TokenKeyword)
			i = end
		default:
			i++
		}
	}
}

// matchingBrace returns the offset of the } closing the { at open, or the
// end of the line.
func matchingBrace(line string, open int) int {
	depth := 0
	for i := open; i < len(line); i++ {
		switch line[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(line)
}

func isLetter(b byte) bool   { return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' }
func isWordByte(b byte) bool { return isLetter(b) || b >= '0' && b <= '9' || b == '_' }

func fill(kinds []editor.TokenKind, start, end int, kind editor.TokenKind) {
	for i := start; i < end && i < len(kinds); i++ {
		kinds[i] = kind
	}
}

// tokenKinds is the kind of each of a line's bytes.
func tokenKinds(tokens []editor.Token, length int) []editor.TokenKind {
	kinds := make([]editor.TokenKind, length)
	start := 0
	for _, t := range tokens {
		fill(kinds, start, t.End, t.Kind)
		start = t.End
	}
	return kinds
}

// appendTokens appends the runs of kinds to tokens.
func appendTokens(tokens []editor.Token, kinds []editor.TokenKind) []editor.Token {
	for i := 0; i < len(kinds); {
		end := i + 1
		for end < len(kinds) && kinds[end] == kinds[i] {
			end++
		}
		tokens = append(tokens, editor.Token{End: end, Kind: kinds[i]})
		i = end
	}
	return tokens
}
