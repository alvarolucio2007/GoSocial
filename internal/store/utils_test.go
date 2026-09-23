package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/testify/require"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB    *sql.DB
	testStore Storage
)

func TestMain(m *testing.M) {
	log.Println("testMain started")
	ctx := context.Background()

	var testContainer *postgres.PostgresContainer
	var err error
	testDB, testContainer, err = setupTestDB()
	if err != nil {
		log.Fatalf("unable to create DB and testcontainers %v", err)
	}
	if testDB == nil {
		log.Fatalf("testDB is nil")
	}
	if testContainer == nil {
		log.Fatalf("testDB is nil")
	}
	testStore = NewPostgresStorage(testDB)

	code := m.Run()
	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
	if err := testDB.Close(); err != nil {
		log.Printf("failed to terminate DB: %v", err)
	}
	os.Exit(code)
}

func setupTestDB() (*sql.DB, *postgres.PostgresContainer, error) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("gosocial_test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(connStr); err != nil {
		return nil, nil, err
	}

	return db, pgContainer, nil
}

func runMigrations(connStr string) error {
	m, err := migrate.New("file:../../cmd/migrate/migrations/", connStr)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}

func createUserTest(t *testing.T, user *User) {
	tx, err := testDB.Begin()
	require.NoError(t, err)
	err = testStore.Users.Create(t.Context(), tx, user)
	require.NoError(t, err)
	_, err = tx.Exec("UPDATE users SET is_active = true WHERE username=$1", user.Username)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
	err = testDB.QueryRowContext(t.Context(), "SELECT id FROM users WHERE email=$1", user.Email).Scan(&user.ID)
	require.NoError(t, err)
}

func createPostTest(t *testing.T, post *Post) {
	err := testStore.Posts.Create(t.Context(), post)
	require.NoError(t, err)
	require.NotZero(t, post.ID, "create not populating post.ID")
	require.False(t, post.CreatedAt.IsZero(), "create not populating post.createdAt")
}
