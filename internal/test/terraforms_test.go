package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fhke/tpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

func TestTerraforms__OK(t *testing.T) {
	goldenDir := "fixtures/multi-file/golden"
	outDir := t.TempDir()

	tf, err := tpp.NewTerraformsForDir("fixtures/multi-file/input")
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

	require.NoError(t, tf.WriteTo(outDir), "It should write data")
	assertDirsMatch(t, goldenDir, outDir)
}

func assertDirsMatch(t *testing.T, src, dst string) {
	srcEnts, err := os.ReadDir(src)
	require.NoError(t, err, "It should read source dir")
	dstEnts, err := os.ReadDir(dst)
	require.NoError(t, err, "It should read destination dir")

	assert.Equal(t, len(srcEnts), len(dstEnts), "Directories should have same number of files")
	for _, srcEnt := range srcEnts {
		if srcEnt.IsDir() {
			t.Fatal("Asserting directories is not supported")
		}
		assertFilesMatch(t, filepath.Join(src, srcEnt.Name()), filepath.Join(dst, srcEnt.Name()))
	}
}

func assertFilesMatch(t *testing.T, src, dst string) {
	srcData, err := os.ReadFile(src)
	require.NoError(t, err, "It should read source file")
	dstData, err := os.ReadFile(dst)
	require.NoError(t, err, "It should read destination file")

	assert.Equalf(t, string(srcData), string(dstData), "Data in file %q should match file %q", dst, src)
}
