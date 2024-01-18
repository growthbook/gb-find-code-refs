package search

import (
	"strings"

	"github.com/growthbook/gb-find-code-refs/aliases"
	"github.com/growthbook/gb-find-code-refs/internal/helpers"
	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/options"
)

type Matcher struct {
	Elements []ElementMatcher
	ctxLines int
}

func NewMultiProjectMatcher(opts options.Options, dir string, flagKeys []string) Matcher {
	elements := make([]ElementMatcher, 0, 1)
	delimiters := strings.Join(GetDelimiters(opts), "")

	projectFlags := flagKeys
	projectAliases := opts.Aliases
	// projectAliases = append(projectAliases, project.Aliases...)
	aliasesByFlagKey, err := aliases.GenerateAliases(projectFlags, projectAliases, dir)
	if err != nil {
		log.Error.Fatalf("failed to generate aliases: %s", err)
	}

	elements = append(elements, NewElementMatcher("default", "", delimiters, projectFlags, aliasesByFlagKey))

	return Matcher{
		ctxLines: opts.ContextLines,
		Elements: elements,
	}
}

func (m Matcher) MatchElement(line, element string) bool {
	for _, em := range m.Elements {
		if e, exists := em.matcherByElement[element]; exists {
			if e.Iter(line).Next() != nil {
				return true
			}
		}
	}

	return false
}

func (m Matcher) GetProjectElementMatcher(projectKey string) *ElementMatcher {
	var elementMatcher ElementMatcher
	for _, element := range m.Elements {
		if element.ProjKey == projectKey {
			elementMatcher = element
			break
		}
	}
	return &elementMatcher
}

func (m Matcher) FindAliases(line, element string) []string {
	matches := make([]string, 0)
	for _, em := range m.Elements {
		matches = append(matches, em.FindAliases(line, element)...)
	}
	return helpers.Dedupe(matches)
}

func (m Matcher) GetElements() (elements [][]string) {
	for _, element := range m.Elements {
		elements = append(elements, element.Elements)
	}
	return elements
}

func buildElementPatterns(flags []string, delimiters string) map[string][]string {
	patternsByFlag := make(map[string][]string, len(flags))
	for _, flag := range flags {
		var patterns []string
		if delimiters != "" {
			patterns = make([]string, 0, len(delimiters)*len(delimiters))
			for _, left := range delimiters {
				for _, right := range delimiters {
					var sb strings.Builder
					sb.Grow(len(flag) + 2)
					sb.WriteRune(left)
					sb.WriteString(flag)
					sb.WriteRune(right)
					patterns = append(patterns, sb.String())
				}
			}
		} else {
			patterns = []string{flag}
		}
		patternsByFlag[flag] = patterns
	}
	return patternsByFlag
}
