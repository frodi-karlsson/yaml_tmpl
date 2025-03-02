package yaml_tmpl

import (
	"fmt"
	"os"
	"strings"
)

// Takes in a path to a yaml template and returns it transpiled to HTML.
func LoadTemplate(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("LoadTemplate failed to read file: %w", err)
	}

	split := strings.Split(string(content), "\n")

	yamlNodes, err := GetYamlNodesFromLines(split)
	if err != nil {
		return "", fmt.Errorf("LoadTemplate failed to get yaml nodes: %w", err)
	}

	out := ""
	for _, yamlNode := range yamlNodes {
		htmlNode := yamlNode.Transpile(nil)
		out += htmlNode.String()
	}

	return out, nil
}
