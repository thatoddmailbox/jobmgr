package server

import (
	"log"
	"net/http"
	"strconv"
)

func StartServer() {
	initRoutes()

	port := 9845
	log.Printf("API server listening on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}
