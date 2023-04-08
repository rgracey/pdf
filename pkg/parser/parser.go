package parser

import (
	"fmt"
	"strings"

	"github.com/rgracey/pdf/pkg/document"
	"github.com/rgracey/pdf/pkg/token"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

type Parser struct {
	tokeniser tokeniser.Tokeniser
	pdfDoc    document.Pdf
}

func NewParser(tokeniser tokeniser.Tokeniser) *Parser {
	return &Parser{
		tokeniser: tokeniser,
	}
}

func (p *Parser) Parse() document.Pdf {
	p.pdfDoc = document.Pdf{}

	t, err := p.tokeniser.NextToken()

	if err != nil {
		panic(err)
	}

	if t.Type != token.COMMENT || !strings.HasPrefix(t.Value.(string), "PDF-") {
		panic("Expected PDF version comment")
	}

	p.pdfDoc.SetVersion(t.Value.(string))

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.EOF {
			break
		}

		p.tokeniser.UnreadToken()

		obj := p.parseObject()

		switch obj.Type {
		case document.INDIRECT_OBJECT:
			p.pdfDoc.AddObject(obj)

		case document.XREFS:
			for _, xref := range obj.Data.([]document.Xref) {
				p.pdfDoc.AddXref(xref)
			}

		case document.DICT:
			p.pdfDoc.SetTrailer(obj.Data.(document.Dict))
		}

	}

	return p.pdfDoc
}

func (p *Parser) parseObject() document.Object {
	tok, err := p.tokeniser.NextToken()

	if err != nil {
		panic(err)
	}

	switch tok.Type {
	case token.KEYWORD:
		switch tok.Value {
		case "stream":
			return p.parseStream()

		case "xref":
			return document.Object{
				Type: document.XREFS,
				Data: p.parseXrefTable(),
			}

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
		return document.Object{
			Type: document.STRING,
			Data: tok.Value,
		}

	case token.NUMBER_FLOAT:
		return document.Object{
			Type: document.NUMBER_FLOAT,
			Data: tok.Value,
		}

	case token.NUMBER_INTEGER:
		gen, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if gen.Type != token.NUMBER_INTEGER {
			p.tokeniser.UnreadToken()
			return document.Object{
				Type: document.NUMBER_INTEGER,
				Data: tok.Value,
			}
		}

		keyword, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if keyword.Type != token.KEYWORD {
			p.tokeniser.UnreadToken()
			p.tokeniser.UnreadToken()
			return document.Object{
				Type: document.NUMBER_INTEGER,
				Data: tok.Value,
			}
		}

		switch keyword.Value {
		case "obj":
			data := p.parseObject()

			t, err := p.tokeniser.NextToken()

			if err != nil {
				panic(err)
			}

			var stream interface{}

			if t.Type == token.KEYWORD && t.Value == "stream" {
				p.tokeniser.UnreadToken()
				stream = p.parseObject()
				t, err = p.tokeniser.NextToken()

				if err != nil {
					panic(err)
				}
			}

			if t.Type != token.KEYWORD || t.Value != "endobj" {
				panic("Expected endobj")
			}

			return document.Object{
				Type: document.INDIRECT_OBJECT,
				Ref: document.ObjectRef{
					Id:         int(tok.Value.(int64)),
					Generation: int(gen.Value.(int64)),
				},
				Header: data,
				Data:   stream,
			}

		case "R":
			return document.Object{
				Type: document.OBJECT_REF,
				Ref: document.ObjectRef{
					Id:         int(tok.Value.(int64)),
					Generation: int(gen.Value.(int64)),
				},
			}
		}
	}

	return document.Object{
		Type: document.UNKNOWN,
		Data: tok.Value,
	}
}

func (p *Parser) parseXrefTable() []document.Xref {
	id, err := p.tokeniser.NextToken()

	if err != nil {
		panic(err)
	}

	if id.Type != token.NUMBER_INTEGER {
		panic("Expected xref id")
	}

	totalObjects, err := p.tokeniser.NextToken()

	if err != nil {
		panic(err)
	}

	if totalObjects.Type != token.NUMBER_INTEGER {
		panic("Expected xref total objects")
	}

	tot := totalObjects.Value.(int64)

	xrefs := make([]document.Xref, tot)

	for i := int64(0); i < tot; i++ {
		offset, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if offset.Type != token.NUMBER_INTEGER {
			panic("Expected xref offset")
		}

		gen, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if gen.Type != token.NUMBER_INTEGER {
			panic("Expected xref gen")
		}

		used, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if used.Type != token.KEYWORD || (used.Value != "n" && used.Value != "f") {
			panic("Expected xref used")
		}

		u := false

		if used.Value == "n" {
			u = true
		}

		xrefs[i] = document.Xref{
			Offset:     offset.Value.(int64),
			Generation: gen.Value.(int64),
			Used:       u,
		}
	}

	return xrefs
}

func (p *Parser) parseStream() document.Object {
	stream := ""

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.KEYWORD && t.Value == "endstream" {
			break
		}

		if t.Value != nil {
			stream += fmt.Sprintf("%v", t.Value)
		}
	}

	return document.Object{
		Type: document.STREAM,
		Data: stream,
	}
}

func (p *Parser) parseFunction() document.Object {
	function := ""

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.FUNCTION_END {
			break
		}

		if t.Value != nil {
			function += t.Value.(string)
		}
	}

	return document.Object{
		Type: document.FUNCTION,
		Data: function,
	}
}

func (p *Parser) parseArray() document.Object {
	arr := []document.Object{}

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.ARRAY_END {
			break
		}

		arr = append(arr, p.parseObject())
	}

	return document.Object{
		Type: document.ARRAY,
		Data: arr,
	}
}

// TODO - correct return type
func (p *Parser) parseDictionary() document.Object {
	dict := document.Dict{}

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.DICT_END {
			break
		}

		if t.Type != token.NAME {
			report(token.Token{Type: token.NAME}, t)
		}

		key := t.Value.(string)

		dict[key] = p.parseObject()
	}

	return document.Object{
		Type: document.DICT,
		Data: dict,
	}
}

func report(expected token.Token, actual token.Token) {
	panic(fmt.Sprintf("\nExpected:\n	%s\nActual:\n	%s", expected, actual))
}
