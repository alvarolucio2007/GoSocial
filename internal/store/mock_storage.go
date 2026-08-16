package store

type MockStorage struct {
	Posts PostRepository
	Users UserRepository
}

func NewMockStorage(mapPost map[int]*Post, mapUser map[int]*User) MockStorage {
	return MockStorage{
		Posts: &MockPostRepository{posts: mapPost},
		Users: &MockUserRepository{users: mapUser},
	}
}
