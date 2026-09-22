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
