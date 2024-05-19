package openapi_conversor

import (
	_ "embed"
	"fmt"
	"strings"
)

func newPtr[T any](value T) (ptr *T) {
	ptr = new(T)
	*ptr = value
	return ptr
}

func newInterface(value any) (ptr *interface{}) {
	return newPtr(value)
}

func getFieldRef(endpoint string, modelName string, isImport bool) string {
	schemaRefPrefix := strings.ReplaceAll(endpoint, "/", ".")
	if isImport {
		return fmt.Sprint("#/components/schemas/", schemaRefPrefix, ".", modelName)
	} else {
		return fmt.Sprint(schemaRefPrefix, ".", modelName)
	}
}
