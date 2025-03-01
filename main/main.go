package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/frodi-karlsson/yaml_website"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Cache-Control", "no-cache")
		filename := "index.yaml"
		filePath := filepath.Join("templates", filename)
		rawContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(w, "Failed to read file: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loaded, err := yaml_website.LoadTemplate(filePath)
		if err != nil {
			fmt.Fprintf(w, "Failed to load template: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		templateData := struct {
			YAML string
		}{
			YAML: string(rawContent),
		}

		template, err := template.New("index").Parse(loaded)
		if err != nil {
			fmt.Fprintf(w, "Failed to parse template: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = template.Execute(w, templateData)
		if err != nil {
			fmt.Fprintf(w, "Failed to execute template: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		http.ServeFile(w, r, path)
	})

	fmt.Println("Listening on port", *port)
	http.ListenAndServe(":"+*port, nil)
}
