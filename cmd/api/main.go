package main

import (
	"time"

	"github.com/ngvanthanggit/RicolaSocial/internal/db"
	"github.com/ngvanthanggit/RicolaSocial/internal/env"
	"github.com/ngvanthanggit/RicolaSocial/internal/mailer"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			RicolaSocial API
//	@description	Hi, this is Thang. Welcome!
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// setting up the configuration
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURl: env.GetString("FRONTEND_URL", "localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://thangitcbg:thangitcbg@localhost:5433/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			fromEmail: env.GetString("SENDGRID_FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			exp: time.Hour * 2, // 2 hours
		},
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// creating a new database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	// creating a new storage with the created database
	store := store.NewStorage(db)

	// timing setup
	if err = store.DBTime.DBTimeSetup(); err != nil {
		logger.Fatal(err)
	}

	mailer := mailer.NewSendgrid(cfg.mail.fromEmail, cfg.mail.sendGrid.apiKey)

	// initialisation the application
	app := application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
