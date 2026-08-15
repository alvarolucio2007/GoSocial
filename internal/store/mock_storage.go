package store

type MockStorage struct {
	PostRepository
	UserRepository
}

func NewMockStorage(mapPost map[int]*Post, mapUser map[int]*User) MockStorage {
	return MockStorage{
		PostRepository: &MockPostRepository{posts: mapPost},
		UserRepository: &MockUserRepository{users: mapUser},
	}
}
