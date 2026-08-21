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
	defer conn.Close()
	store := store.NewPostgresStorage(conn)
	db.Seed(store)
}
