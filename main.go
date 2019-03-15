package main

import (
	"fmt"
)

func main() {
	listenPort := ":5000"

	fmt.Println("Initializing server...")
	a := App{}

	// Uncomment this if using environment variables.
	// a.Initialize(
	//     os.Getenv("APP_DB_USERNAME"),
	//     os.Getenv("APP_DB_PASSWORD"),
	//     os.Getenv("APP_DB_NAME"))

	// Uncomment this if using static values.
	a.Initialize("nvuser", "nvuser", "gonv")

	fmt.Println("Starting up server at localhost", listenPort)
	a.Run(listenPort)
}
