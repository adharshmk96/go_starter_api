/*
Copyright © 2025 Adharsh Manikandan <debugslayer@gmail.com>
*/
package main

import (
	"log"
	"servicehub_api/cmd"
	_ "servicehub_api/docs"

	"github.com/joho/godotenv"
)

// @title			ServiceHub API
// @version		1.0
// @description	This is the API for the ServiceHub platform.
// @host			localhost:8080
// @BasePath		/
// @schemes		http
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not loaded... ignoring")
	}

	cmd.Execute()
}
