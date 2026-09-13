package main

import (
	"encoding/base64"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/alvarolucio2007/GoSocial/internal/db"
	"github.com/alvarolucio2007/GoSocial/internal/env"
	"github.com/alvarolucio2007/GoSocial/internal/mailer"
	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/alvarolucio2007/GoSocial/internal/store/cache"
	"github.com/redis/go-redis/v9"
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

	maxIdleTime, err := time.ParseDuration(env.GetString("DB_MAX_IDLE_TIME", "15m"))
	if err != nil {
		logger.Errorf("couldn't parse maxIdleTime, error: %v", err) // might want to deal w this later oh well
	}
	secretB64 := env.GetString("AUTH_TOKEN_SECRET", "")
	secret, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		logger.Panicw("coudln't decode the string", "error", err)
	}
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://admin:admin@localhost/social?sslmode=disable"),
			maxOpenConn: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdleConn: env.GetInt("DB_MAX_IDLE_CONN", 30),
			maxIdleTime: maxIdleTime,
		},
		env: env.GetString("ENV", "development"),
		redisCfg: redisConfig{
			address:  env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PW", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		mail: mailConfig{
			exp:       3 * 24 * time.Hour, // 3 days
			fromEmail: env.GetString("SENDGRID_FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "root"),
			},
			token: tokenConfig{
				secret: secret,
				exp:    24 * time.Hour,
			},
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
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.New(cfg.redisCfg.address, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
		defer func() {
			if err := rdb.Close(); err != nil {
				logger.Errorf("couldn't close the cache connection: %v\n", err)
			}
		}()
	}
	store := store.NewPostgresStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	pasetoAuthenticator, err := auth.NewPasetoAuthenticator(cfg.auth.token.secret)
	if err != nil {
		logger.Panicf("PANIC: couldn't create PASETO authenticator, error: %v", err)
	}

	app := &application{
		config:        cfg,
		storage:       store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: pasetoAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
