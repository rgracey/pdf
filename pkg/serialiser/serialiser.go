package serialiser

import (
	"fmt"
	"strings"

	"github.com/rgracey/pdf/pkg/ast"
)

type Serialiser interface {
	Serialise(node ast.PdfNode) (string, error)
}

type AstSerialiser struct {
}

func NewSerialiser() Serialiser {
	return &AstSerialiser{}
}

func (s *AstSerialiser) Serialise(node ast.PdfNode) (string, error) {
	switch node.Type() {
	case ast.ROOT:
		data := "%PDF-1.4\n" // TODO - store in root node
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			data += string(serialised)
		}
		return data + "\n%%EOF", nil

	case ast.BOOLEAN:
		switch node.Value().(bool) {
		case true:
			return "true", nil
		case false:
			return "false", nil
		}

	case ast.FLOAT:
		return fmt.Sprintf("%f", node.Value().(float64)), nil

	case ast.INTEGER:
		return fmt.Sprintf("%d", node.Value().(int64)), nil

	case ast.NAME:
		return fmt.Sprintf("/%s", node.Value().(string)), nil

	case ast.DICT:
		dict := ""
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			dict += string(serialised) + " "
		}

		return fmt.Sprintf("<<%s>>", dict), nil

	case ast.STRING:
		return fmt.Sprintf("(%s)", node.Value().(string)), nil

	case ast.FUNCTION:
		function := ""
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			function += string(serialised)
		}

		return fmt.Sprintf("{ %s }", function), nil

	case ast.ARRAY:
		array := ""
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			array += string(serialised) + " "
		}

		return fmt.Sprintf(
			"[%s]",
			strings.TrimRight(array, " "),
		), nil

	case ast.STREAM:
		return fmt.Sprintf("\nstream\n%s\nendstream\n", node.Value().(string)), nil

	case ast.XREFS:
		xrefTable := ""
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			xrefTable += string(serialised) + " "
		}

		return fmt.Sprintf("xref\n0 %d\n%s\n", len(node.Children()), xrefTable), nil

	case ast.INDIRECT_OBJECT:
		id := node.(*ast.IndirectObjectNode).Id()
		generation := node.(*ast.IndirectObjectNode).Gen()

		data := fmt.Sprintf("%d %d obj\n", id, generation)
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			data += string(serialised)
		}

		return data + "\nendobj\n", nil

	case ast.OBJECT_REF:
		return fmt.Sprintf(
			"%d %d R",
			node.(*ast.ObjectRefNode).Id(),
			node.(*ast.ObjectRefNode).Gen(),
		), nil
	}

	return "", fmt.Errorf("unknown node type: %d", node.Type())
}
