package bundler_test

import (
	"errors"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestBundlerRunner(t *testing.T) {
	args := bundler.Args{
		Home: core.HomePage{DocsUrl: "example.com"},
	}
	bundlerRunner := bundler.Runner{Args: args}

	bundler1 := &mocks.BundlerMock{}
	bundler2 := &mocks.BundlerMock{}
	bundler3 := &mocks.BundlerMock{}

	err := bundlerRunner.Run(bundler1, bundler2, bundler3)

	require.NoError(t, err)

	require.Equal(t, args, bundler1.CalledBundleArgs)
	require.Equal(t, args, bundler2.CalledBundleArgs)
	require.Equal(t, args, bundler3.CalledBundleArgs)

}

func TestBundlerRunnerError(t *testing.T) {
	args := bundler.Args{}
	bundlerRunner := bundler.Runner{Args: args}

	expectedErr := errors.New("err")
	bundler1 := &mocks.BundlerMock{}
	bundler2 := &mocks.BundlerMock{ReturnError: expectedErr}

	err := bundlerRunner.Run(bundler1, bundler2)

	require.ErrorIs(t, err, expectedErr)

	require.Equal(t, args, bundler1.CalledBundleArgs)
	require.Equal(t, args, bundler2.CalledBundleArgs)

}
