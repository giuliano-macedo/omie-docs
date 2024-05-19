package bruno_test

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bruno/test/testdata"
	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestBrunoBundler(t *testing.T) {
	brunoBundler := bruno.NewBrunoBundler()
	home := testdata.ReadJsonFile[core.HomePage]("normal", "homepage.json")
	pages := testdata.ReadJsonFile[[]core.Page]("normal", "pages.json")

	fsWriter := mocks.NewFilewriterMock()

	err := brunoBundler.Bundle(bundler.Args{
		Pages:    pages,
		Home:     home,
		FsWriter: fsWriter,
	})
	require.NoError(t, err)

	require.True(t, fsWriter.AllFilesClosed())
	require.Equal(t, fsWriter.FileNamesWritten(), []string{core.BrunoCollectionName})

	writenFile := fsWriter.FileWritten(core.BrunoCollectionName)
	require.NotNil(t, writenFile)

	assertIsValidTarFile(t, writenFile.FileContent)

}

func assertIsValidTarFile(t *testing.T, data []byte) {
	reader := tar.NewReader(bytes.NewReader(data))

	for {
		_, err := reader.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		_, err = io.ReadAll(reader)
		require.NoError(t, err)
	}
}
