package main

import (
	cyoa "chooseAdventure"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const storyTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title> Choose Your Own Adventure </title>
	</head>
	<style>
		body {
			font-family: helvetica, arial, tahoma, sans-serif;
		}
		h1 {
			text-align: center;
			position: relative;
		}
		.page {
			width: 80%;
			max-width: 500px;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
			padding: 80px;
			background: #fffcf6;
			border: 1px solid #eee;
			box-shadow: 0 10px 6px -6px #777;
		}
		ul {
			border-top: 1px dotted #ccc;
			padding: 10px 0 0 0;
			-webkit-padding-start: 0;
		}
		li {
			padding-top: 10px;
		}
		a, a:visited {
			text-decoration: none;
			color: #6295b5;
		}
		a:active, a:hover {
			color: #7792a2;
		}
		p {
			text-indent: 1em;
			color: #616161;
		}
	</style>
	<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Story}}
			<p>{{.}}</p>
		{{end}}
		<ul>
		{{range .Options}}
			<li><a href="/story/{{.Arc}}">{{.Text}}</a></li>
		{{end}}
		</ul>
	<section>
	</body>
</html>
`

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

	tpl := template.Must(template.New("").Parse(storyTemplate))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPath(pathFn))
	//h := cyoa.NewHandler(story, nil)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("starting server on port: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
