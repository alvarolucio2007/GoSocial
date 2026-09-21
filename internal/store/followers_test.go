package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFollowerStore_Follow(t *testing.T) {
	user1 := &User{Username: "testFollow", Email: "testFollow@gmail.com"}
	createUserTest(t, user1)
	user2 := &User{Username: "testFollowed", Email: "testFollowed@gmail.com"}
	createUserTest(t, user2)
	t.Cleanup(func() {
		ctx := context.Background()
		err := testStore.Users.Delete(ctx, user1.ID)
		require.NoError(t, err)
		err = testStore.Users.Delete(ctx, user2.ID)
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
			userID:     user1.ID,
			followerID: user2.ID,
			wantErr:    false,
		}, {
			name:       "invalid userID, valid followerID",
			userID:     user1.ID + 100,
			followerID: user2.ID,
			wantErr:    true,
		}, {
			name:       "valid userID, invalid followerID",
			userID:     user1.ID,
			followerID: user2.ID + 1000,
			wantErr:    true,
		}, {
			name:       "following same user",
			userID:     user1.ID,
			followerID: user1.ID,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			err := testStore.Followers.Unfollow(t.Context(), user1.ID, user2.ID)
			require.NoError(t, err)
		})
	}
	t.Run("following user twice", func(t *testing.T) {
		require.NoError(t, testStore.Followers.Follow(t.Context(), user1.ID, user2.ID))
		require.Error(t, testStore.Followers.Follow(t.Context(), user1.ID, user2.ID))
		require.NoError(t, testStore.Followers.Unfollow(t.Context(), user1.ID, user2.ID))
	})
}

func TestFollowerStore_Unfollow(t *testing.T) {
	user1 := &User{Username: "testFollow", Email: "testFollow@gmail.com"}
	createUserTest(t, user1)
	user2 := &User{Username: "testFollowed", Email: "testFollowed@gmail.com"}
	createUserTest(t, user2)
	err := testStore.Followers.Follow(t.Context(), user1.ID, user2.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx := context.Background()
		err := testStore.Users.Delete(ctx, user1.ID)
		require.NoError(t, err)
		err = testStore.Users.Delete(ctx, user2.ID)
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
			userID:     user1.ID,
			followerID: user2.ID,
			wantErr:    false,
		},
		{
			name:       "invalid userID, valid followerID",
			userID:     user1.ID + 100,
			followerID: user2.ID,
			wantErr:    true,
		},
		{
			name:       "invalid followerID, valid userID",
			userID:     user1.ID,
			followerID: user2.ID + 100,
			wantErr:    true,
		},
		{
			name:       "unfollowing same user",
			userID:     user1.ID,
			followerID: user1.ID,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := testStore.Followers.Unfollow(t.Context(), tt.userID, tt.followerID)
			if gotErr != nil {
				if !tt.wantErr {
					require.NoError(t, gotErr)
				}
				return
			}
			if tt.wantErr {
				require.Fail(t, "unfollow succeded unexpectedly")
			}
		})
	}
	t.Run("unfollowing user twice", func(t *testing.T) {
		require.NoError(t, testStore.Followers.Follow(t.Context(), user1.ID, user2.ID))
		require.NoError(t, testStore.Followers.Unfollow(t.Context(), user1.ID, user2.ID))
		require.Error(t, testStore.Followers.Unfollow(t.Context(), user1.ID, user2.ID))
	})
}
