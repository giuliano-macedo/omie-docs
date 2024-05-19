package openapi_bundler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatFeatureName(t *testing.T) {
	testCases := []struct {
		featureName                  string
		exptectedFormatedFeatureName string
	}{
		{"sim, um cachorro de caça idealmente lá", "sim_um_cachorro_de_caca_idealmente_la"},
		{"Finanças", "financas"},
		{"Compras, Estoque e Produção", "compras_estoque_e_producao"},
		{"Number 1", "number_1"},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprint("case ", i), func(t *testing.T) {
			require.Equal(t, formatFeatureName(testCase.featureName), testCase.exptectedFormatedFeatureName)
		})
	}
}

func TestRemoveJsComments(t *testing.T) {
	sourceCode := `
	// This is a single-line comment
	var x = 10; // Another comment here

	/*
	This is a multi-line comment
	*/
	function test() { /* another comment */
		// Function body
	}`
	expectedJsWithouctComments := `
	
	var x = 10; 

	
	function test() { 
		
	}`
	jsWithoutComments := removeJsComments(sourceCode)

	assert.Equal(t, jsWithoutComments, expectedJsWithouctComments)
}
