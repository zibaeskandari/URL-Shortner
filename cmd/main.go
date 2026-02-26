package main

import (
	"URLShortner/internal/adapters/repository"
	"URLShortner/internal/core/services"
	"URLShortner/internal/infrastructure/persistence/db"
	"URLShortner/pkg"
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("URL Shortener Project Started!")

	configPath := flag.String("config", "./configs/dev.yaml", "Path to config file")
	flag.Parse()

	config, err := pkg.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	postgresDb, err := db.NewPostgresDb(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDb.Close()

	client, err := postgresDb.Open()
	if err != nil {
		log.Fatal(err)
	}

	urlRepository := repository.NewPgUrlRepository(client)
	shortCodeGeneratorService := services.NewCodeGeneratorService(urlRepository, config.ShortCode)
	fmt.Printf("Sample Short Code %s\n", shortCodeGeneratorService.GetShortCode())
}
