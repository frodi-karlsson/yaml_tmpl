package yaml_tmpl

type HtmlNodeType int

const (
	// Unknown node type. Only returned if an error occurs.
	UNKNOWN_HTML_NODE HtmlNodeType = iota
	// A raw node is raw innerText.
	RAW_HTML_NODE
	// A tag node is a node that contains a tag and a list of children nodes.
	TAG_HTML_NODE
	// An attribute node represents an attribute of a tag.
	ATTRIBUTE_HTML_NODE
)

type HtmlNode struct {
	Type HtmlNodeType
	// Only used if Type == TAG_NODE
	Tag string
	// Only used if Type == ATTRIBUTE_NODE
	Attribute string
	// Only used if Type == RAW_NODE
	Content string
	// Only used if Type == TAG_NODE
	Children []*HtmlNode
	// Nil if this is a root node.
	Parent *HtmlNode
}

// Transpiles a raw node to an html node. A raw node is a representation
// of `tag: "content"` in yaml.
func (node *YamlNode) transpileRawNode(parent *HtmlNode) *HtmlNode {
	isRootHtmlElement := parent == nil
	isChildElement := node.Parent != nil && node.Parent.Key == "children"
	isAnyHtmlElement := isRootHtmlElement || isChildElement

	if !isAnyHtmlElement {
		// Handle attributes and innerText separately
		if node.Key == "innerText" {
			return &HtmlNode{
				Type:    RAW_HTML_NODE,
				Content: node.Content,
				Parent:  parent,
			}
		}
		return &HtmlNode{
			Type:      ATTRIBUTE_HTML_NODE,
			Attribute: node.Key,
			Content:   node.Content,
			Parent:    parent,
		}
	}

	rawNode := &HtmlNode{
		Type:    RAW_HTML_NODE,
		Content: node.Content,
		Parent:  parent,
	}

	if node.Key == "raw" {
		return rawNode
	}

	if len(node.Children) == 0 || node.Key == "innerText" {
		return &HtmlNode{
			Type:     TAG_HTML_NODE,
			Tag:      node.Key,
			Children: []*HtmlNode{rawNode},
			Parent:   parent,
		}
	}

	// Handle unexpected case
	return &HtmlNode{
		Type:      ATTRIBUTE_HTML_NODE,
		Attribute: node.Key,
		Content:   node.Content,
		Parent:    parent,
	}
}

// Transpiles a children node to an html node. A children node is a representation
// of `tag: anything: ...` in yaml.
func (node *YamlNode) transpileChildrenNode(parent *HtmlNode) *HtmlNode {
	htmlNode := HtmlNode{
		Type:     TAG_HTML_NODE,
		Tag:      node.Key,
		Children: make([]*HtmlNode, 0, len(node.Children)),
		Parent:   parent,
	}

	for _, child := range node.Children {
		// children: is special syntax to denote child elements.
		if child.Type == CHILDREN_YAML_NODE && child.Key == "children" {
			for _, grandchild := range child.Children {
				htmlNode.Children = append(htmlNode.Children, grandchild.Transpile(&htmlNode))
			}
		} else {
			htmlNode.Children = append(htmlNode.Children, child.Transpile(&htmlNode))
		}
	}

	return &htmlNode
}

// Determines the type of a node based on its content.
func (node *YamlNode) Transpile(parent *HtmlNode) *HtmlNode {
	switch node.Type {
	case RAW_YAML_NODE:
		return node.transpileRawNode(parent)

	case CHILDREN_YAML_NODE:
		return node.transpileChildrenNode(parent)
	default:
		return &HtmlNode{
			Type:   UNKNOWN_HTML_NODE,
			Parent: parent,
		}
	}
}

// Converts an HTML node to a string.
func (node *HtmlNode) String() string {
	switch node.Type {
	case RAW_HTML_NODE:
		return node.Content
	case TAG_HTML_NODE:
		attributes := ""
		children := ""

		for _, node := range node.Children {
			if node.Type == ATTRIBUTE_HTML_NODE {
				attributes += " " + node.Attribute + "=\"" + node.Content + "\""
			} else {
				children += node.String()
			}
		}

		return "<" + node.Tag + attributes + ">" + children + "</" + node.Tag + ">"
	case ATTRIBUTE_HTML_NODE:
		return node.Attribute + "=\"" + node.Content + "\""
	default:
		return ""
	}
}
