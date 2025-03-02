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
	Children []HtmlNode
	// Nil if this is a root node.
	Parent *HtmlNode
}

// Transpiles a raw node to an html node. A raw node is a representation
// of `tag: "content"` in yaml.
func (node *YamlNode) transpileRawNode(parent *HtmlNode) HtmlNode {
	isRootHtmlElement := parent == nil
	isChildElement := node.Parent != nil && node.Parent.Key == "children"
	isAnyHtmlElement := isRootHtmlElement || isChildElement

	if !isAnyHtmlElement {
		// content: is special syntax to specify innerText while setting attributes.
		if node.Key == "innerText" {
			rawNode := HtmlNode{
				Type:    RAW_HTML_NODE,
				Content: node.Content,
				Parent:  parent,
			}
			return rawNode
		}

		return HtmlNode{
			Type:      ATTRIBUTE_HTML_NODE,
			Attribute: node.Key,
			Content:   node.Content,
			Parent:    parent,
		}
	}

	rawNode := HtmlNode{
		Type:    RAW_HTML_NODE,
		Content: node.Content,
	}

	// raw: is special syntax to denote a Raw Text element in html.
	if node.Key == "raw" {
		rawNode.Parent = parent
		return rawNode
	}

	rawHtmlNode := HtmlNode{
		Type:     TAG_HTML_NODE,
		Tag:      node.Key,
		Children: []HtmlNode{rawNode},
		Parent:   parent,
	}
	rawNode.Parent = &rawHtmlNode

	if len(node.Children) == 0 || node.Key == "innerText" {
		return rawHtmlNode
	}

	return HtmlNode{
		Type:      ATTRIBUTE_HTML_NODE,
		Attribute: node.Key,
		Content:   node.Content,
	}
}

// Transpiles a children node to an html node. A children node is a representation
// of `tag: anything: ...` in yaml.
func (node *YamlNode) transpileChildrenNode(parent *HtmlNode) HtmlNode {
	htmlNode := HtmlNode{
		Type:     TAG_HTML_NODE,
		Tag:      node.Key,
		Children: []HtmlNode{},
		Parent:   parent,
	}

	children := []HtmlNode{}
	for _, child := range node.Children {
		// children: is special syntax to denote child elements.
		if child.Type == CHILDREN_YAML_NODE && child.Key == "children" {
			for _, grandchild := range child.Children {
				children = append(children, grandchild.Transpile(&htmlNode))
			}
		} else {
			children = append(children, child.Transpile(&htmlNode))
		}
	}

	htmlNode.Children = children

	return htmlNode
}

// Determines the type of a node based on its content.
func (node *YamlNode) Transpile(parent *HtmlNode) HtmlNode {
	switch node.Type {
	case RAW_YAML_NODE:
		return node.transpileRawNode(parent)

	case CHILDREN_YAML_NODE:
		return node.transpileChildrenNode(parent)
	default:
		return HtmlNode{
			Type:   UNKNOWN_HTML_NODE,
			Parent: parent,
		}
	}
}

func splitHtmlNodesByIsType(nodes []HtmlNode, nodeType HtmlNodeType) (matching []HtmlNode, rest []HtmlNode) {
	matching = []HtmlNode{}
	rest = []HtmlNode{}

	for _, node := range nodes {
		if node.Type == nodeType {
			matching = append(matching, node)
		} else {
			rest = append(rest, node)
		}
	}

	return matching, rest
}

// Converts an HTML node to a string.
func (node *HtmlNode) ToString() string {
	switch node.Type {
	case RAW_HTML_NODE:
		return node.Content
	case TAG_HTML_NODE:
		attributes := ""
		children := ""

		attrChildren, otherChildren := splitHtmlNodesByIsType(node.Children, ATTRIBUTE_HTML_NODE)

		for _, node := range attrChildren {
			if node.Type == ATTRIBUTE_HTML_NODE {
				attributes += " " + node.Attribute + "=\"" + node.Content + "\""
			}
		}

		for _, child := range otherChildren {
			children += child.ToString()
		}

		return "<" + node.Tag + attributes + ">" + children + "</" + node.Tag + ">"
	case ATTRIBUTE_HTML_NODE:
		return node.Attribute + "=\"" + node.Content + "\""
	default:
		return ""
	}
}
