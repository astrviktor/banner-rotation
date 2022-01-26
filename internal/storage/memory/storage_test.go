package memorystorage

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("not implemented", func(t *testing.T) {
		require.NoError(t, nil)
		require.True(t, true)
		require.Equal(t, 0, 0)
	})
}
