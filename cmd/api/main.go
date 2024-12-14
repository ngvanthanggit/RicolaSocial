package main

import (
	"log"
)

func main() {
	config := config{
		Addr: ":8080",
	}
	app := application{
		config: config,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
