package tpp

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/sets"
)

type Terraforms map[string]*Terraform

func NewTerraformsForDir(dir string) (Terraforms, error) {
	paths, err := getTfFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	tfs := make(Terraforms, len(paths))
	for _, path := range paths {
		tf, err := NewTerraformForFile(path)
		if err != nil {
			return nil, fmt.Errorf("error creating Terraform from file %q; %w", path, err)
		}
		tfs[path] = tf
	}

	return tfs, nil
}

func (t Terraforms) GetModuleSources() ([]string, error) {
	srcs := sets.New[string]()
	for path, tf := range t {
		tfSrc, err := tf.GetModuleSources()
		if err != nil {
			return nil, fmt.Errorf("error getting sources for file %q: %w", path, err)
		}
		srcs.Insert(tfSrc...)
	}
	return srcs.UnsortedList(), nil
}

func (t Terraforms) SetModuleSources(mods map[string]ModuleSource) error {
	for path, tf := range t {
		err := tf.SetModuleSources(mods)
		if err != nil {
			return fmt.Errorf("error setting sources for file %q: %w", path, err)
		}
	}
	return nil
}

func (t Terraforms) WriteTo(dir string) error {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	for path, tf := range t {
		newPath := filepath.Join(dir, filepath.Base(path))
		if _, err := tf.WriteFile(newPath); err != nil {
			return fmt.Errorf("error writing Terraform from source file %q to %q: %w", path, newPath, err)
		}
	}

	return nil
}
