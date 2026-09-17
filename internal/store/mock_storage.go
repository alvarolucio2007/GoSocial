package store

func NewMockStorage(mapPost map[int64]*Post, mapUser map[int64]*User, mapComments map[int64]*Comment, mapFollowers map[FollowKey]struct{}, mapRoles map[int]*Role) Storage {
	return Storage{
		Posts:     &MockPostRepository{posts: mapPost, users: mapUser},
		Users:     &MockUserRepository{users: mapUser},
		Comments:  &MockCommentRepository{comments: mapComments, posts: mapPost},
		Followers: &MockFollowerRepository{followers: mapFollowers},
		Roles:     &MockRoleRepository{roles: mapRoles},
	}
}
