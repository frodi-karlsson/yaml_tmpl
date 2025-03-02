package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/frodi-karlsson/yaml_tmpl"
)

func main() {
	static := flag.Bool("static", false, "Build static page")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	if !*static {
		startServer(port)
	} else {
		buildStatic()
	}
}

func startServer(port *string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Cache-Control", "no-cache")

		path := filepath.Join("templates", "index.yaml")

		stylesCss, err := os.ReadFile("static/style.css")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read style.css: %v", err), http.StatusInternalServerError)
			return
		}

		template, templateData, err := loadTemplate(path, "index", string(stylesCss))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load template: %v", err), http.StatusInternalServerError)
			return
		}

		err = template.Execute(w, templateData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
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

func buildStatic() error {
	outDir := filepath.Join("..", "docs")

	files, err := filepath.Glob("templates/*.yaml")
	if err != nil {
		return fmt.Errorf("BuildStatic failed to get files: %w", err)
	}

	// Delete and recreate dist directory
	exists, err := os.Stat(outDir)
	if err == nil {
		if exists.IsDir() {
			err = os.RemoveAll(outDir)
			if err != nil {
				return fmt.Errorf("BuildStatic failed to remove directory: %w", err)
			}
		}
	}

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return fmt.Errorf("BuildStatic failed to create directory: %w", err)
	}

	// Create out/static dir
	err = os.MkdirAll(filepath.Join(outDir, "static"), 0755)
	if err != nil {
		return fmt.Errorf("BuildStatic failed to create directory: %w", err)
	}

	// Copy /static/ files to dist
	staticFiles, err := filepath.Glob("static/*")
	if err != nil {
		return fmt.Errorf("BuildStatic failed to get static files: %w", err)
	}

	for _, staticFile := range staticFiles {
		content, err := os.ReadFile(staticFile)
		if err != nil {
			return fmt.Errorf("BuildStatic failed to read file: %w", err)
		}

		err = os.WriteFile(filepath.Join(outDir, staticFile), content, 0644)
		if err != nil {
			return fmt.Errorf("BuildStatic failed to write file: %w", err)
		}
	}

	styleCss, err := os.ReadFile("static/style.css")
	if err != nil {
		return fmt.Errorf("BuildStatic failed to read style.css: %w", err)
	}

	for _, file := range files {
		// Write yaml
		template, templateData, err := loadTemplate(file, filepath.Base(file), string(styleCss))
		if err != nil {
			return fmt.Errorf("BuildStatic failed to load template: %w", err)
		}

		name := filepath.Base(file)
		outPath := filepath.Join(outDir, name[:len(name)-5]+".html")
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("BuildStatic failed to create file: %w", err)
		}

		err = template.Execute(out, templateData)
		if err != nil {
			return fmt.Errorf("BuildStatic failed to execute template: %w", err)
		}

		out.Close()
	}

	return nil
}

// Load a template from a file and return a parsed template and the data to be used with it.
// You can also pass in a CSS string to be used in the template.
func loadTemplate(path string, name string, css string) (*template.Template, interface{}, error) {
	rawContent, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadTemplate failed to read file: %w", err)
	}

	loaded, err := yaml_tmpl.LoadTemplate(path)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadTemplate failed to load template: %w", err)
	}

	templateData := struct {
		YAML string
		CSS  string
	}{
		YAML: string(rawContent),
		CSS:  css,
	}

	template, err := template.New(name).Parse(loaded)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadTemplate failed to parse template: %w", err)
	}

	return template, templateData, nil
}
