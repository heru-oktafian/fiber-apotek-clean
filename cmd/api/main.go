package main

import (
	"fmt"
	"log"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/bootstrap"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/console"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(console.StartupBanner(app.Config.AppName, "v0.2", "0.0.0.0", app.Config.ServerPort, len(app.Fiber.GetRoutes()), false))
	log.Fatal(app.Fiber.Listen(":" + app.Config.ServerPort))
}
