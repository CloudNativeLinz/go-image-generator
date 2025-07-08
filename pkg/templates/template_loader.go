package templates

import (
	"io/ioutil"
	"path/filepath"
)

// LoadTemplates loads image templates from the specified directory.
func LoadTemplates(templateDir string) ([]string, error) {
	var templates []string

	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			templates = append(templates, filepath.Join(templateDir, file.Name()))
		}
	}

	return templates, nil
}
