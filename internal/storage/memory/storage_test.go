package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("not implemented", func(t *testing.T) {
		require.NoError(t, nil)
		require.True(t, true)
		require.Equal(t, 0, 0)
	})
}
