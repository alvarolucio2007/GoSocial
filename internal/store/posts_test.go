package store

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestCreatePost(t *testing.T) {
	user := &User{Username: "testCreatePost", Email: "testCreatePost@gmail.com"}
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(context.Background(), tx, user)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testCreate'")
	require.NoError(t, err)
	t.Run("create valid post with valid userID", func(t *testing.T) {
		post := &Post{Content: "testContent", Title: "testTitle", UserID: user.ID, Tags: nil}
		err := testStore.Posts.Create(t.Context(), post)
		require.NoError(t, err)

		var postID int64
		err = testDB.QueryRow("SELECT id FROM posts WHERE content='testContent'").Scan(&postID)
		require.NoError(t, err)
		postRead, err := testStore.Posts.Read(t.Context(), postID)
		require.NoError(t, err)
		require.Equal(t, post.Content, postRead.Content)
		require.Equal(t, post.Title, postRead.Title)
		require.Equal(t, post.UserID, postRead.UserID)
	})
	t.Run("create invalid post, invalid userID", func(t *testing.T) {
		post := &Post{Content: "testContent", Title: "testTitle", UserID: user.ID - 1, Tags: nil}
		err := testStore.Posts.Create(t.Context(), post)
		require.Error(t, err)
	})
	t.Cleanup(func() {
		_, err = testDB.Exec("DELETE FROM posts WHERE id=1")
		require.NoError(t, err)
		_, err := testDB.Exec("DELETE FROM users WHERE email='testCreatePost@gmail.com'")
		require.NoError(t, err)
	})
}

func TestUpdatePost(t *testing.T) {
	user := &User{Username: "testUpdatePost", Email: "testUpdatePost@gmail.com"}
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(context.Background(), tx, user)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testUpdatePost'")
	require.NoError(t, err)

	post := &Post{Content: "testContent", Title: "testTitle", UserID: user.ID, Tags: nil}
	err = testStore.Posts.Create(t.Context(), post)
	require.NoError(t, err)

	var postID int64
	err = testDB.QueryRow("SELECT id FROM posts WHERE content='testContent'").Scan(&postID)
	require.NoError(t, err)
	t.Run("edit valid post", func(t *testing.T) {
		postUpdated := &Post{ID: postID, Content: "updatedContent", Title: "updatedTitle", Tags: []string{"Test1", "Test2"}}
		err := testStore.Posts.Update(t.Context(), postUpdated)
		require.NoError(t, err)

		postRead, err := testStore.Posts.Read(t.Context(), postID)
		require.NoError(t, err)
		require.NotNil(t, postRead)

		require.Equal(t, postUpdated.Content, postRead.Content)
		require.Equal(t, postUpdated.Title, postRead.Title)
		require.Equal(t, postUpdated.Tags, postRead.Tags)
	})
	t.Run("edit a post which doesn't exist", func(t *testing.T) {
		postUpdated := &Post{ID: postID + 1, Content: "updatedContent", Title: "updatedTitle", Tags: []string{"Test1", "Test2"}}
		err := testStore.Posts.Update(t.Context(), postUpdated)
		require.Error(t, err)
	})
}

func TestDeletePost(t *testing.T) {
}
