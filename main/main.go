package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/frodi-karlsson/yaml_website"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Cache-Control", "no-cache")
		yaml_website.ServeYamlFile(w, "index.yaml")
	})

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		http.ServeFile(w, r, path)
	})

	fmt.Println("Listening on port", *port)
	http.ListenAndServe(":"+*port, nil)
}
