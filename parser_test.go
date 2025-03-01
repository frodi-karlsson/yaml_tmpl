package yaml_website_test

import (
	"testing"

	"github.com/frodi-karlsson/yaml_website"
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
	"    - nested1: \"value\"",
	"    - nested2: \"value\"",
}

var DOCUMENT_NODE = []string{
	"head:",
	"  children:",
	"    - title: \"Stupid YAML Website\"",
	"    - link:",
	"        rel: \"stylesheet\"",
	"        type: \"text/css\"",
	"        href: \"/static/style.css\"",
	"body:",
	"  children:",
	"    - h1:",
	"        class: \"title\"",
	"        text: \"Welcome to the Stupid YAML Website\"",
	"    - p: \"The template is written in YAML like God intended\"",
}

func TestParseSimpleDoubleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_website.GetYamlNodesFromLines(SIMPLE_DOUBLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	if node.Key != "tag" {
		t.Errorf("Expected key to be 'tag', got %s", node.Key)
	}

	if node.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected type to be RAW_NODE, got %d", node.Type)
	}

	if node.Content != "value" {
		t.Errorf("Expected content to be 'value', got %s", node.Content)
	}
}

func TestParseSimpleSingleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_website.GetYamlNodesFromLines(SIMPLE_SINGLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	if node.Key != "tag" {
		t.Errorf("Expected key to be 'tag', got %s", node.Key)
	}

	if node.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected type to be RAW_NODE, got %d", node.Type)
	}

	if node.Content != "value" {
		t.Errorf("Expected content to be 'value', got %s", node.Content)
	}
}

func TestParseEscapedDoubleQuoteNode(t *testing.T) {
	// Test a simple raw node
	nodes, err := yaml_website.GetYamlNodesFromLines(ESCAPED_DOUBLE_QUOTE_RAW_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	if node.Key != "tag" {
		t.Errorf("Expected key to be 'tag', got %s", node.Key)
	}

	if node.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected type to be RAW_NODE, got %d", node.Type)
	}

	if node.Content != "value \"with escaped quotes\"" {
		t.Errorf("Expected content to be 'value \"with escaped quotes\"', got %s", node.Content)
	}
}

func TestParseSimpleChildrenNode(t *testing.T) {
	// Test a simple children node
	nodes, err := yaml_website.GetYamlNodesFromLines(SIMPLE_CHILDREN_NODE)
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	if node.Key != "tag" {
		t.Errorf("Expected key to be 'tag', got %s", node.Key)
	}

	if node.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected type to be CHILDREN_NODE, got %d", node.Type)
	}

	if len(node.Children) != 1 {
		t.Errorf("Expected children to have 1 element, got %d", len(node.Children))
	}

	child := node.Children[0]

	if child.Key != "child" {
		t.Errorf("Expected child key to be 'child', got %s", child.Key)
	}

	if child.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected child type to be RAW_NODE, got %d", child.Type)
	}

	if child.Content != "value" {
		t.Errorf("Expected child content to be 'value', got %s", child.Content)
	}
}

func TestParseNestedChildrenNode(t *testing.T) {
	// Test a nested children node
	nodes, err := yaml_website.GetYamlNodesFromLines(NESTED_CHILDREN_NODE)

	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]

	if node.Key != "tag" {
		t.Errorf("Expected key to be 'tag', got %s", node.Key)
	}

	if node.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected type to be CHILDREN_NODE, got %d", node.Type)
	}

	if len(node.Children) != 1 {
		t.Errorf("Expected children to have 1 element, got %d", len(node.Children))
	}

	child := node.Children[0]

	if child.Key != "child" {
		t.Errorf("Expected child key to be 'child', got %s", child.Key)
	}

	if child.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected child type to be CHILDREN_NODE, got %d", child.Type)
	}

	if len(child.Children) != 2 {
		t.Errorf("Expected child to have 2 children, got %d", len(child.Children))
	}

	nested1 := child.Children[0]

	if nested1.Key != "nested1" {
		t.Errorf("Expected nested1 key to be 'nested1', got %s", nested1.Key)
	}

	if nested1.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected nested1 type to be RAW_NODE, got %d", nested1.Type)
	}

	if nested1.Content != "value" {
		t.Errorf("Expected nested1 content to be 'value', got %s", nested1.Content)
	}

	nested2 := child.Children[1]

	if nested2.Key != "nested2" {
		t.Errorf("Expected nested2 key to be 'nested2', got %s", nested2.Key)
	}

	if nested2.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected nested2 type to be RAW_NODE, got %d", nested2.Type)
	}

	if nested2.Content != "value" {
		t.Errorf("Expected nested2 content to be 'value', got %s", nested2.Content)
	}
}

func TestParseDocumentNode(t *testing.T) {
	// Test a full document node
	nodes, err := yaml_website.GetYamlNodesFromLines(DOCUMENT_NODE)

	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	head := nodes[0]

	if head.Key != "head" {
		t.Errorf("Expected head key to be 'head', got %s", head.Key)
	}

	if head.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected head type to be CHILDREN_NODE, got %d", head.Type)
	}

	if len(head.Children) != 1 {
		t.Errorf("Expected head to have 1 child, got %d", len(head.Children))
	}

	headChild := head.Children[0]

	if headChild.Key != "children" {
		t.Errorf("Expected head child key to be 'children', got %s", headChild.Key)
	}

	if headChild.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected head child type to be CHILDREN_NODE, got %d", headChild.Type)
	}

	if len(headChild.Children) != 2 {
		t.Errorf("Expected head child to have 2 children, got %d", len(headChild.Children))
	}

	title := headChild.Children[0]

	if title.Key != "title" {
		t.Errorf("Expected title key to be 'title', got %s", title.Key)
	}

	if title.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected title type to be RAW_NODE, got %d", title.Type)
	}

	if title.Content != "Stupid YAML Website" {
		t.Errorf("Expected title content to be 'Stupid YAML Website', got %s", title.Content)
	}

	link := headChild.Children[1]

	if link.Key != "link" {
		t.Errorf("Expected link key to be 'link', got %s", link.Key)
	}

	if link.Type != yaml_website.CHILDREN_YAML_NODE {
		t.Errorf("Expected link type to be CHILDREN_NODE, got %d", link.Type)
	}

	if len(link.Children) != 3 {
		t.Errorf("Expected link to have 3 children, got %d", len(link.Children))
	}

	rel := link.Children[0]

	if rel.Key != "rel" {
		t.Errorf("Expected rel key to be 'rel', got %s", rel.Key)
	}

	if rel.Type != yaml_website.RAW_YAML_NODE {
		t.Errorf("Expected rel type to be RAW_NODE, got %d", rel.Type)
	}

	if rel.Content != "stylesheet" {
		t.Errorf("Expected rel content to be 'stylesheet', got %s", rel.Content)
	}

	linkType := link.Children[1]

	if linkType.Key != "type" {
		t.Errorf("Expected type key to be 'type', got %s", linkType.Key)
	}
}
