package tpp

import (
	"fmt"
	"os"
	"path/filepath"
)

var validTerraformExtensions = []string{
	".tf",
	".tf.json",
}

func getTfFiles(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory entries: %w", err)
	}

	files := make([]string, 0, len(ents))
	for _, ent := range ents {
		if hasTerraformExtension(ent.Name()) && !ent.IsDir() {
			files = append(files, filepath.Join(dir, ent.Name()))
		}
	}

	return files, nil
}

func hasTerraformExtension(path string) bool {
	for _, ext := range validTerraformExtensions {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}
