package token

import "fmt"

// Type represents the type of a token
type Type int

const (
	EOF Type = iota

	COMMENT // %

	NUMBER_INTEGER // 123 -123 +123
	NUMBER_FLOAT   // 123.456 -123.456 +123.456

	BOOLEAN // true false

	DICT_START // <<
	DICT_END   // >>

	ARRAY_START // [
	ARRAY_END   // ]

	STRING_LITERAL // (the string)

	FUNCTION_START // {
	FUNCTION_END   // }

	NAME // /Name

	KEYWORD // obj endobj R stream endstream xref trailer startxref

	STREAM // The body of a stream

	DELIMITER    // Any delmiter character
	REGULAR_CHAR // Any non whitespace or delimiter character
)

// Token represents a grouping of characters that have a meaning.
type Token struct {
	Type  Type
	Value interface{}
}

func (t Token) String() string {
	var tokenType = "UNKNOWN"

	switch t.Type {
	case EOF:
		tokenType = "EOF"
	case COMMENT:
		tokenType = "COMMENT"
	case BOOLEAN:
		tokenType = "BOOLEAN"
	case NUMBER_INTEGER:
		tokenType = "INTEGER"
	case NUMBER_FLOAT:
		tokenType = "FLOAT"
	case DICT_START:
		tokenType = "DICT_START"
	case DICT_END:
		tokenType = "DICT_END"
	case ARRAY_START:
		tokenType = "ARRAY_START"
	case ARRAY_END:
		tokenType = "ARRAY_END"
	case STRING_LITERAL:
		tokenType = "STRING_LITERAL"
	case FUNCTION_START:
		tokenType = "FUNCTION_START"
	case FUNCTION_END:
		tokenType = "FUNCTION_END"
	case NAME:
		tokenType = "NAME"
	case KEYWORD:
		tokenType = "KEYWORD"
	}

	return fmt.Sprintf("TYPE: %s, VALUE: %v", tokenType, t.Value)
}
