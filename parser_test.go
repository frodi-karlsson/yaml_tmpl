package yaml_tmpl_test

import (
	"fmt"
	"testing"

	"github.com/frodi-karlsson/yaml_tmpl"
)

var SIMPLE_DOUBLE_QUOTE_RAW_NODE = []string{
	"tag: \"value\"",
}

var SIMPLE_SINGLE_QUOTE_RAW_NODE = []string{
	"tag: 'value'",
}

var ESCAPED_DOUBLE_QUOTE_RAW_NODE = []string{
	"tag: \"value \\\"with escaped quotes\\\"\"",
}

var SIMPLE_CHILDREN_NODE = []string{
	"tag:",
	"  child: \"value\"",
}

var NESTED_CHILDREN_NODE = []string{
	"tag:",
	"  child:",
	"    - nested1: \"value\" # Anything after a non-quoted # is ignored by the parser",
	"    - nested2: \"value\"",
}

var DOCUMENT_NODE = []string{
	"head:",
	"  children:",
	"    - title: \"Stupid YAML Website\"",
	"    - link:",
	"      rel: \"stylesheet\"",
	"      type: \"text/css\"",
	"      href: \"/static/style.css\"",
	"body:",
	"  children:",
	"    - h1:",
	"        class: \"title\"",
	"        text: \"Welcome to the Stupid YAML Website\"",
	"    - p: \"The template is written in YAML like God intended\"",
}

var RAW_NODE_WITH_ANCHOR_NAME = []string{
	"tag: &anchor \"value\"",
}

var SIMPLE_CHILDREN_NODE_WITH_ANCHOR_NAME = []string{
	"tag: &anchor",
	"  child: \"value\"",
}

func TestParseSimpleDoubleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(SIMPLE_DOUBLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:     "tag",
		Type:    yaml_tmpl.RAW_YAML_NODE,
		Content: "value",
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseSimpleSingleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(SIMPLE_SINGLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]
	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:     "tag",
		Type:    yaml_tmpl.RAW_YAML_NODE,
		Content: "value",
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseEscapedDoubleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(ESCAPED_DOUBLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]
	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:     "tag",
		Type:    yaml_tmpl.RAW_YAML_NODE,
		Content: "value \"with escaped quotes\"",
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseSimpleChildrenNode(t *testing.T) {
	// Test a simple children node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(SIMPLE_CHILDREN_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:  "tag",
		Type: yaml_tmpl.CHILDREN_YAML_NODE,
		Children: []yaml_tmpl.YamlNode{
			{
				Key:     "child",
				Type:    yaml_tmpl.RAW_YAML_NODE,
				Content: "value",
			},
		},
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseNestedChildrenNode(t *testing.T) {
	// Test a nested children node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(NESTED_CHILDREN_NODE)

	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]
	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:  "tag",
		Type: yaml_tmpl.CHILDREN_YAML_NODE,
		Children: []yaml_tmpl.YamlNode{
			{
				Key:  "child",
				Type: yaml_tmpl.CHILDREN_YAML_NODE,
				Children: []yaml_tmpl.YamlNode{
					{
						Key:     "nested1",
						Type:    yaml_tmpl.RAW_YAML_NODE,
						Content: "value",
					},
					{
						Key:     "nested2",
						Type:    yaml_tmpl.RAW_YAML_NODE,
						Content: "value",
					},
				},
			},
		},
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseDocumentNode(t *testing.T) {
	// Test a full document node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(DOCUMENT_NODE)

	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	head := nodes[0]

	res, msg := expectYamlNodeToEqual(t, head, yaml_tmpl.YamlNode{
		Key:  "head",
		Type: yaml_tmpl.CHILDREN_YAML_NODE,
		Children: []yaml_tmpl.YamlNode{
			{
				Key:  "children",
				Type: yaml_tmpl.CHILDREN_YAML_NODE,
				Children: []yaml_tmpl.YamlNode{
					{
						Key:     "title",
						Type:    yaml_tmpl.RAW_YAML_NODE,
						Content: "Stupid YAML Website",
					},
					{
						Key:  "link",
						Type: yaml_tmpl.CHILDREN_YAML_NODE,
						Children: []yaml_tmpl.YamlNode{
							{
								Key:     "rel",
								Type:    yaml_tmpl.RAW_YAML_NODE,
								Content: "stylesheet",
							},
							{
								Key:     "type",
								Type:    yaml_tmpl.RAW_YAML_NODE,
								Content: "text/css",
							},
							{
								Key:     "href",
								Type:    yaml_tmpl.RAW_YAML_NODE,
								Content: "/static/style.css",
							},
						},
					},
				},
			},
		},
	})

	if !res {
		t.Error(msg)
	}

	body := nodes[1]

	res, msg = expectYamlNodeToEqual(t, body, yaml_tmpl.YamlNode{
		Key:  "body",
		Type: yaml_tmpl.CHILDREN_YAML_NODE,
		Children: []yaml_tmpl.YamlNode{
			{
				Key:  "children",
				Type: yaml_tmpl.CHILDREN_YAML_NODE,
				Children: []yaml_tmpl.YamlNode{
					{
						Key:  "h1",
						Type: yaml_tmpl.CHILDREN_YAML_NODE,
						Children: []yaml_tmpl.YamlNode{
							{
								Key:     "class",
								Type:    yaml_tmpl.RAW_YAML_NODE,
								Content: "title",
							},
							{
								Key:     "text",
								Type:    yaml_tmpl.RAW_YAML_NODE,
								Content: "Welcome to the Stupid YAML Website",
							},
						},
					},
					{
						Key:     "p",
						Type:    yaml_tmpl.RAW_YAML_NODE,
						Content: "The template is written in YAML like God intended",
					},
				},
			},
		},
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseRawNodeWithAnchorName(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(RAW_NODE_WITH_ANCHOR_NAME)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:        "tag",
		Type:       yaml_tmpl.RAW_YAML_NODE,
		Content:    "value",
		AnchorName: "anchor",
	})

	if !res {
		t.Error(msg)
	}
}

func TestParseSimpleChildrenNodeWithAnchorName(t *testing.T) {
	// Test a simple children node
	nodes, err := yaml_tmpl.GetYamlNodesFromLines(SIMPLE_CHILDREN_NODE_WITH_ANCHOR_NAME)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	res, msg := expectYamlNodeToEqual(t, node, yaml_tmpl.YamlNode{
		Key:        "tag",
		Type:       yaml_tmpl.CHILDREN_YAML_NODE,
		AnchorName: "anchor",
		Children: []yaml_tmpl.YamlNode{
			{
				Key:     "child",
				Type:    yaml_tmpl.RAW_YAML_NODE,
				Content: "value",
			},
		},
	})

	if !res {
		t.Error(msg)
	}
}

func expectYamlNodeToEqual(t *testing.T, node yaml_tmpl.YamlNode, expected yaml_tmpl.YamlNode) (bool, string) {
	return _expectYamlNodeToEqual(t, node, expected, "")
}

func _expectYamlNodeToEqual(t *testing.T, node yaml_tmpl.YamlNode, expected yaml_tmpl.YamlNode, path string) (bool, string) {
	pathLogSuffix := "root"
	if path != "" {
		pathLogSuffix = path
	}

	if node.Key != expected.Key {
		return false, fmt.Sprintf("Expected key to be %s, got %s at %s", expected.Key, node.Key, pathLogSuffix)
	}

	if node.Type != expected.Type {
		return false, fmt.Sprintf("Expected type to be %d, got %d at %s", expected.Type, node.Type, pathLogSuffix)
	}

	if node.Content != expected.Content {
		return false, fmt.Sprintf("Expected content to be %s, got %s at %s", expected.Content, node.Content, pathLogSuffix)
	}

	if len(node.Children) != len(expected.Children) {
		return false, fmt.Sprintf("Expected %d children, got %d at %s", len(expected.Children), len(node.Children), pathLogSuffix)
	}

	for i, child := range node.Children {
		if i >= len(expected.Children) {
			return false, fmt.Sprintf("Unexpected child at %d at %s", i, pathLogSuffix)
		}

		res, err := expectYamlNodeToEqual(t, child, expected.Children[i])
		if !res {
			return false, fmt.Sprintf("Unexpected result: %s for:\n%v", err, child)
		}
	}

	return true, ""
}

func BenchmarkParseSimpleDoubleQuoteNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		yaml_tmpl.GetYamlNodesFromLines(SIMPLE_DOUBLE_QUOTE_RAW_NODE)
	}
}

func BenchmarkParseSimpleSingleQuoteNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		yaml_tmpl.GetYamlNodesFromLines(SIMPLE_SINGLE_QUOTE_RAW_NODE)
	}
}

func BenchmarkParseEscapedDoubleQuoteNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		yaml_tmpl.GetYamlNodesFromLines(ESCAPED_DOUBLE_QUOTE_RAW_NODE)
	}
}

func BenchmarkParseSimpleChildrenNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		yaml_tmpl.GetYamlNodesFromLines(SIMPLE_CHILDREN_NODE)
	}
}

func BenchmarkParseNestedChildrenNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		yaml_tmpl.GetYamlNodesFromLines(NESTED_CHILDREN_NODE)
	}
}
