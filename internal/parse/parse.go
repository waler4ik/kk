package parse

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/waler4ik/kk/internal/tmpl"
)

func RouterType(path string) (string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
	if err != nil {
		return "", fmt.Errorf("parse file: %w", err)
	}
	for _, i := range f.Imports {
		if i.Path.Value == "\"github.com/go-chi/chi/v5\"" {
			return tmpl.Chi, nil
		} else if i.Path.Value == "\"github.com/gin-gonic/gin\"" {
			return tmpl.Gin, nil
		}
	}
	return "", fmt.Errorf("router type not found")
}
