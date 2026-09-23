package cache

import (
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func TestRedisUserCache_Set(t *testing.T) {
	user := &store.User{ID: 1, Username: "TestSetCache", Email: "TestSetCache@gmail.com"}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		user    *store.User
		wantErr bool
	}{
		{
			"adding normal user",
			user,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := testCache.Users.Set(t.Context(), tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Set() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Set() succeeded unexpectedly")
			}
		})
	}
}
