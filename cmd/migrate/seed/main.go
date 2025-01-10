package main

import (
	"log"

	"github.com/ngvanthanggit/RicolaSocial/internal/db"
	"github.com/ngvanthanggit/RicolaSocial/internal/env"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://thangitcbg:thangitcbg@localhost:5433/socialnetwork?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.SeedHandler(store)
}
