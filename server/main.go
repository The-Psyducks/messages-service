//coverage:ignore
package main

import (
	"log"
	"messages/src/router"
)

const PORT = ":8080"

func main() {

	r, err := router.NewRouter(router.DEFAULT)
	if err != nil {
		log.Fatalln("Error creating router: ", err)
	}

	log.Println("Starting at port ", PORT)
	if err = r.Run(PORT); err != nil {
		log.Fatalln("Error starting server: ", err)
	}

}
