package yaml_website

import (
	"fmt"
	"os"
	"strings"
)

var QUOTE_TYPES = [...]string{"'", "\""}

type YamlNodeType int

const (
	// Unknown node type. Only returned if an error occurs.
	UNKNOWN_YAML_NODE YamlNodeType = iota
	// A raw node is a node that contains a single line with a tag and a value.
	RAW_YAML_NODE
	// A children node is a node that contains a tag and a list of children nodes.
	CHILDREN_YAML_NODE
)

type Element struct {
	Tag      string
	Class    string
	Children []Element
}

type YamlNode struct {
	Key  string
	Type YamlNodeType
	// Only used if Type == CHILDREN_NODE
	Children []YamlNode
	// Only used if Type == RAW_NODE
	//
	// This contains the lines up until an indentation of the same level or lower.
	Content string
}

// Reads a template and returns the content as a string.
func readTemplate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("ReadTemplate failed to read file: %w", err)
	}

	return string(content), nil
}

// Splits a template into lines.
func splitTemplate(template string) []string {
	return strings.Split(template, "\n")
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
	var element []string = nil

	for _, line := range lines {
		if element == nil {
			element = []string{line}
			continue
		}

		indentation := len(line) - len(strings.TrimLeft(line, " "))
		isTopLevel := indentation == topLevelIndent

		if isTopLevel {
			elements = append(elements, element)
			element = []string{line}
		} else {
			element = append(element, line)
		}
	}

	if element != nil {
		elements = append(elements, element)
	}

	return elements, nil
}

// Determines the type of a yaml node.
func determineNodeType(lines []string) (YamlNodeType, error) {
	if len(lines) == 0 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: no lines")
	}
	// The first node is the definition of the node itself.
	// E.g: "tag: value" with an optional indentation and "-"

	// If the definition line contains a quotation as defined in QUOTE_TYPES, it is a raw node.
	for _, quote := range QUOTE_TYPES {
		if strings.Contains(lines[0], quote) {
			return RAW_YAML_NODE, nil
		}
	}

	// If the next line has a higher indentation, it is a children node.
	if len(lines) > 1 {
		firstIndentation := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
		secondIndentation := len(lines[1]) - len(strings.TrimLeft(lines[1], " "))

		if secondIndentation > firstIndentation {
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
	}

	return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: Could not determine node type in %v", lines)
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

	for _, char := range rightHandSide {
		if char == '\'' || char == '"' {
			if !insideQuote {
				insideQuote = true
				quoteType = string(char)
			} else if quoteType == string(char) {
				insideQuote = false
			}
		} else if insideQuote {
			value += string(char)
		}
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

func parseChildrenNode(lines []string) (YamlNode, error) {
	if len(lines) == 0 {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: no lines")
	}

	childLines, err := collectGroups(lines[1:])
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	children := make([]YamlNode, 0)

	for _, childLines := range childLines {
		childNode, err := parseNode(childLines)
		if err != nil {
			return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
		}

		children = append(children, childNode...)
	}

	key, err := extractKey(lines[0])
	if err != nil {
		return YamlNode{}, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	return YamlNode{
		Key:      key,
		Type:     CHILDREN_YAML_NODE,
		Children: children,
	}, nil
}

func parseRawNode(lines []string) (YamlNode, error) {
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
	}, nil
}

func parseNode(lines []string) ([]YamlNode, error) {
	nodeType, err := determineNodeType(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseNode failed: %w", err)
	}

	if nodeType == RAW_YAML_NODE {
		node, err := parseRawNode(lines)
		if err != nil {
			return nil, fmt.Errorf("ParseNode failed: %w", err)
		}

		return []YamlNode{node}, nil
	}

	if nodeType == CHILDREN_YAML_NODE {
		node, err := parseChildrenNode(lines)
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
		node, err := parseNode(topLevelLines)
		if err != nil {
			return nil, fmt.Errorf("GetYamlNodesFromLines failed: %w", err)
		}

		nodes = append(nodes, node...)
	}

	return nodes, nil
}

// Gets a yaml file as an array of yaml nodes.
func GetYamlNodesFromFile(filename string) ([]YamlNode, error) {
	template, err := readTemplate(filename)
	if err != nil {
		return nil, fmt.Errorf("GetYamlNodes failed to read template: %w", err)
	}

	lines := splitTemplate(template)
	return GetYamlNodesFromLines(lines)
}
