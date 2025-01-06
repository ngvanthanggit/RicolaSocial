package main

import (
	"log"

	"github.com/ngvanthanggit/RicolaSocial/internal/db"
	"github.com/ngvanthanggit/RicolaSocial/internal/env"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

const version = "0.0.1"

func main() {
	// setting up the configuration
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://thangitcbg:thangitcbg@localhost:5433/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// creating a new database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("database connection pool established")

	// creating a new storage with the created database
	store := store.NewStorage(db)

	// initialisation the application
	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
