package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
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
		<title>ENGi {{.Name}}</title>
		<link rel="icon" type="image/png" href="/favicon.png">
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
	static := flag.String("static", "data", "The relative path to your assets")
	host := flag.String("host", "127.0.0.1", "The host at which to serve your games")
	port := flag.Int("port", 8080, "The port at which to serve your games")
	flag.Parse()

	http.HandleFunc("/", programHandler)
	http.Handle("/favicon.png", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "."}))
	http.HandleFunc(fmt.Sprintf("/%s/", path.Clean(*static)), staticHandler)

	fmt.Printf("Now open your browser to http://%s:%d!\n", *host, *port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
}
