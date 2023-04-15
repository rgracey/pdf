package parser_test

import (
	"testing"

	"github.com/rgracey/pdf/pkg/ast"
	"github.com/rgracey/pdf/pkg/parser"
	"github.com/rgracey/pdf/pkg/token"
)

func TestParser_ParsesVersion(t *testing.T) {
	tokeniser := &mockTokeniser{
		tokens: []token.Token{
			{Type: token.COMMENT, Value: "PDF-1.7"},
		},
	}

	parser := parser.NewParser(tokeniser)
	root := parser.Parse()

	expectNode(t, root, ast.ROOT, "PDF-1.7")
	expectChildren(t, root, []ast.PdfNode{})
}

func TestParser_ParsesIndirectObject(t *testing.T) {
	tokeniser := &mockTokeniser{
		tokens: []token.Token{
			{Type: token.NUMBER_INTEGER, Value: int64(1)},
			{Type: token.NUMBER_INTEGER, Value: int64(0)},
			{Type: token.KEYWORD, Value: "obj"},
			{Type: token.DICT_START},
			{Type: token.NAME, Value: "Type"},
			{Type: token.NAME, Value: "Catalog"},
			{Type: token.DICT_END},
			{Type: token.KEYWORD, Value: "endobj"},
		},
	}

	parser := parser.NewParser(tokeniser)
	root := parser.Parse()

	expectNode(t, root, ast.ROOT, nil)
	expectChildren(t, root, []ast.PdfNode{
		ast.NewIndirectObjectNode(1, 0),
	})

	obj := root.Children()[0]

	expectNode(t, obj, ast.INDIRECT_OBJECT, nil)
	expectChildren(t, obj, []ast.PdfNode{
		ast.NewDictNode(),
	})

	dict := obj.Children()[0]

	expectNode(t, dict, ast.DICT, nil)
	expectChildren(t, dict, []ast.PdfNode{
		ast.NewNameNode("Type"),
		ast.NewNameNode("Catalog"),
	})
}

// expectChildren checks that a node has the expected number and type of children
func expectChildren(t *testing.T, node ast.PdfNode, expectedChildren []ast.PdfNode) {
	if len(node.Children()) != len(expectedChildren) {
		t.Errorf("Expected %v children, got %v", len(expectedChildren), len(node.Children()))
	}

	for i, child := range node.Children() {
		expectNode(t, child, expectedChildren[i].Type(), expectedChildren[i].Value())
	}
}

// expectNode checks that a node is of the expected type and has the expected
// value if supplied
func expectNode(t *testing.T, node ast.PdfNode, expectedType ast.Type, expectedValue interface{}) {
	if node.Type() != expectedType {
		t.Errorf("Expected node type %v, got %v", expectedType, node.Type())
	}

	if expectedValue != nil && node.Value() != expectedValue {
		t.Errorf("Expected node value %v, got %v", expectedValue, node.Value())
	}
}

type mockTokeniser struct {
	tokens  []token.Token
	current int
}

func (t *mockTokeniser) NextToken() (token.Token, error) {
	if t.current >= len(t.tokens) {
		return token.Token{}, nil
	}

	tok := t.tokens[t.current]
	t.current++
	return tok, nil
}

func (t *mockTokeniser) UnreadToken() {
	if t.current == 0 {
		panic("No tokens to unread")
	}

	t.current--
}
