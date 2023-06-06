package parser

import (
	"go/ast"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	reAnnotationData            = regexp.MustCompile("( [a-zA-Z-]+=)+")
	annotationParamsReplaceExpr = []byte("\n$1")
)

// Parse comments to see if there is an annotation matching the given one (starting with //).
// If that's the case, parse the annotation and fill the target with the parsed data.
//
// For now, you must provide tag template in the target struct with `ini:"nameofthefield"`.
func ParseAnnotation[T any](annotation string, docs *ast.CommentGroup, target T) (bool, error) {
	for _, doc := range docs.List {
		idx := strings.Index(doc.Text, annotation)

		if idx < 0 {
			continue
		}

		annotationParameters := doc.Text[idx+len(annotation):]
		multilineData := reAnnotationData.ReplaceAll([]byte(annotationParameters), annotationParamsReplaceExpr)
		cfg, err := ini.Load(multilineData)

		if err != nil {
			return true, err
		}

		if err = cfg.MapTo(target); err != nil {
			return true, err
		}

		return true, nil
	}

	return false, nil
}
