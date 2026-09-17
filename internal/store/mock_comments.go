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
	return post.Comments, nil
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *Comment) error {
	if _, exists := m.comments[comment.ID]; exists {
		return ErrConflict
	}
	m.comments[comment.ID] = comment
	return nil
}
