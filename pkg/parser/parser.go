package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rgracey/pdf/pkg/token"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

type Parser struct {
	tokeniser tokeniser.Tokeniser
}

func NewParser(tokeniser tokeniser.Tokeniser) *Parser {
	return &Parser{
		tokeniser: tokeniser,
	}
}

func (p *Parser) Parse() {
	t := p.tokeniser.NextToken()

	if t.Type != token.COMMENT || !strings.HasPrefix(t.Value.(string), "PDF-") {
		panic("Expected PDF version comment")
	}

	fmt.Println(t.Value)

	for {
		t := p.tokeniser.PeekToken()

		if t.Type == token.EOF {
			break
		}

		obj := p.parseObject()

		fmt.Println(obj)
	}
}

type ObjectReference struct {
	Id  interface{}
	Gen interface{}
}

type Object struct {
	objRef ObjectReference
	data   interface{}
	stream interface{}
}

func (p *Parser) parseObject() interface{} {
	tok := p.tokeniser.NextToken()

	switch tok.Type {
	case token.KEYWORD:
		switch tok.Value {
		case "stream":
			stream := ""
			for {
				t := p.tokeniser.NextToken()
				// Need to check for newline after endstream? In case endstream is in the stream text?
				if t.Type == token.KEYWORD && t.Value == "endstream" {
					break
				}
				if t.Value != nil {
					stream += t.Value.(string)
				}
			}

			return stream

		case "xref":
			id := p.tokeniser.NextToken()
			if id.Type != token.NUMBER {
				panic("Expected xref id")
			}

			totalObjects := p.tokeniser.NextToken()
			if totalObjects.Type != token.NUMBER {
				panic("Expected xref total objects")
			}

			tot, err := strconv.Atoi(totalObjects.Value.(string))

			if err != nil {
				panic(err)
			}

			xrefs := make([]interface{}, tot)

			for i := 0; i < tot; i++ {
				offset := p.tokeniser.NextToken()
				if offset.Type != token.NUMBER {
					panic("Expected xref offset")
				}

				gen := p.tokeniser.NextToken()
				if gen.Type != token.NUMBER {
					panic("Expected xref gen")
				}

				used := p.tokeniser.NextToken()

				if used.Type != token.KEYWORD || (used.Value != "n" && used.Value != "f") {
					panic("Expected xref used")
				}

				xrefs[i] = map[string]interface{}{
					"offset": offset.Value,
					"gen":    gen.Value,
					"used":   used.Value,
				}
			}

			return xrefs

		case "trailer":
			p.tokeniser.NextToken()
			return p.parseDictionary()
		}

	case token.DICT_START:
		return p.parseDictionary()

	case token.ARRAY_START:
		return p.parseArray()

	case token.FUNCTION_START:
		return p.parseFunction()

	case token.STRING_LITERAL:
		return tok.Value

	case token.NUMBER:
		gen := p.tokeniser.NextToken()

		if gen.Type != token.NUMBER {
			p.tokeniser.UnreadToken()
			return tok.Value
		}

		keyword := p.tokeniser.NextToken()

		if keyword.Type != token.KEYWORD {
			p.tokeniser.UnreadToken()
			p.tokeniser.UnreadToken()
			return tok.Value
		}

		switch keyword.Value {
		case "obj":
			data := p.parseObject()

			t := p.tokeniser.NextToken()

			var stream interface{}

			if t.Type == token.KEYWORD && t.Value == "stream" {
				p.tokeniser.UnreadToken()
				stream = p.parseObject()
				t = p.tokeniser.NextToken()
			}

			if t.Type != token.KEYWORD || t.Value != "endobj" {
				panic("Expected endobj")
			}

			return Object{
				objRef: ObjectReference{
					Id:  tok.Value,
					Gen: gen.Value,
				},
				data:   data,
				stream: stream,
			}

		case "R":
			return ObjectReference{
				Id:  tok.Value,
				Gen: gen.Value,
			}
		}
	}

	return tok.Value
}

func (p *Parser) parseFunction() interface{} {
	function := ""

	for {
		t := p.tokeniser.NextToken()
		if t.Type == token.FUNCTION_END {
			break
		}

		if t.Value != nil {
			function += t.Value.(string)
		}
	}

	return function
}

func (p *Parser) parseArray() []interface{} {
	arr := []interface{}{}

	for {
		t := p.tokeniser.NextToken()
		if t.Type == token.ARRAY_END {
			break
		}

		arr = append(arr, p.parseObject())
	}

	return arr
}

// TODO - correct return type
func (p *Parser) parseDictionary() interface{} {
	dict := map[string]interface{}{}

	for {
		t := p.tokeniser.NextToken()
		if t.Type == token.DICT_END {
			break
		}

		if t.Type != token.NAME {
			report(token.Token{Type: token.NAME}, t)
		}

		key := t.Value.(string)

		dict[key] = p.parseObject()
	}

	return dict
}

func report(expected token.Token, actual token.Token) {
	panic(fmt.Sprintf("\nExpected:\n	%s\nActual:\n	%s", expected, actual))
}
