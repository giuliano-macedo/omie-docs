package openapi_conversor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) { assert.Equal(t, *newPtr("hello"), "hello") })
	t.Run("int", func(t *testing.T) { assert.Equal(t, *newPtr(1), 1) })
}

func TestNewInterface(t *testing.T) {
	assert.Equal(t, (*newInterface("hello")).(string), "hello")
}

func TestGetFieldRef(t *testing.T) {
	t.Run("import", func(t *testing.T) {
		assert.Equal(t, getFieldRef("endpoint/test", "HelloWorld", true), "#/components/schemas/endpoint.test.HelloWorld")
	})
	t.Run("export", func(t *testing.T) {
		assert.Equal(t, getFieldRef("endpoint/test", "HelloWorld", false), "endpoint.test.HelloWorld")
	})
}
