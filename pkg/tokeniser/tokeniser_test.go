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

	expectTokens(t, tokeniser, expectedTokens)
}

func TestTokeniser_UnreadToken(t *testing.T) {
	pdf := strings.NewReader("1 0 obj")

	tokeniser := tokeniser.NewTokeniser(pdf)

	for i := 0; i < 3; i++ {
		tokeniser.NextToken()
	}

	for i := 0; i < 3; i++ {
		tokeniser.UnreadToken()
	}

	expected := []token.Token{
		{Type: token.NUMBER_INTEGER, Value: int64(1)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "obj"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_SkipsWhitespace(t *testing.T) {
	pdf := strings.NewReader("\n \x00           1 0 obj")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.NUMBER_INTEGER, Value: int64(1)},
		{Type: token.NUMBER_INTEGER, Value: int64(0)},
		{Type: token.KEYWORD, Value: "obj"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesStreams(t *testing.T) {
	pdf := strings.NewReader("stream\nThis is a stream\nendstream")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.KEYWORD, Value: "stream"},
		{Type: token.STREAM, Value: "This is a stream"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesStrings(t *testing.T) {
	pdf := strings.NewReader("(This is a string)")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.STRING_LITERAL, Value: "This is a string"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesComments(t *testing.T) {
	pdf := strings.NewReader("%This is a comment\n")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.COMMENT, Value: "This is a comment"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesNames(t *testing.T) {
	pdf := strings.NewReader("/Catalog")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.NAME, Value: "Catalog"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesDict(t *testing.T) {
	pdf := strings.NewReader("<< /Type /Catalog /Outlines 2 0 R /Pages 3 0 R >>")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
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
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesArray(t *testing.T) {
	pdf := strings.NewReader("[ 1 2 3 ]")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.ARRAY_START, Value: "["},
		{Type: token.NUMBER_INTEGER, Value: int64(1)},
		{Type: token.NUMBER_INTEGER, Value: int64(2)},
		{Type: token.NUMBER_INTEGER, Value: int64(3)},
		{Type: token.ARRAY_END, Value: "]"},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesBooleans(t *testing.T) {
	pdf := strings.NewReader("true false")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.BOOLEAN, Value: true},
		{Type: token.BOOLEAN, Value: false},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesFloat(t *testing.T) {
	pdf := strings.NewReader("1.0000 1.5 1.738478")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.NUMBER_FLOAT, Value: float64(1.0000)},
		{Type: token.NUMBER_FLOAT, Value: float64(1.5)},
		{Type: token.NUMBER_FLOAT, Value: float64(1.738478)},
	}

	expectTokens(t, tokeniser, expected)
}

func TestTokeniser_HandlesInteger(t *testing.T) {
	pdf := strings.NewReader("1 2 3")

	tokeniser := tokeniser.NewTokeniser(pdf)

	expected := []token.Token{
		{Type: token.NUMBER_INTEGER, Value: int64(1)},
		{Type: token.NUMBER_INTEGER, Value: int64(2)},
		{Type: token.NUMBER_INTEGER, Value: int64(3)},
	}

	expectTokens(t, tokeniser, expected)
}

func expectTokens(t *testing.T, tok tokeniser.Tokeniser, expected []token.Token) {
	for _, expectedToken := range expected {
		actual, err := tok.NextToken()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if expectedToken.Type != actual.Type {
			t.Errorf("Expected token type %v, got %v", expectedToken.Type, actual.Type)
		}

		if expectedToken.Value != actual.Value {
			t.Errorf("Expected token value \"%v\", got \"%v\"", expectedToken.Value, actual.Value)
		}
	}
}
