package store

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestFollowerStore_Follow(t *testing.T) {
	user1 := &User{Username: "testFollow", Email: "testFollow@gmail.com"}
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(t.Context(), tx, user1)
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testFollow'")
	require.NoError(t, err)
	user2 := &User{Username: "testFollowed", Email: "testFollowed@gmail.com"}
	err = testStore.Users.Create(t.Context(), tx, user2)
	require.NoError(t, err)
	_, err = testDB.Exec("UPDATE users SET is_active = true WHERE username='testFollowed'")
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx := context.Background()
		err = testStore.Users.Delete(ctx, 1)
		require.NoError(t, err)
		err = testStore.Users.Delete(ctx, 2)
		require.NoError(t, err)
	})
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID     int64
		followerID int64
		wantErr    bool
	}{
		{
			name:       "valid userID and followerID",
			userID:     1,
			followerID: 2,
			wantErr:    false,
		}, {
			name:       "invalid userID, valid followerID",
			userID:     100,
			followerID: 2,
			wantErr:    true,
		}, {
			name:       "invalid userID, valid followerID",
			userID:     999,
			followerID: 2,
			wantErr:    true,
		}, {
			name:       "valid userID, invalid followerID",
			userID:     1,
			followerID: 999,
			wantErr:    true,
		}, {
			name:       "following same user",
			userID:     1,
			followerID: 1,
			wantErr:    true,
		}, {
			name:       "following user twice",
			userID:     1,
			followerID: 2,
			wantErr:    true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == 5 {
				gotErr := testStore.Followers.Follow(t.Context(), tt.userID, tt.followerID)
				require.NoError(t, gotErr)
			}
			gotErr := testStore.Followers.Follow(t.Context(), tt.userID, tt.followerID)
			if gotErr != nil {
				if !tt.wantErr {
					require.NoError(t, gotErr)
				}
				return
			}
			if tt.wantErr {
				require.Fail(t, "follow succeded unexpectedly")
			}
			err := testStore.Followers.Unfollow(t.Context(), 1, 2)
			require.NoError(t, err)
		})
	}
}
