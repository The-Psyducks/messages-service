package main

import (
	"log"
	"messages/src/router"
)

const PORT = ":8080"

func main() {

	router, err := router.NewRouter()
	if err != nil {
		log.Fatalln("Error creating router: ", err)
	}

	log.Println("Starting at port ", PORT)
	router.Run(PORT)

}
