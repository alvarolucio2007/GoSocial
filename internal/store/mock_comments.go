package store

import "context"

type MockCommentRepository struct {
	comments map[int]*Comment
}

func (m *MockCommentRepository) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	return nil, nil
}
