package main

import (
	"log"
)

func main() {
	s := newServer("api/v1", 3000)
	s.AddGlobalMiddleware(s.Logger)

	log.Fatalln(s.ListenAndServe())
}
