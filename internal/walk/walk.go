package walk

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var doNotEditPrefix = []byte("// Code generated by kk; DO NOT EDIT.")
var maxPermissions fs.FileMode = 0750

type pathChecker interface {
	Skip(path string) bool
	Rename(path string) (renamed string)
}

func Walk(content embed.FS, templateDir, targetDir string, data any, pathCkeck pathChecker) error {
	if err := os.MkdirAll(targetDir, maxPermissions); err != nil {
		return fmt.Errorf("MkdirAll %s: %s", targetDir, err)
	}

	if err := fs.WalkDir(content, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if pathCkeck.Skip(path) {
			return nil
		}

		renamedPath := pathCkeck.Rename(path)

		if d.IsDir() {
			if err := os.MkdirAll(replacePrefix(renamedPath, templateDir, targetDir), maxPermissions); err != nil {
				return fmt.Errorf("MkdirAll %s: %s", renamedPath, err)
			}
		} else if strings.HasSuffix(path, ".tmpl") {
			tt := template.Must(template.ParseFS(content, path))
			for _, t := range tt.Templates() {
				filePath := filepath.Clean(strings.TrimSuffix(replacePrefix(renamedPath, templateDir, targetDir), ".tmpl"))
				if _, err := os.Stat(filePath); err == nil {
					if currentFile, err := os.ReadFile(filePath); err != nil {
						return fmt.Errorf("ReadFile %s: %s", filePath, err)
					} else if !bytes.HasPrefix(currentFile, doNotEditPrefix) {
						continue
					}
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("stat: %w", err)
				}

				f, err := os.Create(filePath)
				if err != nil {
					return fmt.Errorf("create %s: %s", filePath, err)
				}
				defer f.Close()
				if err := t.Execute(f, data); err != nil {
					return fmt.Errorf("execute %s: %s", path, err)
				}
			}
		} else {
			filePath := filepath.Clean(replacePrefix(path, templateDir, targetDir))
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("create %s: %s", filePath, err)
			}
			defer f.Close()
		}
		return nil
	}); err != nil {
		return fmt.Errorf("WalkDir %s: %s", templateDir, err)
	}
	return nil
}

func replacePrefix(path, oldPrefix, newPrefix string) string {
	return strings.Replace(path, oldPrefix, newPrefix, 1)
}
