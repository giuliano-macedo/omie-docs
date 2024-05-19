package testdata

import (
	"encoding/json"
	"os"
	"path"
)

func getFilePath(caseName, fileName string) string {
	return path.Join("testdata", caseName, fileName)
}

func ReadFile(caseName, fileName string) string {
	data, err := os.ReadFile(getFilePath(caseName, fileName))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func ReadJsonFile[T any](caseName, fileName string) (res T) {
	file, err := os.Open(getFilePath(caseName, fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		panic(err)
	}
	return res
}
