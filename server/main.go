package main

import (
	"log"
	"messages/src/router"
)

const PORT = ":8080"

func main() {

	router := router.NewRouter()
	log.Println("Starting at port ", PORT)
	router.Run(PORT)

}
