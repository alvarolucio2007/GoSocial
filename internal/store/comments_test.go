package store

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestCommentStore_Create(t *testing.T) {
	user := &User{Username: "testCommentUser", Email: "testCommentUser@gmail.com"}
	createUserTest(t, user)

	post := &Post{UserID: user.ID, Content: "test post"}
	createPostTest(t, post)

	t.Cleanup(func() {
		ctx := context.Background()
		_, err := testDB.ExecContext(ctx, "DELETE FROM comments WHERE user_id = $1", user.ID)
		require.NoError(t, err)
		_, err = testDB.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", post.ID)
		require.NoError(t, err)
		require.NoError(t, testStore.Users.Delete(ctx, user.ID))
	})

	tests := []struct {
		name    string
		comment *Comment
		wantErr bool
	}{
		{
			name: "valid comment",
			comment: &Comment{
				PostID:  post.ID,
				UserID:  user.ID,
				Content: "hello world",
			},
			wantErr: false,
		},
		{
			name: "invalid postID",
			comment: &Comment{
				PostID:  post.ID + 99999,
				UserID:  user.ID,
				Content: "orphan comment",
			},
			wantErr: true,
		},
		{
			name: "invalid userID",
			comment: &Comment{
				PostID:  post.ID,
				UserID:  user.ID + 99999,
				Content: "ghost user",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			comment: &Comment{
				PostID:  post.ID,
				UserID:  user.ID,
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testStore.Comments.Create(t.Context(), tt.comment)

			if tt.wantErr {
				require.Errorf(t, err, "required error in %v function", tt.name)
				return
			}
			if err != nil {
				require.NoErrorf(t, err, "required no error in %v function", tt.name)
			}

			require.NotZero(t, tt.comment.ID, "Create should populate ID via RETURNING")
			require.Equal(t, tt.comment.PostID, tt.comment.PostID)
			require.Equal(t, tt.comment.UserID, tt.comment.UserID)
		})
	}
}

func TestCommentStore_GetByPostID(t *testing.T) {
	user1 := &User{Username: "commentAuthor1", Email: "commentAuthor1@gmail.com"}
	createUserTest(t, user1)

	user2 := &User{Username: "commentAuthor2", Email: "commentAuthor2@gmail.com"}
	createUserTest(t, user2)

	post := &Post{UserID: user1.ID, Content: "post with comments"}
	createPostTest(t, post)

	otherPost := &Post{UserID: user1.ID, Content: "post without comments"}
	createPostTest(t, otherPost)

	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = testDB.ExecContext(ctx, "DELETE FROM comments WHERE post_id = $1", post.ID)
		_, _ = testDB.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", post.ID)
		_, _ = testDB.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", otherPost.ID)
		require.NoError(t, testStore.Users.Delete(ctx, user1.ID))
		require.NoError(t, testStore.Users.Delete(ctx, user2.ID))
	})

	c1 := &Comment{PostID: post.ID, UserID: user1.ID, Content: "first"}
	c2 := &Comment{PostID: post.ID, UserID: user2.ID, Content: "second"}
	c3 := &Comment{PostID: post.ID, UserID: user1.ID, Content: "third"}

	require.NoError(t, testStore.Comments.Create(t.Context(), c1))
	require.NoError(t, testStore.Comments.Create(t.Context(), c2))
	require.NoError(t, testStore.Comments.Create(t.Context(), c3))
	tests := []struct {
		name      string
		postID    int64
		wantCount int
		wantErr   bool
	}{
		{
			name:      "post with 3 comments",
			postID:    post.ID,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "post with no comments",
			postID:    otherPost.ID,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "nonexistent post",
			postID:    post.ID + 99999,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testStore.Comments.GetByPostID(t.Context(), tt.postID)

			if tt.wantErr {
				require.Errorf(t, err, "required error in %v function", tt.name)
				return
			}
			require.NoError(t, err)

			require.Len(t, got, tt.wantCount)

			if tt.wantCount == 3 {
				require.Equal(t, c3.Content, got[0].Content)
				require.Equal(t, c2.Content, got[1].Content)
				require.Equal(t, c1.Content, got[2].Content)

				require.Equal(t, user2.Username, got[1].User.Username)
				require.Equal(t, user2.ID, got[1].User.ID)
			}
		})
	}
}
