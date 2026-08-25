package main

import (
	"log"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/db"
	"github.com/alvarolucio2007/GoSocial/internal/env"
	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func main() {
	addr := env.GetString("DB_URL", "postgres://admin:admin@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 30, 30, 15*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error while seeding: %v", err)
		}
	}()
	store := store.NewPostgresStorage(conn)
	if err := db.Seed(store); err != nil {
		log.Printf("error while seeding: %v", err)
	}
}
