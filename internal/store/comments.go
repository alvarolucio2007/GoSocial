package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type PostWithComment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}
type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}
type CommentStore struct {
	db *sql.DB
}
type CommentRepository interface {
	GetByPostID(context.Context, int64) ([]Comment, error)
	Create(context.Context, *Comment) error
}

var ErrNoContent = errors.New("no content")

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id,user_id,content)
	VALUES ($1,$2,$3)
	RETURNING id,created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	if comment.Content == "" {
		return ErrNoContent
	}
	if err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(
		&comment.ID, &comment.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username FROM comments c
								JOIN users on users.id = c.user_id
								WHERE c.post_id=$1
								ORDER BY c.created_at DESC,c.id DESC;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error in getByPostID function (comments) while closing rows:%v", err)
		}
	}()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username); err != nil {
			return nil, err
		}
		c.User.ID = c.UserID
		comments = append(comments, c)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return comments, nil
}
