package store

import (
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    PaginatedFeedQuery
		wantErr bool
	}{
		{
			name:  "default when empty",
			query: "",
			want:  PaginatedFeedQuery{Sort: "desc"},
		},
		{
			name:  "valid limit and offset",
			query: "?limit=20&offset=40",
			want:  PaginatedFeedQuery{Limit: 20, Offset: 40, Sort: "desc"},
		},
		{
			name:  "sort asc",
			query: "?sort=ASC",
			want:  PaginatedFeedQuery{Sort: "asc"},
		},
		{
			name:  "tags separated by comma",
			query: "?tags=go,backend",
			want:  PaginatedFeedQuery{Sort: "desc", Tags: []string{"go", "backend"}},
		},
		{
			name:    "invalid limit",
			query:   "?limit=abc",
			wantErr: true,
		},
		{
			name:    "negative limit",
			query:   "?limit=-1",
			wantErr: true,
		},
		{
			name:    "invalid sort",
			query:   "?sort=random",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/feed"+tt.query, nil)
			var fq PaginatedFeedQuery
			got, err := fq.Parse(req)

			if tt.wantErr {
				require.Errorf(t, err, "required error in %v function", tt.name)
				return
			}
			if err != nil {
				require.NoErrorf(t, err, "required no error in %v function", tt.name)
			}
			require.Equal(t, tt.want.Limit, got.Limit)

			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Sort, got.Sort)
			require.EqualValues(t, tt.want.Tags, got.Tags)
		})
	}
}
