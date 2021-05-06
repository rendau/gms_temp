package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	prepareDbForNewTest()

	require.True(t, true, true)
}
