package tokeniser_test

import (
	"strings"
	"testing"

	"github.com/rgracey/pdf/pkg/token"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

func TestTokeniser_NextToken(t *testing.T) {
	pdf := strings.NewReader(`
		%PDF-1.3
		%����

		1 0 obj
		<<
		/Type /Catalog
		/Outlines 2 0 R
		/Pages 3 0 R
		>>
		endobj

		3 0 obj
		<<
		/Type /Pages
		/Count 2
		/Kids [ 4 0 R 6 0 R ]
		/Something (This is a string)
		>>
		endobj
		`)

	tokeniser := tokeniser.NewTokeniser(pdf)

	expectedTokens := []token.Token{
		{Type: token.COMMENT, Value: "PDF-1.3"},
		{Type: token.COMMENT, Value: "����"},
		{Type: token.NUMBER_INTEGER, Value: int64(1)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "obj"},
		{Type: token.DICT_START, Value: "<<"},
		{Type: token.NAME, Value: "Type"},
		{Type: token.NAME, Value: "Catalog"},
		{Type: token.NAME, Value: "Outlines"},
		{Type: token.NUMBER_INTEGER, Value: int64(2)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "R"},
		{Type: token.NAME, Value: "Pages"},
		{Type: token.NUMBER_INTEGER, Value: int64(3)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "R"},
		{Type: token.DICT_END, Value: ">>"},
		{Type: token.KEYWORD, Value: "endobj"},
		{Type: token.NUMBER_INTEGER, Value: int64(3)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "obj"},
		{Type: token.DICT_START, Value: "<<"},
		{Type: token.NAME, Value: "Type"},
		{Type: token.NAME, Value: "Pages"},
		{Type: token.NAME, Value: "Count"},
		{Type: token.NUMBER_INTEGER, Value: int64(2)},
		{Type: token.NAME, Value: "Kids"},
		{Type: token.ARRAY_START, Value: "["},
		{Type: token.NUMBER_INTEGER, Value: int64(4)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "R"},
		{Type: token.NUMBER_INTEGER, Value: int64(6)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "R"},
		{Type: token.ARRAY_END, Value: "]"},
		{Type: token.NAME, Value: "Something"},
		{Type: token.STRING_LITERAL, Value: "This is a string"},
		{Type: token.DICT_END, Value: ">>"},
		{Type: token.KEYWORD, Value: "endobj"},
	}

	for _, expectedToken := range expectedTokens {
		tok, err := tokeniser.NextToken()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		compareToken(t, expectedToken, tok)
	}
}

func TestTokeniser_UnreadToken(t *testing.T) {
	pdf := strings.NewReader("1 0 obj")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := make([]token.Token, 3)

	for i := 0; i < 3; i++ {
		tok, _ := tokeniser.NextToken()
		expected[i] = tok
	}

	for i := 0; i < 3; i++ {
		tokeniser.UnreadToken()
	}

	for i := 0; i < 3; i++ {
		tok, _ := tokeniser.NextToken()
		compareToken(t, expected[i], tok)
	}
}

func compareToken(t *testing.T, expected token.Token, actual token.Token) {
	if expected.Type != actual.Type {
		t.Errorf("Expected token type %v, got %v", expected.Type, actual.Type)
	}

	if expected.Value != actual.Value {
		t.Errorf("Expected token value %v, got %v", expected.Value, actual.Value)
	}
}
