package test

import (
	"path/filepath"
	"testing"

	"github.com/fhke/tpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

func TestTerraform__OK(t *testing.T) {
	goldenPath := "fixtures/single-file/golden.tf"
	outFile := filepath.Join(t.TempDir(), "out.tf")

	tf, err := tpp.NewTerraformForFile("fixtures/single-file/input.tf")
	require.NoError(t, err, "It should parse file")

	modSrcs, err := tf.GetModuleSources()
	require.NoError(t, err, "It should get module sources")
	assert.Equal(
		t,
		sets.New("s3::somebucket/a/sdf", "s3::otherbucket/a/sdf"),
		sets.New(modSrcs...),
		"Module sources should be correct",
	)

	err = tf.SetModuleSources(map[string]tpp.ModuleSource{
		"s3::somebucket/a/sdf": {
			Source: "somebucket/replaced",
		},
		"s3::otherbucket/a/sdf": {
			Source:  "somebucket/replaced/2",
			Version: ptr.To("1.0.0"),
		},
	})
	require.NoError(t, err, "It should update module sources")

	wrLen, err := tf.WriteFile(outFile)
	require.NoError(t, err, "It should write file")
	assert.Equal(t, len(tf.Bytes()), int(wrLen), "It should write correct number of bytes")

	assertFilesMatch(t, goldenPath, outFile)
}
