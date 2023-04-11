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
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%%%s\n", node.Value()))

		indirectObjectOffsets := []int{}
		trailer := ""

		for _, child := range node.Children() {
			switch child.Type() {
			case ast.INDIRECT_OBJECT:
				indirectObjectOffsets = append(indirectObjectOffsets, sb.Len())

			case ast.TRAILER:
				// Serialise the trailer now as we need to output it
				// after the xref table
				t, err := s.Serialise(child)

				if err != nil {
					return "", err
				}

				trailer = t
				continue
			}

			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised)
		}

		xrefTableStartOffset := sb.Len()

		return fmt.Sprintf(
			"%s%s\n%s\nstartxref\n%d\n%%EOF",
			sb.String(),
			createXrefTable(indirectObjectOffsets),
			trailer,
			xrefTableStartOffset,
		), nil

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
		sb := strings.Builder{}
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised + " ")
		}

		return fmt.Sprintf("<<%s>>", sb.String()), nil

	case ast.STRING:
		return fmt.Sprintf("(%s)", node.Value().(string)), nil

	case ast.FUNCTION:
		sb := strings.Builder{}
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised)
		}

		return fmt.Sprintf("{ %s }", sb.String()), nil

	case ast.ARRAY:
		sb := strings.Builder{}
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised + " ")
		}

		return fmt.Sprintf(
			"[%s]",
			strings.TrimRight(sb.String(), " "),
		), nil

	case ast.STREAM:
		return fmt.Sprintf("\nstream\n%s\nendstream\n", node.Value().(string)), nil

	case ast.XREFS:
		// We don't serialise the xrefs from the AST,
		// we generate them from the indirect objects
		// (done in the root node serialisation)

	case ast.TRAILER:
		sb := strings.Builder{}

		for _, child := range node.Children() {
			if child.Type() != ast.DICT {
				continue
			}

			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised + " ")
		}

		return fmt.Sprintf("trailer\n%s\n", sb.String()), nil

	case ast.INDIRECT_OBJECT:
		id := node.(*ast.IndirectObjectNode).Id()
		generation := node.(*ast.IndirectObjectNode).Gen()
		sb := strings.Builder{}

		sb.WriteString(fmt.Sprintf("%d %d obj\n", id, generation))
		for _, child := range node.Children() {
			serialised, _ := s.Serialise(child)
			sb.WriteString(serialised)
		}

		return sb.String() + "\nendobj\n", nil

	case ast.OBJECT_REF:
		return fmt.Sprintf(
			"%d %d R",
			node.(*ast.ObjectRefNode).Id(),
			node.(*ast.ObjectRefNode).Gen(),
		), nil
	}

	return "", fmt.Errorf("unknown node type: %d", node.Type())
}

func createXrefTable(offsets []int) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("xref\n0 %d\n0000000000 65535 f\n", len(offsets)+1))

	for _, offset := range offsets {
		sb.WriteString(fmt.Sprintf("%010d 00000 n\n", offset))
	}

	return sb.String()
}
