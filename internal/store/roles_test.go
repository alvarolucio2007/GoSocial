package store

import (
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestGetByName(t *testing.T) {
	t.Run("select all 3 roles", func(t *testing.T) {
		roleAdmin, err := testStore.Roles.GetByName(t.Context(), "admin")
		require.NoError(t, err)
		require.NotNil(t, roleAdmin)
		roleModerator, err := testStore.Roles.GetByName(t.Context(), "moderator")
		require.NoError(t, err)
		require.NotNil(t, roleModerator)

		roleUser, err := testStore.Roles.GetByName(t.Context(), "user")
		require.NoError(t, err)
		require.NotNil(t, roleUser)
	})
	t.Run("select non-existent role", func(t *testing.T) {
		roleUnknown, err := testStore.Roles.GetByName(t.Context(), "test")
		require.Error(t, err)
		require.Nil(t, roleUnknown)
	})
}
