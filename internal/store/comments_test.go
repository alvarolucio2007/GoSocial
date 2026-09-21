package store_test

import (
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func TestCommentStore_Create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		comment *store.Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s store.CommentStore
			gotErr := s.Create(t.Context(), tt.comment)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
		})
	}
}
