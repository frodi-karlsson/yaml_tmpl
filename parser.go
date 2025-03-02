package yaml_tmpl

import (
	"fmt"
	"strings"
)

var _QUOTE_TYPES = [...]string{"'", "\""}

type YamlNodeType int

const (
	// Unknown node type. Only returned if an error occurs.
	UNKNOWN_YAML_NODE YamlNodeType = iota
	// A raw node is a node that contains a single line with a tag and a value.
	RAW_YAML_NODE
	// A children node is a node that contains a tag and a list of children nodes.
	CHILDREN_YAML_NODE
)

type YamlNode struct {
	Key  string
	Type YamlNodeType
	// Only used if Type == CHILDREN_NODE
	Children []YamlNode
	// Only used if Type == RAW_NODE
	//
	// This contains the lines up until an indentation of the same level or lower.
	Content string
	// Nil for a root node.
	Parent *YamlNode
}

// Splits a group of yaml lines into groups of direct children.
func collectGroups(lines []string) ([][]string, error) {
	length := len(lines)
	if length == 0 {
		return [][]string{}, nil
	}

	if length == 1 {
		return [][]string{lines}, nil
	}

	topLevelIndent := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))

	var elements = make([][]string, 0)
	var element = []string{}

	for _, line := range lines {
		if len(element) == 0 {
			element = []string{line}
			continue
		}

		indentation := len(line) - len(strings.TrimLeft(line, " "))
		isTopLevel := indentation == topLevelIndent

		// If the line is at the same indentation as the first line, we have a new element.
		if isTopLevel {
			elements = append(elements, element)
			element = []string{line}
		} else if indentation < topLevelIndent {
			return nil, fmt.Errorf("CollectGroups failed: indentation is less than top level")
		} else {
			element = append(element, line)
		}
	}

	// As we reach the end of the file, we'll always be in the process of building an element.
	// So we need to append it to the elements list.
	elements = append(elements, element)

	return elements, nil
}

// Determines the type of a yaml node.
//
// The first line passed will be the definition for the parent node,
// and all following nodes are children.
func determineNodeType(lines []string) (YamlNodeType, error) {
	if len(lines) == 0 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: no lines")
	}

	// If the definition line contains a quotation as defined in QUOTE_TYPES, it is a raw node.
	for _, quote := range _QUOTE_TYPES {
		if strings.Contains(lines[0], quote) {
			return RAW_YAML_NODE, nil
		}
	}

	// If we don't have more than one line and it's not raw, it must be unknown.
	if len(lines) < 2 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: Could not determine node type of one line and no quotes in %v", lines)
	}

	firstIndentation := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
	secondIndentation := len(lines[1]) - len(strings.TrimLeft(lines[1], " "))

	// If the next line has a higher indentation, it is a children node.
	// If it doesn't, we have an empty node and resolve it as an empty string raw node.
	if secondIndentation <= firstIndentation {
		return RAW_YAML_NODE, nil
	}

	// If there are no lines at the same indentation as the first line, it is a children node.
	nonFirstLineAtFirstIndentationExists := false

	for _, line := range lines[1:] {
		indentation := len(line) - len(strings.TrimLeft(line, " "))

		if indentation == firstIndentation {
			nonFirstLineAtFirstIndentationExists = true
			break
		}
	}

	if !nonFirstLineAtFirstIndentationExists {
		return CHILDREN_YAML_NODE, nil
	}

	return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: There are lines at the same indentation as the first line in %v", lines)
}

// Extracts the content of a raw node.
func extractRawContent(lines []string) (string, error) {
	if len(lines) == 0 {
		return "", fmt.Errorf("ExtractRawContent failed: no lines")
	}

	if len(lines) > 1 {
		return "", fmt.Errorf("ExtractRawContent failed: too many lines")
	}

	definition := lines[0]

	split := strings.SplitN(definition, ":", 2)
	if len(split) < 2 {
		return "", fmt.Errorf("ExtractRawContent failed: no colon in %s", definition)
	}

	rightHandSide := split[1]

	var insideQuote = false
	var quoteType = ""
	var value = ""
	var escaped = false

	for _, char := range rightHandSide {
		// Treat escaped quotes as a single character.
		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		if (char == '\'' || char == '"') && !escaped {
			if !insideQuote {
				insideQuote = true
				quoteType = string(char)
			} else if quoteType == string(char) {
				insideQuote = false
			}
		} else if insideQuote {
			value += string(char)
		}

		escaped = false
	}

	if insideQuote {
		return "", fmt.Errorf("ExtractRawContent failed: missing closing quote")
	}

	return value, nil
}

// Extracts the key of a node.
func extractKey(line string) (string, error) {
	split := strings.SplitN(line, ":", 2)
	if len(split) < 2 {
		return "", fmt.Errorf("ExtractKey failed: no colon in %s", line)
	}
	leftHandSide := split[0]
	trimmedWhitespace := strings.Trim(leftHandSide, " ")
	withoutDash := strings.TrimLeft(trimmedWhitespace, "- ")
	return withoutDash, nil
}

func parseChildrenNode(lines []string, parent *YamlNode) (YamlNode, error) {
	if len(lines) == 0 {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: no lines")
	}

	childLines, err := collectGroups(lines[1:])
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	key, err := extractKey(lines[0])
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	childrenNode := YamlNode{
		Key:    key,
		Type:   CHILDREN_YAML_NODE,
		Parent: parent,
	}

	children := make([]YamlNode, 0)

	for _, childLines := range childLines {
		childNode, err := parseNode(childLines, &childrenNode)
		if err != nil {
			return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
		}

		children = append(children, childNode...)
	}

	childrenNode.Children = children
	return childrenNode, nil
}

func parseRawNode(lines []string, parent *YamlNode) (YamlNode, error) {
	if len(lines) == 0 {
		return YamlNode{}, fmt.Errorf("ParseRawNode failed: no lines")
	}

	content, err := extractRawContent(lines)
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	key, err := extractKey(lines[0])
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	return YamlNode{
		Key:     key,
		Type:    RAW_YAML_NODE,
		Content: content,
		Parent:  parent,
	}, nil
}

func parseNode(lines []string, parent *YamlNode) ([]YamlNode, error) {
	nodeType, err := determineNodeType(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseNode failed: %w", err)
	}

	if nodeType == RAW_YAML_NODE {
		node, err := parseRawNode(lines, parent)
		if err != nil {
			return nil, fmt.Errorf("ParseNode failed: %w", err)
		}

		return []YamlNode{node}, nil
	}

	if nodeType == CHILDREN_YAML_NODE {
		node, err := parseChildrenNode(lines, parent)
		if err != nil {
			return nil, fmt.Errorf("ParseNode failed: %w", err)
		}

		return []YamlNode{node}, nil
	}

	return nil, fmt.Errorf("ParseNode failed: unknown node type")
}

func getNonEmptyLines(lines []string) []string {
	nonEmptyLines := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.Trim(line, " ")
		if len(trimmed) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return nonEmptyLines
}

func GetYamlNodesFromLines(lines []string) ([]YamlNode, error) {
	groups, err := collectGroups(getNonEmptyLines(lines))
	if err != nil {
		return nil, fmt.Errorf("GetYamlNodesFromLines failed: %w", err)
	}

	nodes := make([]YamlNode, 0)

	for _, topLevelLines := range groups {
		node, err := parseNode(topLevelLines, nil)
		if err != nil {
			return nil, fmt.Errorf("GetYamlNodesFromLines failed: %w", err)
		}

		nodes = append(nodes, node...)
	}

	return nodes, nil
}
