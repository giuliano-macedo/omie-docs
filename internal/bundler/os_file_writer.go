package bundler

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

type OsFileWriter struct {
	basePath     string
	prettifyJson bool
}

func NewOsFileWriter(basePath string, prettifyJson bool) (FSWriter, error) {
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &OsFileWriter{basePath: basePath, prettifyJson: prettifyJson}, nil
}

func (fsWriter *OsFileWriter) MkdirAll(name string) error {
	return os.MkdirAll(path.Join(fsWriter.basePath, name), os.ModePerm)
}

func (fsWriter *OsFileWriter) Create(name string) (io.WriteCloser, error) {
	return os.Create(path.Join(fsWriter.basePath, name))
}

func (fsWriter *OsFileWriter) SaveJsonFile(name string, value interface{}) error {
	file, err := fsWriter.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	if fsWriter.prettifyJson {
		encoder.SetIndent("", "    ")
	}
	return encoder.Encode(value)
}
