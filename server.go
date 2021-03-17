package main

import (
	"log"
	"net/http"
	"os"

	"github.com/paul/go-server/web_client"
)

const defaultPort = "8081"

func main() {
	web_client.InitMaps()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/go-client", web_client.GraphQLQueryHandler())
	http.Handle("/go-client/mutation", web_client.GraphQLMutationHandler())

	log.Printf("connect to http://localhost:%s/go-client to play", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
