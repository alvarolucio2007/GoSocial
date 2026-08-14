package main

import (
	"log"

	"github.com/alvarolucio2007/GoSocial/internal/env"
	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	store := store.NewPostgresStorage(nil)
	app := &application{
		cfg,
		store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
