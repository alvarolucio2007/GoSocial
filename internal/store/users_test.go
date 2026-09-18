package store

import (
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestCreateUser(t *testing.T) {
	user := &User{Username: "testCreate", Email: "test@gmail.com"}
	t.Run("user creation", func(t *testing.T) {
		tx, err := testDB.Begin()
		require.NoError(t, err)
		err = testStore.Users.Create(t.Context(), tx, user)
		require.NoError(t, err)
		err = tx.Commit()
		require.NoError(t, err)
	})
	t.Run("read created user to check whether user was truly created", func(t *testing.T) {
		_, err := testDB.Exec("UPDATE users SET is_active = true WHERE username='testCreate'")
		require.NoError(t, err)

		userReceived, err := testStore.Users.GetByEmail(t.Context(), user.Email)
		require.NoError(t, err)
		require.Equal(t, user.Username, userReceived.Username)
		require.Equal(t, user.Email, userReceived.Email)
	})
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE email='test@gmail.com'")
		require.NoError(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	user := &User{Username: "testUpdate", Email: "testUpdate@gmail.com"}
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(t.Context(), tx, user)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testUpdate'")
	require.NoError(t, err)
	t.Run("user update", func(t *testing.T) {
		err = testStore.Users.Update(t.Context(), user)
		require.NoError(t, err)
		userReceived, err := testStore.Users.GetByEmail(t.Context(), user.Email)
		require.NoError(t, err)
		require.Equal(t, user.Username, userReceived.Username)
		require.Equal(t, user.Email, userReceived.Email)
	})
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE email='testUpdate@gmail.com'")
		require.NoError(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	user := &User{Username: "testDelete", Email: "testDelete@gmail.com"}
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(t.Context(), tx, user)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testDelete'")
	require.NoError(t, err)
	var userID int64
	err = testDB.QueryRow("SELECT id FROM users WHERE email='testDelete@gmail.com'").Scan(&userID)
	require.NoError(t, err)
	t.Run("testing user deletion", func(t *testing.T) {
		err = testStore.Users.Delete(t.Context(), userID)
		require.NoError(t, err)
		user, err := testStore.Users.Read(t.Context(), userID)
		require.ErrorIs(t, err, ErrUserNotFound)
		require.Nil(t, user)
	})
}
