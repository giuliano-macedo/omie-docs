package bundler_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getDirTemp() string {
	dirTemp, err := os.MkdirTemp("", "os_file_writer_test")
	if err != nil {
		panic(err)
	}
	return dirTemp
}

var dirTemp = getDirTemp()

func TestOsFileWriter(t *testing.T) {
	defer os.RemoveAll(dirTemp)
	getFileWriter := func(prettifyJson bool) bundler.FSWriter {
		fileWriter, err := bundler.NewOsFileWriter(path.Join(dirTemp, "bundle"), prettifyJson)
		if err != nil {
			panic(err)
		}
		return fileWriter
	}

	t.Run("Create", func(t *testing.T) {
		fileWriter := getFileWriter(false)
		data := []byte("hello world")

		writer, err := fileWriter.Create("foo.txt")
		require.NoError(t, err)

		_, err = writer.Write(data)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		dataRead, err := os.ReadFile(path.Join(dirTemp, "bundle", "foo.txt"))
		require.NoError(t, err)
		require.Equal(t, data, dataRead)
	})
	t.Run("SaveJson", func(t *testing.T) {
		runSaveJson := func(t *testing.T, prettifyJson bool) {
			fileWriter := getFileWriter(prettifyJson)
			data := map[string]interface{}{
				"Hello": "world",
				"n":     float64(2),
			}

			require.NoError(t, fileWriter.SaveJsonFile("test.json", data))

			file, err := os.Open(path.Join(dirTemp, "bundle", "test.json"))
			require.NoError(t, err)
			decoder := json.NewDecoder(file)

			var dataRead map[string]interface{}
			decoder.Decode(&dataRead)
			require.Equal(t, data, dataRead)

			require.NoError(t, file.Close())
		}
		t.Run("Without prettifyJson", func(t *testing.T) { runSaveJson(t, false) })
		t.Run("With prettifyJson", func(t *testing.T) { runSaveJson(t, true) })
	})
	t.Run("mkDirAll", func(t *testing.T) {
		fileWriter := getFileWriter(false)
		err := fileWriter.MkdirAll("dir1/dir2")
		if !assert.NoError(t, err) {
			t.Fail()
		}

		stat, err := os.Stat(path.Join(dirTemp, "bundle", "dir1/dir2"))
		if !assert.NoError(t, err) {
			t.Fail()
		}

		if !assert.True(t, stat.IsDir()) {
			t.Fail()
		}
	})
}

func TestOsFileWriterErr(t *testing.T) {
	invalidPath := string([]byte{'a', 0, 'b'})
	t.Run("On New", func(t *testing.T) {
		_, err := bundler.NewOsFileWriter(invalidPath, false)
		require.Error(t, err)
	})
	t.Run("On saveJson", func(t *testing.T) {
		fsWriter, err := bundler.NewOsFileWriter(dirTemp, false)
		require.NoError(t, err)

		require.Error(t, fsWriter.SaveJsonFile(invalidPath, "hello"))
	})
}
