package main

import (
	"log"

	"github.com/ngvanthanggit/RicolaSocial/internal/env"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

func main() {
	config := config{
		Addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := application{
		config: config,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
