package main

import (
	"log"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Fiber.Listen(":" + app.Config.ServerPort))
}
