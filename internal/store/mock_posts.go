package store

import "context"

type MockPostRepository struct {
	posts map[int]*Post
}

func (m *MockPostRepository) Create(ctx context.Context, post *Post) error {
	m.posts[int(post.ID)] = post
	return nil
}

func (m *MockPostRepository) Read(ctx context.Context, id int) (*Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return p, nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *Post) error {
	oldPost, err := m.Read(ctx, int(post.ID))
	if err != nil {
		return ErrPostNotFound // literally the only error that can be thrown by m.Read
	}
	p := *oldPost
	if post.Content != "" {
		p.Content = post.Content
	}
	if post.Title != "" {
		p.Title = post.Title
	}
	if post.Tags != nil {
		p.Tags = post.Tags
	}

	m.posts[int(post.ID)] = &p
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, idPost int) error {
	_, err := m.Read(ctx, idPost)
	if err != nil {
		return ErrPostNotFound
	}
	delete(m.posts, idPost)
	return nil
}

func (m *MockPostRepository) GetUserFeed(ctx context.Context, idUser int64, fn PaginatedFeedQuery) ([]PostWithMetadata, error) {
	return nil, nil
}
