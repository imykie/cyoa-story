package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"html/template"
)

var tpl *template.Template

func init(){
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

const defaultHandlerTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title> Choose Your Own Adventure </title>
	</head>
	<body>
	<h1>{{.Title}}</h1>
	{{range .Story}}
		<p>{{.}}</p>
	{{end}}
	<ul>
	{{range .Options}}
		<li><a href="/{{.Arc}}">{{.Text}}</a></li>
	{{end}}
	</ul>
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
	s Story
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
