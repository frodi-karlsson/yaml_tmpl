package yaml_website

import (
	"strings"
)

type HtmlNodeType int

const (
	// Unknown node type. Only returned if an error occurs.
	UNKNOWN_NODE HtmlNodeType = iota
	// A raw node is raw innerText.
	RAW_NODE
	// A tag node is a node that contains a tag and a list of children nodes.
	TAG_NODE
	// An attribute node represents an attribute of a tag.
	ATTRIBUTE_NODE
)

var HTML_TAGS = [...]string{
	"a",
	"abbr",
	"address",
	"area",
	"article",
	"aside",
	"audio",
	"b",
	"base",
	"bdi",
	"bdo",
	"blockquote",
	"body",
	"br",
	"button",
	"canvas",
	"caption",
	"cite",
	"code",
	"col",
	"colgroup",
	"data",
	"datalist",
	"dd",
	"del",
	"details",
	"dfn",
	"dialog",
	"div",
	"dl",
	"dt",
	"em",
	"embed",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"head",
	"header",
	"hgroup",
	"hr",
	"html",
	"i",
	"iframe",
	"img",
	"input",
	"ins",
	"kbd",
	"label",
	"legend",
	"li",
	"link",
	"main",
	"map",
	"mark",
	"meta",
	"meter",
	"nav",
	"noscript",
	"object",
	"ol",
	"optgroup",
	"option",
	"output",
	"p",
	"param",
	"picture",
	"pre",
	"progress",
	"q",
	"rp",
	"rt",
	"ruby",
	"s",
	"samp",
	"script",
	"section",
	"select",
	"slot",
	"small",
	"source",
	"span",
	"strong",
	"style",
	"sub",
	"summary",
	"sup",
	"table",
	"tbody",
	"td",
	"template",
	"textarea",
	"tfoot",
	"th",
	"thead",
	"time",
	"title",
	"tr",
	"track",
	"u",
	"ul",
	"var",
	"video",
	"wbr",
}

var HTML_ATTRIBUTES = [...]string{
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"align",
	"alt",
	"async",
	"autocomplete",
	"autofocus",
	"autoplay",
	"charset",
	"checked",
	"cite",
	"class",
	"color",
	"cols",
	"colspan",
	"content",
	"contenteditable",
	"controls",
	"coords",
	"data",
	"datetime",
	"default",
	"defer",
	"dir",
	"dirname",
	"disabled",
	"download",
	"draggable",
	"dropzone",
	"enctype",
	"for",
	"form",
	"formaction",
	"headers",
	"height",
	"hidden",
	"high",
	"href",
	"hreflang",
	"http-equiv",
	"icon",
	"id",
	"ismap",
	"itemprop",
	"keytype",
	"kind",
	"label",
	"lang",
	"language",
	"list",
	"loop",
	"low",
	"manifest",
	"max",
	"maxlength",
	"media",
	"method",
	"min",
	"multiple",
	"muted",
	"name",
	"novalidate",
	"open",
	"optimum",
	"pattern",
	"ping",
	"placeholder",
	"poster",
	"preload",
	"radiogroup",
	"readonly",
	"rel",
	"required",
	"reversed",
	"rows",
	"rowspan",
	"sandbox",
	"scope",
	"scoped",
	"seamless",
	"selected",
	"shape",
	"size",
	"sizes",
	"span",
	"spellcheck",
	"src",
	"srcdoc",
	"srclang",
	"srcset",
	"start",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"type",
	"usemap",
	"value",
	"width",
	"wrap",
}

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
}

func isHtmlTag(tag string) bool {
	for _, htmlTag := range HTML_TAGS {
		if strings.EqualFold(htmlTag, tag) {
			return true
		}
	}
	return false
}

func isHtmlAttribute(attribute string) bool {
	for _, htmlAttribute := range HTML_ATTRIBUTES {
		if strings.EqualFold(htmlAttribute, attribute) {
			return true
		}
	}
	return false
}

// Deternimes the type of a node based on its content.
func TranspileNode(node YamlNode) HtmlNode {
	switch node.Type {
	case RAW_YAML_NODE:
		if isHtmlTag(node.Key) {
			rawNode := HtmlNode{
				Type:    RAW_NODE,
				Content: node.Content,
			}

			return HtmlNode{
				Type:     TAG_NODE,
				Tag:      node.Key,
				Children: []HtmlNode{rawNode},
			}
		} else if isHtmlAttribute(node.Key) {
			return HtmlNode{
				Type:      ATTRIBUTE_NODE,
				Attribute: node.Key,
				Content:   node.Content,
			}
		} else {
			return HtmlNode{
				Type:    RAW_NODE,
				Content: node.Content,
			}
		}
	case CHILDREN_YAML_NODE:
		children := []HtmlNode{}
		for _, child := range node.Children {
			children = append(children, TranspileNode(child))
		}

		return HtmlNode{
			Type:     TAG_NODE,
			Tag:      node.Key,
			Children: children,
		}
	default:
		return HtmlNode{
			Type: UNKNOWN_NODE,
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
func HtmlNodeToString(node HtmlNode) string {
	switch node.Type {
	case RAW_NODE:
		return node.Content
	case TAG_NODE:
		attributes := ""
		children := ""

		attrChildren, otherChildren := splitHtmlNodesByIsType(node.Children, ATTRIBUTE_NODE)

		for _, node := range attrChildren {
			if node.Type == ATTRIBUTE_NODE {
				attributes += " " + node.Attribute + "=\"" + node.Content + "\""
			}
		}

		for _, child := range otherChildren {
			children += HtmlNodeToString(child)
		}

		if len(attributes) > 0 {
			attributes = " " + attributes
		}

		return "<" + node.Tag + attributes + ">" + children + "</" + node.Tag + ">"
	case ATTRIBUTE_NODE:
		return node.Attribute + "=\"" + node.Content + "\""
	default:
		return ""
	}
}
