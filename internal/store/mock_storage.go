package store

type MockStorage struct {
	PostRepository
	UserRepository
}

func NewMockStorage() MockStorage {
	return MockStorage{
		PostRepository: nil, // i have to add soon, it'll take a tiny bit of time
		UserRepository: nil,
	}
}
