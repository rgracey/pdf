package token

import "fmt"

type Type int

const (
	EOF Type = iota

	COMMENT // %

	NUMBER // 123

	DICT_START // <<
	DICT_END   // >>

	ARRAY_START // [
	ARRAY_END   // ]

	STRING_START // (
	STRING_END   // )

	FUNCTION_START // {
	FUNCTION_END   // }

	NAME // /Name

	KEYWORD // obj endobj stream endstream R

	UNKNOWN // Anything else. Possibly the body of a stream (if encoded)
)

type Token struct {
	Type  Type
	Value interface{}
}

type IndiectObject struct {
	Id  uint
	Gen uint
	// TODO - more
}

// String method
func (t Token) String() string {
	var tokenType = "UNKNOWN"

	switch t.Type {
	case EOF:
		tokenType = "EOF"
	case COMMENT:
		tokenType = "COMMENT"
	case NUMBER:
		tokenType = "NUMBER"
	case DICT_START:
		tokenType = "DICT_START"
	case DICT_END:
		tokenType = "DICT_END"
	case ARRAY_START:
		tokenType = "ARRAY_START"
	case ARRAY_END:
		tokenType = "ARRAY_END"
	case STRING_START:
		tokenType = "STRING_START"
	case STRING_END:
		tokenType = "STRING_END"
	case NAME:
		tokenType = "NAME"
	case KEYWORD:
		tokenType = "KEYWORD"
	case UNKNOWN:
		tokenType = "UNKNOWN"
	}

	return fmt.Sprintf("TYPE: %s, VALUE: %v", tokenType, t.Value)
}
