package core_test

import (
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/stretchr/testify/require"
)

func TestRenderFeatureDescription(t *testing.T) {
	feature := core.Feature{
		Description: "Feature description",
	}
	t.Run("with feature", func(t *testing.T) {
		description := core.RenderFeatureDescription(&feature, core.OpenAPI)
		require.Contains(t, description, feature.Description)
		require.Contains(t, description, "https://github.com/giuliano-macedo/omie-docs")
		require.Contains(t, description, "não oficial da API da Omie")
	})
	t.Run("without feature", func(t *testing.T) {
		description := core.RenderFeatureDescription(nil, core.OpenAPI)
		require.NotContains(t, description, feature.Description)
		require.Contains(t, description, "não oficial da API da Omie")
	})
	t.Run("openAPI", func(t *testing.T) {
		description := core.RenderFeatureDescription(nil, core.OpenAPI)
		require.Contains(t, description, "https://github.com/OAI/OpenAPI-Specification/issues/1635")
		require.Contains(t, description, "OpenAPI")
	})

	t.Run("postman", func(t *testing.T) {
		description := core.RenderFeatureDescription(nil, core.Postman)
		require.NotContains(t, description, "https://github.com/OAI/OpenAPI-Specification/issues/1635")
		require.Contains(t, description, "Postman")
	})
	t.Run("bruno", func(t *testing.T) {
		description := core.RenderFeatureDescription(nil, core.Bruno)
		require.NotContains(t, description, "https://github.com/OAI/OpenAPI-Specification/issues/1635")
		require.Contains(t, description, "Bruno")
	})
}
