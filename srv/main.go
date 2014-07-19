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
	// Determine the program to build.
	path := r.URL.Path[1:]
	if len(path) == 0 {
		path = "main"
	}

	// Open the program's source.
	f, err := os.Open(fmt.Sprintf("%s.go", path))
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}

	// Compile the program's source.
	var out bytes.Buffer
	err = gopherjslib.Build(f, &out, nil)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}
	script := out.String()

	// Plug the compiled js into the template.
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
	port := flag.Int("port", 8080, "Port to serve on")
	flag.Parse()

	http.HandleFunc("/", programHandler)
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
