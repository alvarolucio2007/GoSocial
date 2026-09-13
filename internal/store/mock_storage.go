package store

func NewMockStorage(mapPost map[int]*Post, mapUser map[int]*User, mapComments map[int]*Comment, mapFollowers map[int]*Follower, mapRoles map[int]*Role) Storage {
	return Storage{
		Posts:     &MockPostRepository{posts: mapPost},
		Users:     &MockUserRepository{users: mapUser},
		Comments:  &MockCommentRepository{comments: mapComments},
		Followers: &MockFollowerRepository{followers: mapFollowers},
		Roles:     &MockRoleRepository{roles: mapRoles},
	}
}
