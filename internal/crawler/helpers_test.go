package crawler

import (
	"errors"
	"testing"

	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetPageClose(t *testing.T) {
	reader := mocks.NewReadCloseMock("")

	_, err := parsePage(reader)
	if !assert.NoError(t, err) {
		return
	}

	reader.AssertCloseCalled(t)
}

func TestGetPageCloseEvenOnError(t *testing.T) {
	reader := mocks.NewErrReadCloseMock(errors.New(""))

	_, err := parsePage(reader)
	if !assert.Error(t, err) {
		return
	}

	reader.AssertCloseCalled(t)
}

func TestNewSelectorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("new selector did not panic!")
		}
	}()
	newSelector("h1..invalid")
}

func TestGetChildrenOnNil(t *testing.T) {
	assert.Empty(t, getChildrenNodes(nil))
}
