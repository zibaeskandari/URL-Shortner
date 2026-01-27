package main

import (
	"URLShortner/pkg"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("URL Shortener Project Started!")
	config, err := pkg.LoadConfig("./configs/" + getConfigProfile() + ".yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Config: %+v\n", config)
}

func getConfigProfile() string {
	env := os.Getenv("APP_SERVER_ENV")
	if env == "" {
		env = "dev"
	}
	return env
}
