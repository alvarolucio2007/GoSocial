package main

import (
	"log"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/db"
	"github.com/alvarolucio2007/GoSocial/internal/env"
	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func main() {
	maxOpenConn, _ := env.GetInt("DB_MAX_OPEN_CONN", 30)
	maxIdleConn, _ := env.GetInt("DB_MAX_IDLE_CONN", 30)
	maxIdleTime, err := time.ParseDuration(env.GetString("DB_MAX_IDLE_TIME", "15m"))
	if err != nil {
		log.Panicf("PANIC: couldn't parse maxIdleTime, error: %v", err) // might want to deal w this later oh well
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://admin:admin@localhost/social?sslmode=disable"),
			maxOpenConn: maxOpenConn,
			maxIdleConn: maxIdleConn,
			maxIdleTime: maxIdleTime,
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConn, cfg.db.maxIdleConn, cfg.db.maxIdleTime)
	if err != nil {
		log.Panicf("PANIC: couldn't connect to database, error: %v", err)
	}
	defer db.Close()
	store := store.NewPostgresStorage(db)
	app := &application{
		cfg,
		store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
