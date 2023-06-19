package parser

import (
	"fmt"
	"regexp"
)

const directivePrefix = "ease"

var (
	reKeyValueExtractor = regexp.MustCompile(`([\w]*?)=([^ ]*)`)
	reDirectiveName     = regexp.MustCompile(fmt.Sprintf(`^%s:(\w+)`, directivePrefix))
)

type Directive struct {
	Name   string            // Name of the directive
	Params map[string]string // Key values of parsed directive params
}

// Try to parse a directive from a sanitized comment (without the //).
func tryParseDirective(comment string) *Directive {
	matches := reDirectiveName.FindStringSubmatch(comment)

	if len(matches) < 2 {
		return nil
	}

	return &Directive{
		Name:   matches[1],
		Params: parseDirectiveParams(comment),
	}
}

// Parse raw directive params into a map of key / values.
func parseDirectiveParams(params string) map[string]string {
	// FIXME: This is sufficient for now, but we might want to use a real parser in the future
	kv := reKeyValueExtractor.FindAllStringSubmatch(params, -1)
	result := make(map[string]string, len(kv))

	// Convert the regex result into a map
	for _, match := range kv {
		if len(match) != 3 {
			continue
		}

		result[match[1]] = match[2]
	}

	return result
}
