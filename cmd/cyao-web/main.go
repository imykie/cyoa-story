package main

import (
	cyoa "chooseAdventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	file := flag.String("file", "gopher.json", "the JSON file with the cyao story")
	port := flag.Int("port", 3000, "port to start the cyoa server on")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story, nil)
	fmt.Printf("starting server on port: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
