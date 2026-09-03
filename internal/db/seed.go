package db

import (
	"context"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-faker/faker/v4"
)

const NUMBER_GENERATED = 100

func Seed(store store.Storage) error {
	ctx := context.Background()

	users, err := generateUsers(NUMBER_GENERATED)
	if err != nil {
		return err
	}
	var userIDs []int64
	for _, user := range users {
		if err := store.Users.Create(ctx, nil, user); err != nil {
			return err
		}
		userIDs = append(userIDs, user.ID)
	}
	posts, err := generatePosts(NUMBER_GENERATED, userIDs)
	if err != nil {
		return err
	}
	var postIDs []int64
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			return err
		}
		postIDs = append(postIDs, post.ID)
	}
	comments, err := generateComment(NUMBER_GENERATED, userIDs, postIDs)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			return err
		}
	}
	return nil
}

type fakeUsers struct {
	Username string `faker:"username"`
	Email    string `faker:"email"`
	Password string `faker:"word"` // should work, doesn't matter that much what's the password
}

func generateUsers(numUsers int) ([]*store.User, error) {
	users := make([]*store.User, numUsers)
	for i := range numUsers {
		fakeUser := fakeUsers{}
		err := faker.FakeData(&fakeUser)
		if err != nil {
			return nil, err
		}
		addSeedUser := &store.User{Username: fakeUser.Username, Email: fakeUser.Email}
		if err := addSeedUser.Password.Set(fakeUser.Password); err != nil {
			return nil, err
		}
		users[i] = addSeedUser
	}
	return users, nil
}

type fakePosts struct {
	Content string   `faker:"paragraph"`
	Title   string   `faker:"sentence"`
	UserID  int64    // gotta get this from users, shouldn't be hard
	Tags    []string // as this doesn't have a direct faker alternative i gotta make this from scratch
}

func generatePosts(numPosts int, userIDs []int64) ([]*store.Post, error) {
	posts := make([]*store.Post, numPosts)
	for i := range numPosts {
		fakePost := fakePosts{}
		err := faker.FakeData(&fakePost)
		if err != nil {
			return nil, err
		}
		fakePost.UserID = userIDs[i%len(userIDs)]
		fakePost.Tags = []string{faker.Word(), faker.Word(), faker.Word()}
		posts[i] = &store.Post{Content: fakePost.Content, Title: fakePost.Title, UserID: fakePost.UserID, Tags: fakePost.Tags}
	}
	return posts, nil
}

type fakeComment struct {
	PostID  int64  // get this from post
	UserID  int64  // get this from user
	Content string `faker:"sentence"`
}

func generateComment(numComments int, userIDs, postIDs []int64) ([]*store.Comment, error) {
	comments := make([]*store.Comment, numComments)
	for i := range numComments {
		fakeComment := fakeComment{}
		if err := faker.FakeData(&fakeComment); err != nil {
			return nil, err
		}
		fakeComment.UserID = userIDs[i%len(userIDs)]
		fakeComment.PostID = postIDs[i%len(postIDs)]
		comments[i] = &store.Comment{UserID: fakeComment.UserID, PostID: fakeComment.PostID, Content: fakeComment.Content}
	}
	return comments, nil
}
