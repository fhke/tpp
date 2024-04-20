package tpp

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	blockTypeModule = "module"
	attrNameSource  = "source"
	attrNameVersion = "version"
)

type (
	Terraform struct {
		fi        *hclwrite.File
		namespace string
	}
	ModuleSource struct {
		Source  string
		Version *string
	}
)

func NewTerraformForFile(path, namespace string) (*Terraform, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	fi, hclErr := hclwrite.ParseConfig(data, path, hcl.InitialPos)
	if hclErr != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &Terraform{fi: fi, namespace: namespace}, nil
}

func (t *Terraform) GetModuleSources() ([]string, error) {
	srcSet := sets.New[string]()
	if err := t.iterModules(func(source string, _ *hclwrite.Block) error {
		srcSet.Insert(source)
		return nil
	}); err != nil {
		return nil, err
	}
	return srcSet.UnsortedList(), nil
}

func (t *Terraform) SetModuleSources(mods map[string]ModuleSource) error {
	return t.iterModules(func(source string, bl *hclwrite.Block) error {
		modSrc, ok := mods[source]
		if !ok {
			return fmt.Errorf("no replacement found for %q", source)
		}

		setStringAttr(bl, attrNameSource, modSrc.Source)
		if modSrc.Version != nil {
			setStringAttr(bl, attrNameVersion, *modSrc.Version)
		} else {
			bl.Body().RemoveAttribute(attrNameVersion)
		}

		return nil
	})
}

func (t *Terraform) Bytes() []byte {
	return t.fi.Bytes()
}

func (t *Terraform) Write(wr io.Writer) (int64, error) {
	return t.fi.WriteTo(wr)
}

func (t *Terraform) WriteFile(path string) (int64, error) {
	fi, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("error creating file: %w", err)
	}
	defer fi.Close()
	return t.Write(fi)
}

func (t *Terraform) isManagedModule(name string) bool {
	if t.namespace == "" {
		return true
	}
	return strings.HasPrefix(name, t.namespace+"::")
}
