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
	directives := make([]*Directive, 0)

	for _, doc := range docs.List {
		for _, name := range names {
			idx := strings.Index(doc.Text, directivePrefix+name)

			if idx < 0 {
				continue
			}

			params := doc.Text[idx+len(directivePrefix+name):]

			// Parse params with a regex to extract values from key=value otherkey=other value into a map
			// This is sufficient for now, but we might want to use a real parser in the future
			kv := reKeyValueExtractor.FindAllStringSubmatch(params, -1)

			// Convert the regex result into a map
			paramsMap := make(map[string]string)

			for _, match := range kv {
				if len(match) != 3 {
					continue
				}

				paramsMap[match[1]] = match[2]
			}

			directives = append(directives, &Directive{
				Name:   name,
				Params: paramsMap,
			})
		}
	}

	return directives
}
