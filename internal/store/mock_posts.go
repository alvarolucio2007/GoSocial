package store

import "context"

type MockPostRepository struct {
	posts map[int64]*Post
	users map[int64]*User
}

func (m *MockPostRepository) Create(ctx context.Context, post *Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Read(ctx context.Context, id int64) (*Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return p, nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *Post) error {
	oldPost, err := m.Read(ctx, post.ID)
	if err != nil {
		return ErrPostNotFound
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

	m.posts[post.ID] = &p
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, idPost int64) error {
	_, err := m.Read(ctx, idPost)
	if err != nil {
		return ErrPostNotFound
	}
	delete(m.posts, idPost)
	return nil
}

func (m *MockPostRepository) GetUserFeed(ctx context.Context, idUser int64, fn PaginatedFeedQuery) ([]PostWithMetadata, error) {
	if _, exist := m.users[idUser]; !exist {
		return nil, ErrUserNotFound
	}
	// TODO: Implement this somehow.
}
