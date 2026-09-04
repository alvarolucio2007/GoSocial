package main

import (
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/db"
	"github.com/alvarolucio2007/GoSocial/internal/env"
	"github.com/alvarolucio2007/GoSocial/internal/store"
	_ "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const version = "0.0.1"

//	@title			GoSocial API
//	@description	This is a sample server GoSocial server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	MIT
//	@license.url	https://mit-license.org/

// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg_zap := zap.NewProductionConfig()

	logger := zap.Must(cfg_zap.Build(
		zap.AddStacktrace(zapcore.FatalLevel),
	)).Sugar()

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Errorf("Couldn't sync logger: %s", err)
		}
	}()

	maxOpenConn, _ := env.GetInt("DB_MAX_OPEN_CONN", 30)
	maxIdleConn, _ := env.GetInt("DB_MAX_IDLE_CONN", 30)
	maxIdleTime, err := time.ParseDuration(env.GetString("DB_MAX_IDLE_TIME", "15m"))
	if err != nil {
		logger.Errorf("couldn't parse maxIdleTime, error: %v", err) // might want to deal w this later oh well
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://admin:admin@localhost/social?sslmode=disable"),
			maxOpenConn: maxOpenConn,
			maxIdleConn: maxIdleConn,
			maxIdleTime: maxIdleTime,
		},
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		mail: mailConfig{
			exp: 3 * 24 * time.Hour, // 3 days
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConn, cfg.db.maxIdleConn, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panicf("PANIC: couldn't connect to database, error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("couldn't close the database: %v\n", err)
		}
	}()
	store := store.NewPostgresStorage(db)

	app := &application{
		config:  cfg,
		storage: store,
		logger:  logger,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
