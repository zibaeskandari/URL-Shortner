package main

import (
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
	fmt.Printf("Config: %+v\n", config)
}
