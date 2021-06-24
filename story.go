package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

const defaultHandlerTemplate = `
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
			<li><a href="/{{.Arc}}">{{.Text}}</a></li>
		{{end}}
		</ul>
	<section>
	</body>
</html>
`

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

type HandlerOptions func(h *handler)

func WithTemplate(t *template.Template) HandlerOptions {
	return func(h *handler) {
		h.t = t
	}
}

func WithPath(fn func(r *http.Request) string) HandlerOptions {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func NewHandler(s Story, opts ...HandlerOptions) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
