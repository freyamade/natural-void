package naturalvoid

import (
	"html/template"
	"os"
	"path/filepath"
)

func ParseTemplatesInDir(dir string) (*template.Template, error) {
	// Given a directory, find all the template files inside and pass them to a template.ParseFiles call
	paths := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}
