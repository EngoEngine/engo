package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-on/gopherjslib"
	"html/template"
	"net/http"
	"os"
	"path"
)

var page = `
<!DOCTYPE html>
<html>
	<head>
    <title>{{.Name}}</title>
		<style>
		html, body {
			padding: 0;
			margin: 0;
			font-family: Arial;
			background: #000000;
			{{if .Error}}
			background: #f4f4f4;
			{{end}}
			width: 100%;
			height: 100%;
			overflow: hidden;
			font-size: 18px;
		}
		div#error {
			color: #4b5464;
			font-size: 1.75em;
			text-align: center;
			padding: 80px;
		}
		</style>
	</head>
  <body>
		{{if .Error}}
		<div id="error">{{.Error}}</div>
		{{end}}
		<script>
			{{.Script}}
		</script>
  </body>
</html>
`

type Game struct {
	Name   string
	Script template.JS
	Error  error
}

func programHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the html template.
	t, err := template.New("path").Parse(page)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}

	// Determine the program to build.
	path := r.URL.Path[1:]
	if len(path) == 0 {
		path = "main"
	}

	// Open the program's source.
	f, err := os.Open(fmt.Sprintf("%s.go", path))
	if err != nil {
		t.Execute(w, &Game{"Error", "", err})
		return
	}

	// Compile the program's source.
	var out bytes.Buffer
	err = gopherjslib.Build(f, &out, nil)
	if err != nil {
		t.Execute(w, &Game{"Error", "", err})
		return
	}
	script := out.String()

	// Plug in the compiled javascript.
	t.Execute(w, &Game{path, template.JS(script), nil})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	static := flag.String("static", "data", "Path to static files")
	port := flag.Int("port", 8080, "Port to serve on")
	flag.Parse()

	http.HandleFunc("/", programHandler)
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
