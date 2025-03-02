package yaml_tmpl

import (
	"fmt"
	"strings"
)

var _QUOTE_TYPES = [...]rune{'\'', '"'}

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

func getIndentation(line string) int {
	length := 0
	for _, char := range line {
		if char == ' ' {
			length++
		} else if char == '\t' {
			length += 4
		} else {
			break
		}
	}
	return length
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

	topLevelIndent := getIndentation(lines[0])

	var elements = make([][]string, 0, length)
	var element = make([]string, 0, length)
	elementLength := 0

	for _, line := range lines {
		if elementLength == 0 {
			element = make([]string, 0, length)
			element = append(element, line)
			elementLength++
			continue
		}

		indentation := getIndentation(line)

		if indentation < topLevelIndent {
			return nil, fmt.Errorf("CollectGroups failed: indentation level is lower than top level")
		}

		isTopLevel := indentation == topLevelIndent

		// If the line is at the same indentation as the first line, we have a new element.
		if isTopLevel {
			elements = append(elements, element)
			element = []string{line}
			elementLength = 1
		} else {
			element = append(element, line)
			elementLength++
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
	lineLength := len(lines)

	if lineLength == 0 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: no lines")
	}

	// If the definition line contains a quotation as defined in QUOTE_TYPES, it is a raw node.
	for _, char := range lines[0] {
		for _, quoteType := range _QUOTE_TYPES {
			if char == quoteType {
				return RAW_YAML_NODE, nil
			}
		}
	}

	// If we don't have more than one line and it's not raw, it must be unknown.
	if lineLength < 2 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: Could not determine node type of one line and no quotes in %v", lines)
	}

	firstIndentation := getIndentation(lines[0])
	secondIndentation := getIndentation(lines[1])

	// If the next line has a higher indentation, it is a children node.
	// If it doesn't, we have an empty node and resolve it as an empty string raw node.
	if secondIndentation <= firstIndentation {
		return RAW_YAML_NODE, nil
	}

	// If there are no lines at the same indentation as the first line, it is a children node.
	nonFirstLineAtFirstIndentationExists := false

	for _, line := range lines[1:] {
		indentation := getIndentation(line)

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
	lineLength := len(lines)

	if lineLength == 0 {
		return "", fmt.Errorf("ExtractRawContent failed: no lines")
	}

	if lineLength > 1 {
		return "", fmt.Errorf("ExtractRawContent failed: too many lines")
	}

	definition := lines[0]

	colonIndex := strings.Index(definition, ":")
	if colonIndex == -1 {
		return "", fmt.Errorf("ExtractRawContent failed: no colon in %s", definition)
	}

	rightHandSide := definition[colonIndex+1:]

	var insideQuote = false
	var quoteType rune
	var value = make([]rune, 0, len(rightHandSide))
	var escaped = false

	for _, char := range rightHandSide {
		if char == '#' && !insideQuote {
			break
		}

		// Treat escaped quotes as a single character.
		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		if (char == '\'' || char == '"') && !escaped {
			if !insideQuote {
				insideQuote = true
				quoteType = char
			} else if quoteType == char {
				insideQuote = false
			}
		} else if insideQuote {
			value = append(value, char)
		}

		escaped = false
	}

	if insideQuote {
		return "", fmt.Errorf("ExtractRawContent failed: missing closing quote")
	}

	return string(value), nil
}

// Extracts the key of a node.
func extractKey(line string) (string, error) {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", fmt.Errorf("ExtractKey failed: no colon in %s", line)
	}

	leftHandSide := line[:colonIndex]
	trimmed := strings.TrimLeft(leftHandSide, "- ")
	return trimmed, nil
}

func parseChildrenNode(lines []string, parent *YamlNode) (*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseChildrenNode failed: no lines")
	}

	childLines, err := collectGroups(lines[1:])
	if err != nil {
		return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	key, err := extractKey(lines[0])
	if err != nil {
		return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	var childrenNode YamlNode
	childrenNode.Key = key
	childrenNode.Type = CHILDREN_YAML_NODE
	childrenNode.Parent = parent

	children := make([]YamlNode, 0, len(childLines))

	for _, childLines := range childLines {
		childNode, err := parseNode(childLines, &childrenNode)
		if err != nil {
			return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
		}

		children = append(children, *childNode)
	}

	childrenNode.Children = children
	return &childrenNode, nil
}

func parseRawNode(lines []string, parent *YamlNode) (*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseRawNode failed: no lines")
	}

	content, err := extractRawContent(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	key, err := extractKey(lines[0])
	if err != nil {
		return nil, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	return &YamlNode{
		Key:     key,
		Type:    RAW_YAML_NODE,
		Content: content,
		Parent:  parent,
	}, nil
}

func parseNode(lines []string, parent *YamlNode) (*YamlNode, error) {
	nodeType, err := determineNodeType(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseNode failed: %w", err)
	}
	switch nodeType {
	case RAW_YAML_NODE:
		return parseRawNode(lines, parent)
	case CHILDREN_YAML_NODE:
		return parseChildrenNode(lines, parent)
	default:
		return nil, fmt.Errorf("ParseNode failed: unknown node type")
	}
}

func getNonEmptyLines(lines []string) []string {
	nonEmptyLines := make([]string, 0, len(lines))

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

	nodes := make([]YamlNode, 0, len(groups))

	for _, topLevelLines := range groups {
		node, err := parseNode(topLevelLines, nil)
		if err != nil {
			return nil, fmt.Errorf("GetYamlNodesFromLines failed: %w", err)
		}

		nodes = append(nodes, *node)
	}

	return nodes, nil
}
