package parser

import (
	"go/ast"
	"regexp"
	"strings"
)

const directivePrefix = "//ease:"

var reKeyValueExtractor = regexp.MustCompile(`([\w]*?)=([^ ]*)`)

type Directive struct {
	Name   string            // Name of the directive
	Params map[string]string // Key values of parsed directive params
}

// Parse comments to see if there is a directive matching given ones.
// If that's the case, parse its params.
func ParseDirectives(docs *ast.CommentGroup, names ...string) []*Directive {
	if docs == nil {
		return nil
	}

	directives := make([]*Directive, 0)

	for _, doc := range docs.List {
		for _, name := range names {
			idx := strings.Index(doc.Text, directivePrefix+name)

			if idx < 0 {
				continue
			}

			params := doc.Text[idx+len(directivePrefix+name):]

			directives = append(directives, &Directive{
				Name:   name,
				Params: parseDirectiveParams(params),
			})
		}
	}

	return directives
}

// Parse raw directive params into a map of key / values.
func parseDirectiveParams(params string) map[string]string {
	result := make(map[string]string)

	// FIXME: This is sufficient for now, but we might want to use a real parser in the future
	kv := reKeyValueExtractor.FindAllStringSubmatch(params, -1)

	// Convert the regex result into a map
	for _, match := range kv {
		if len(match) != 3 {
			continue
		}

		result[match[1]] = match[2]
	}

	return result
}
