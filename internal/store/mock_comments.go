package store

import (
	"context"
)

type MockCommentRepository struct {
	comments map[int64]*Comment
	posts    map[int64]*Post
}

func (m *MockCommentRepository) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	post, exists := m.posts[postID]
	if !exists {
		return nil, ErrNotFound
	}
	result := make([]Comment, len(post.Comments))
	for c := range post.Comments {
		result = append(result, *m.comments[int64(c)])
	}
	return result, nil
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *Comment) error {
	_, exists := m.comments[comment.ID]
	if !exists {
		return ErrConflict
	}
	m.comments[comment.ID] = comment
	return nil
}
