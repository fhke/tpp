package tpp

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/utils/ptr"
)

func (t *Terraform) iterModules(iter func(source string, bl *hclwrite.Block) error) error {
	for _, bl := range t.fi.Body().Blocks() {
		if !isModuleBlock(bl) {
			continue
		}

		if len(bl.Labels()) == 0 {
			return ErrNoModuleName
		}
		modName := bl.Labels()[0]

		src := getAttrQuotedValue(bl.Body(), attrNameSource)
		if src == nil {
			return fmt.Errorf("could not parse source for module %q", modName)
		}

		if !t.isManagedModule(*src) {
			continue
		}

		if err := iter(*src, bl); err != nil {
			return fmt.Errorf("error for module %q: %w", modName, err)
		}
	}
	return nil
}

func isModuleBlock(block *hclwrite.Block) bool {
	return block.Type() == blockTypeModule
}

func getAttrQuotedValue(body *hclwrite.Body, attrName string) *string {
	attrVal := body.GetAttribute(attrName)
	if attrVal == nil {
		return nil
	}

	for _, ident := range attrVal.BuildTokens(nil) {
		if ident.Type == hclsyntax.TokenQuotedLit {
			return ptr.To(string(ident.Bytes))
		}
	}

	return nil
}

func setStringAttr(bl *hclwrite.Block, name, val string) {
	bl.Body().SetAttributeValue(name, cty.StringVal(val))
}
