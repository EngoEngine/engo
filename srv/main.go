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
			background: #222222;
			width: 100%;
			height: 100%;
			overflow: hidden;
		}
		</style>
	</head>
  <body>
		<script>
			{{.Script}}
		</script>
  </body>
</html>
`

type Game struct {
	Name   string
	Script template.JS
}

func programHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	f, err := os.Open(fmt.Sprintf("%s.go", path))
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}
	var out bytes.Buffer
	err = gopherjslib.Build(f, &out, nil)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}
	script := out.String()

	t, err := template.New("path").Parse(page)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}
	t.Execute(w, &Game{path, template.JS(script)})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	static := flag.String("static", "data", "Path to static files")
	flag.Parse()

	http.HandleFunc("/", programHandler)
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)
	http.ListenAndServe(":8080", nil)
}
