package yaml_tmpl_test

import (
	"fmt"
	"testing"

	"github.com/frodi-karlsson/yaml_tmpl"
)

var SIMPLE_RAW_TAG_HTML_NODE = yaml_tmpl.YamlNode{
	Key:     "tag",
	Type:    yaml_tmpl.RAW_YAML_NODE,
	Content: "value",
}

func getSimpleRawTagHtmlNode() yaml_tmpl.YamlNode {
	return yaml_tmpl.YamlNode{
		Key:     "tag",
		Type:    yaml_tmpl.RAW_YAML_NODE,
		Content: "value",
	}
}

func getSimpleChildrenHtmlNode() yaml_tmpl.YamlNode {
	tag := yaml_tmpl.YamlNode{
		Key:  "tag",
		Type: yaml_tmpl.CHILDREN_YAML_NODE,
	}

	children := yaml_tmpl.YamlNode{
		Key:    "children",
		Type:   yaml_tmpl.CHILDREN_YAML_NODE,
		Parent: &tag,
	}

	child := yaml_tmpl.YamlNode{
		Key:     "child",
		Type:    yaml_tmpl.RAW_YAML_NODE,
		Content: "value",
		Parent:  &children,
	}

	children.Children = []yaml_tmpl.YamlNode{child}

	tag.Children = []yaml_tmpl.YamlNode{children}

	return tag
}

func TestParseSimpleRawTagNode(t *testing.T) {
	SIMPLE_RAW_TAG_HTML_NODE := getSimpleRawTagHtmlNode()

	transpiled := SIMPLE_RAW_TAG_HTML_NODE.Transpile(nil)

	res, err := expectHtmlNodeToEqual(t, *transpiled, yaml_tmpl.HtmlNode{
		Type: yaml_tmpl.TAG_HTML_NODE,
		Tag:  "tag",
		Children: []*yaml_tmpl.HtmlNode{
			{
				Type:      yaml_tmpl.RAW_HTML_NODE,
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

	transpiled := SIMPLE_RAW_TAG_HTML_NODE.Transpile(nil)
	html := transpiled.String()
	expected := "<tag>value</tag>"
	if html != expected {
		t.Errorf("Expected %s, got %s", expected, html)
	}
}

func TestParseSimpleHtmlChildrenNode(t *testing.T) {
	SIMPLE_CHILDREN_HTML_NODE := getSimpleChildrenHtmlNode()

	transpiled := SIMPLE_CHILDREN_HTML_NODE.Transpile(nil)
	res, err := expectHtmlNodeToEqual(t, *transpiled, yaml_tmpl.HtmlNode{
		Type: yaml_tmpl.TAG_HTML_NODE,
		Tag:  "tag",
		Children: []*yaml_tmpl.HtmlNode{
			{
				Type:      yaml_tmpl.TAG_HTML_NODE,
				Tag:       "child",
				Attribute: "",
				Children: []*yaml_tmpl.HtmlNode{
					{
						Type:      yaml_tmpl.RAW_HTML_NODE,
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

func expectHtmlNodeToEqual(t *testing.T, node yaml_tmpl.HtmlNode, expected yaml_tmpl.HtmlNode) (bool, string) {
	return _expectHtmlNodeToEqual(t, node, expected, "")
}

func _expectHtmlNodeToEqual(t *testing.T, node yaml_tmpl.HtmlNode, expected yaml_tmpl.HtmlNode, path string) (bool, string) {
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
		res, err := _expectHtmlNodeToEqual(t, *child, *expected.Children[i], fmt.Sprintf("%s.children[%d]", pathLogSuffix, i))
		if !res {
			return res, err
		}
	}

	return true, ""
}

func BenchmarkTranspileSimpleRawTagNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SIMPLE_RAW_TAG_HTML_NODE.Transpile(nil)
	}
}

func BenchmarkTranspileSimpleChildrenNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleChildrenHtmlNode := getSimpleChildrenHtmlNode()
		simpleChildrenHtmlNode.Transpile(nil)
	}
}

func BenchmarkTranspileSimpleChildrenNodeWithGrandchildren(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tag := yaml_tmpl.YamlNode{
			Key:  "tag",
			Type: yaml_tmpl.CHILDREN_YAML_NODE,
		}

		children := yaml_tmpl.YamlNode{
			Key:    "children",
			Type:   yaml_tmpl.CHILDREN_YAML_NODE,
			Parent: &tag,
		}

		grandchildren := yaml_tmpl.YamlNode{
			Key:    "children",
			Type:   yaml_tmpl.CHILDREN_YAML_NODE,
			Parent: &children,
		}

		grandchild := yaml_tmpl.YamlNode{
			Key:     "grandchild",
			Type:    yaml_tmpl.RAW_YAML_NODE,
			Content: "value",
			Parent:  &grandchildren,
		}

		grandchildren.Children = []yaml_tmpl.YamlNode{grandchild}

		children.Children = []yaml_tmpl.YamlNode{grandchildren}

		tag.Children = []yaml_tmpl.YamlNode{children}

		tag.Transpile(nil)
	}
}
