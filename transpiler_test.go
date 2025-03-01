package yaml_website_test

import (
	"fmt"
	"testing"

	"github.com/frodi-karlsson/yaml_website"
)

var SIMPLE_RAW_TAG_HTML_NODE = yaml_website.YamlNode{
	Key:     "tag",
	Type:    yaml_website.RAW_YAML_NODE,
	Content: "value",
}

func getSimpleRawTagHtmlNode() yaml_website.YamlNode {
	return yaml_website.YamlNode{
		Key:     "tag",
		Type:    yaml_website.RAW_YAML_NODE,
		Content: "value",
	}
}

func getSimpleChildrenHtmlNode() yaml_website.YamlNode {
	tag := yaml_website.YamlNode{
		Key:  "tag",
		Type: yaml_website.CHILDREN_YAML_NODE,
	}

	children := yaml_website.YamlNode{
		Key:    "children",
		Type:   yaml_website.CHILDREN_YAML_NODE,
		Parent: &tag,
	}

	child := yaml_website.YamlNode{
		Key:     "child",
		Type:    yaml_website.RAW_YAML_NODE,
		Content: "value",
		Parent:  &children,
	}

	children.Children = []yaml_website.YamlNode{child}

	tag.Children = []yaml_website.YamlNode{children}

	return tag
}

func TestParseSimpleRawTagNode(t *testing.T) {
	SIMPLE_RAW_TAG_HTML_NODE := getSimpleRawTagHtmlNode()

	transpiled := yaml_website.TranspileNode(SIMPLE_RAW_TAG_HTML_NODE, nil)

	res, err := expectHtmlNodeToEqual(t, transpiled, yaml_website.HtmlNode{
		Type: yaml_website.TAG_HTML_NODE,
		Tag:  "tag",
		Children: []yaml_website.HtmlNode{
			{
				Type:      yaml_website.RAW_HTML_NODE,
				Tag:       "",
				Attribute: "",
				Content:   "value",
			},
		},
	})

	if !res {
		t.Errorf("Got unexpected result: %s for:\n%v", err, transpiled)
	}
}

func TestPrintSimpleHtmlRawTagNode(t *testing.T) {
	SIMPLE_RAW_TAG_HTML_NODE := getSimpleRawTagHtmlNode()

	transpiled := yaml_website.TranspileNode(SIMPLE_RAW_TAG_HTML_NODE, nil)
	html := yaml_website.HtmlNodeToString(transpiled)
	expected := "<tag>value</tag>"
	if html != expected {
		t.Errorf("Expected %s, got %s", expected, html)
	}
}

func TestParseSimpleHtmlChildrenNode(t *testing.T) {
	SIMPLE_CHILDREN_HTML_NODE := getSimpleChildrenHtmlNode()

	transpiled := yaml_website.TranspileNode(SIMPLE_CHILDREN_HTML_NODE, nil)
	res, err := expectHtmlNodeToEqual(t, transpiled, yaml_website.HtmlNode{
		Type: yaml_website.TAG_HTML_NODE,
		Tag:  "tag",
		Children: []yaml_website.HtmlNode{
			{
				Type:      yaml_website.TAG_HTML_NODE,
				Tag:       "child",
				Attribute: "",
				Children: []yaml_website.HtmlNode{
					{
						Type:      yaml_website.RAW_HTML_NODE,
						Tag:       "",
						Attribute: "",
						Content:   "value",
						Children:  nil,
					},
				},
			},
		},
	})

	if !res {
		t.Errorf("Got unexpected result: %s for:\n%v", err, transpiled)
	}
}

func expectHtmlNodeToEqual(t *testing.T, node yaml_website.HtmlNode, expected yaml_website.HtmlNode) (bool, string) {
	return _expectHtmlNodeToEqual(t, node, expected, "")
}

func _expectHtmlNodeToEqual(t *testing.T, node yaml_website.HtmlNode, expected yaml_website.HtmlNode, path string) (bool, string) {
	pathLogSuffix := "root"
	if path != "" {
		pathLogSuffix = path
	}

	if node.Type != expected.Type {
		return false, fmt.Sprintf("Expected type to be %d, got %d at %s", expected.Type, node.Type, pathLogSuffix)
	}

	if node.Tag != expected.Tag {
		return false, fmt.Sprintf("Expected tag to be %s, got %s at %s", expected.Tag, node.Tag, pathLogSuffix)
	}

	if node.Attribute != expected.Attribute {
		return false, fmt.Sprintf("Expected attribute to be %s, got %s at %s", expected.Attribute, node.Attribute, pathLogSuffix)
	}

	if node.Content != expected.Content {
		return false, fmt.Sprintf("Expected content to be %s, got %s at %s", expected.Content, node.Content, pathLogSuffix)
	}

	if len(node.Children) != len(expected.Children) {
		return false, fmt.Sprintf("Expected %d children, got %d at %s", len(expected.Children), len(node.Children), pathLogSuffix)
	}

	for i, child := range node.Children {
		res, err := _expectHtmlNodeToEqual(t, child, expected.Children[i], fmt.Sprintf("%s.children[%d]", pathLogSuffix, i))
		if !res {
			return res, err
		}
	}

	return true, ""
}
