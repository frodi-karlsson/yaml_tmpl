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

	// Below types are internal and will never be returned.

	// An alias node is a node that contains a tag and a reference to another node.
	_ALIAS_YAML_NODE
	// An override node is an alias node with the tag "<<"
	_OVERRIDE_YAML_NODE
)

type YamlNode struct {
	Key  string
	Type YamlNodeType
	// Only used if Type == CHILDREN_NODE
	Children []*YamlNode
	// Only used if Type == RAW_NODE
	//
	// This contains the lines up until an indentation of the same level or lower.
	Content string
	// Nil for a root node.
	Parent *YamlNode
	// Empty string if this node is not an anchor
	AnchorName string
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

	definition := lines[0]

	// If the definition line contains a quotation as defined in QUOTE_TYPES, it is a raw node.
	for _, char := range definition {
		for _, quoteType := range _QUOTE_TYPES {
			if char == quoteType {
				return RAW_YAML_NODE, nil
			}
		}
	}

	// If the definition line contains an asterisk, it can be an alias or an override.
	if strings.Contains(definition, "*") {
		key, err := parseKey(definition)
		if err != nil {
			return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: %w", err)
		}

		if key == "<<" {
			return _OVERRIDE_YAML_NODE, nil
		}
		return _ALIAS_YAML_NODE, nil
	}

	// If we don't have more than one line and it's not raw, it must be unknown.
	if lineLength < 2 {
		return UNKNOWN_YAML_NODE, fmt.Errorf("DetermineNodeType failed: Could not determine node type of one line and no quotes in %v", lines)
	}

	firstIndentation := getIndentation(definition)
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

	colonIndex := strings.IndexRune(definition, ':')
	if colonIndex == -1 {
		return "", fmt.Errorf("ExtractRawContent failed: no colon in %s", definition)
	}

	rightHandSide := definition[colonIndex+1:]

	var insideQuote = false
	var quoteType rune
	var value = make([]rune, 0, len(rightHandSide))
	var escaped = false

	for _, char := range rightHandSide {
		if char == '#' && !insideQuote { // The rest of the line is a comment.
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

// Parses the key of a node.
func parseKey(line string) (string, error) {
	colonIndex := strings.IndexRune(line, ':')
	if colonIndex == -1 {
		return "", fmt.Errorf("ExtractKey failed: no colon in %s", line)
	}

	leftHandSide := line[:colonIndex]
	trimmed := strings.TrimLeft(leftHandSide, "- ")
	return trimmed, nil
}

func isSpecial(byte byte) bool {
	return byte == ' ' || byte == '\t' || byte == '\n' || byte == '\r' || byte == '#' || byte == ':' || byte == '&' || byte == '*'
}

// Extracts the anchor name from a line and returns (modifiedLine, anchorName).
func extractAnchorName(line string) (string, string) {
	anchorIndex := strings.IndexRune(line, '&')
	if anchorIndex == -1 {
		return line, ""
	}

	lineLength := len(line)

	if anchorIndex == lineLength-1 {
		return line[:anchorIndex], ""
	}

	anchorName := make([]rune, 0, lineLength-anchorIndex-1)

	var endIndex int
	for index, char := range line[anchorIndex+1:] {
		endIndex = index + anchorIndex + 1
		if isSpecial(byte(char)) {
			break
		}
		anchorName = append(anchorName, char)
	}

	return line[:anchorIndex] + line[endIndex:], string(anchorName)
}

func parseChildrenNode(lines []string, parent *YamlNode, anchorMap map[string]*YamlNode) (*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseChildrenNode failed: no lines")
	}

	definition := lines[0]

	childLines, err := collectGroups(lines[1:])
	if err != nil {
		return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	definition, anchorName := extractAnchorName(definition)

	key, err := parseKey(definition)
	if err != nil {
		return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
	}

	var childrenNode YamlNode
	childrenNode.Key = key
	childrenNode.Type = CHILDREN_YAML_NODE
	childrenNode.Parent = parent
	childrenNode.AnchorName = anchorName

	children := make([]*YamlNode, 0, len(childLines))

	for _, childLines := range childLines {
		childNodes, err := parseNode(childLines, &childrenNode, anchorMap)
		if err != nil {
			return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
		}

		children = append(children, childNodes...)
	}

	childrenNode.Children = children

	if anchorName != "" {
		anchorMap[anchorName] = &childrenNode
	}

	return &childrenNode, nil
}

func parseRawNode(lines []string, parent *YamlNode, anchorMap map[string]*YamlNode) (*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseRawNode failed: no lines")
	}

	content, err := extractRawContent(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	definition, anchorName := extractAnchorName(lines[0])

	key, err := parseKey(definition)
	if err != nil {
		return nil, fmt.Errorf("ParseRawNode failed: %w", err)
	}

	nodePtr := &YamlNode{
		Key:        key,
		Type:       RAW_YAML_NODE,
		Content:    content,
		Parent:     parent,
		AnchorName: anchorName,
	}

	if anchorName != "" {
		anchorMap[anchorName] = nodePtr
	}

	return nodePtr, nil
}

func getAnchor(definition string, anchorMap map[string]*YamlNode) (*YamlNode, error) {
	asteriskIndex := strings.IndexRune(definition, '*')
	if asteriskIndex == -1 {
		return nil, fmt.Errorf("GetAnchor failed: no asterisk")
	}

	anchorName := make([]rune, 0, len(definition)-asteriskIndex-1)
	for _, char := range definition[asteriskIndex+1:] {
		if isSpecial(byte(char)) {
			break
		}
		anchorName = append(anchorName, char)
	}

	anchor, exists := anchorMap[string(anchorName)]
	if !exists {
		allValues := ""
		for key, value := range anchorMap {
			allValues += key + ", " + value.Key + "\n"
		}
		return nil, fmt.Errorf("GetAnchor failed: anchor not found, all values:\n%s", allValues)
	}

	return anchor, nil
}

// Parses an alias node.
func parseAliasNode(lines []string, parent *YamlNode, anchorMap map[string]*YamlNode) (*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseAliasNode failed: no lines")
	}

	definition := lines[0]
	key, err := parseKey(definition)
	if err != nil {
		return nil, fmt.Errorf("ParseAliasNode failed: %w", err)
	}

	anchor, err := getAnchor(definition, anchorMap)
	if err != nil {
		return nil, fmt.Errorf("ParseAliasNode failed: %w", err)
	}

	childrenNode := YamlNode{
		Key:    key,
		Type:   CHILDREN_YAML_NODE,
		Parent: parent,
	}

	copy := *anchor
	copy.Parent = &childrenNode
	childrenNode.Children = []*YamlNode{&copy}

	return &childrenNode, nil
}

// Parses an override node.
func parseOverrideNode(lines []string, parent *YamlNode, anchorMap map[string]*YamlNode) ([]*YamlNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("ParseOverrideNode failed: no lines")
	}

	definition := lines[0]

	anchor, err := getAnchor(definition, anchorMap)
	if err != nil {
		return nil, fmt.Errorf("ParseOverrideNode failed: %w", err)
	}

	if anchor.Type != CHILDREN_YAML_NODE {
		return nil, fmt.Errorf("ParseOverrideNode failed: anchor is not a children node")
	}

	childNodes := make([]*YamlNode, 0, len(anchor.Children))
	for _, child := range anchor.Children {
		copy := *child
		copy.Parent = parent
		childNodes = append(childNodes, &copy)
	}

	return childNodes, nil
}

// Parses a node. It returns an array, because override nodes may return multiple nodes.
func parseNode(lines []string, parent *YamlNode, anchorMap map[string]*YamlNode) ([]*YamlNode, error) {
	nodeType, err := determineNodeType(lines)
	if err != nil {
		return nil, fmt.Errorf("ParseNode failed: %w", err)
	}
	switch nodeType {
	case RAW_YAML_NODE:
		node, err := parseRawNode(lines, parent, anchorMap)
		if err != nil {
			return nil, fmt.Errorf("ParseRawNode failed: %w", err)
		}
		if node.AnchorName != "" {
			return []*YamlNode{}, nil
		}
		return []*YamlNode{node}, nil
	case CHILDREN_YAML_NODE:
		node, err := parseChildrenNode(lines, parent, anchorMap)
		if err != nil {
			return nil, fmt.Errorf("ParseChildrenNode failed: %w", err)
		}
		if node.AnchorName != "" {
			return []*YamlNode{}, nil
		}
		return []*YamlNode{node}, nil
	case _ALIAS_YAML_NODE:
		node, err := parseAliasNode(lines, parent, anchorMap)
		if err != nil {
			return nil, fmt.Errorf("ParseAliasNode failed: %w", err)
		}
		return []*YamlNode{node}, nil
	case _OVERRIDE_YAML_NODE:
		return parseOverrideNode(lines, parent, anchorMap)
	default:
		return nil, fmt.Errorf("ParseNode failed: unknown node type")
	}
}

func getNonEmptyLines(lines []string) []string {
	nonEmptyLines := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.Trim(line, " ")
		if len(trimmed) > 0 && trimmed[0] != '#' {
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

	anchorMap := make(map[string]*YamlNode)

	for _, topLevelLines := range groups {
		childNodes, err := parseNode(topLevelLines, nil, anchorMap)
		if err != nil {
			return nil, fmt.Errorf("GetYamlNodesFromLines failed: %w", err)
		}

		for _, childNode := range childNodes {
			nodes = append(nodes, *childNode)
		}
	}

	return nodes, nil
}
