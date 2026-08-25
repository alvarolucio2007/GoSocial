package store

type MockStorage struct {
	Posts     PostRepository
	Users     UserRepository
	Comments  CommentRepository
	Followers FollowerRepository
}

func NewMockStorage(mapPost map[int]*Post, mapUser map[int]*User, mapComments map[int]*Comment) MockStorage {
	return MockStorage{
		Posts:    &MockPostRepository{posts: mapPost},
		Users:    &MockUserRepository{users: mapUser},
		Comments: &MockCommentRepository{comments: mapComments},
	}
}
