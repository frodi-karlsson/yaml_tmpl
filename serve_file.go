package yaml_website

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func ServeYamlFile(w http.ResponseWriter, filename string) {
	filePath := filepath.Join("templates", filename)
	yamlNodes, err := GetYamlNodesFromFile(filePath)
	if err != nil {
		fmt.Fprintf(w, "Failed to get yaml nodes from file: %v", err)
		return
	}

	html := []HtmlNode{}
	for _, node := range yamlNodes {
		html = append(html, TranspileNode(node, nil))
	}

	str := ""
	for _, line := range html {
		str += HtmlNodeToString(line) + "\n"
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(str))
}
