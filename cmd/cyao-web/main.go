package main

import (
	cyao "chooseAdventure"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	file := flag.String("file", "gopher.json", "the JSON file with the cyao story")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f)
	var story cyao.Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
