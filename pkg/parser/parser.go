package parser

import (
	"fmt"

	"github.com/rgracey/pdf/pkg/ast"
	"github.com/rgracey/pdf/pkg/token"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

type Parser struct {
	tokeniser   tokeniser.Tokeniser
	ast         ast.PdfNode
	current     ast.PdfNode
	parentStack []ast.PdfNode
}

func NewParser(tokeniser tokeniser.Tokeniser) *Parser {
	root := ast.NewRootNode()

	return &Parser{
		tokeniser: tokeniser,
		ast:       root,
		current:   root,
	}
}

func (p *Parser) Parse() ast.PdfNode {
	for {
		tok, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if tok.Type == token.EOF {
			break
		}

		switch tok.Type {
		case token.KEYWORD:
			switch tok.Value {
			case "obj", "R":
				// Try to pop the last 2 children off the current node
				// If they're both integers, then we have an indirect object or
				// a reference to an indirect object
				if len(p.current.Children()) < 2 {
					panic("Unexpected obj")
				}

				gen := p.current.Children()[len(p.current.Children())-1]
				id := p.current.Children()[len(p.current.Children())-2]

				if gen.Type() != ast.INTEGER || id.Type() != ast.INTEGER {
					panic("Unexpected obj")
				}

				p.current.RemoveChild(len(p.current.Children()) - 1)
				p.current.RemoveChild(len(p.current.Children()) - 1)

				switch tok.Value {
				case "obj":
					obj := ast.NewIndirectObjectNode(
						id.Value().(int64),
						gen.Value().(int64),
					)
					p.push(obj)

				case "R":
					refObj := ast.NewObjectRefNode(
						id.Value().(int64),
						gen.Value().(int64),
					)
					p.current.AddChild(refObj)
				}

			case "endobj":
				p.pop()

			case "stream":
				stream := p.parseStream()
				p.push(ast.NewStreamNode(stream))

			case "endstream":
				p.pop()

			case "xref":
				p.push(ast.NewXRefsNode())

			case "startxref":
				p.pop()

			case "trailer":
				// Do no special parsing for trailer for now
				// maybe introduce a trailer node type to more
				// easily access the trailer dictionary?
			}

		case token.NAME:
			p.current.AddChild(ast.NewNameNode(tok.Value.(string)))

		// TODO - Handle null otherwise dictionaries could have trouble
		// (uneven number of children)
		// Could potentially just return null as default?
		// case token.NULL:

		case token.DICT_START:
			p.push(ast.NewDictNode())

		case token.DICT_END:
			p.pop()

		case token.ARRAY_START:
			p.push(ast.NewArrayNode())

		case token.ARRAY_END:
			p.pop()

		case token.STRING_LITERAL:
			p.current.AddChild(ast.NewStringNode(tok.Value.(string)))

		case token.NUMBER_FLOAT:
			p.current.AddChild(ast.NewFloatNode(tok.Value.(float64)))

		case token.NUMBER_INTEGER:
			p.current.AddChild(ast.NewIntegerNode(tok.Value.(int64)))
		}

	}

	return p.ast
}

// push pushes a node onto the parent stack and sets it as the current node
func (p *Parser) push(node ast.PdfNode) {
	p.current.AddChild(node)
	p.parentStack = append(p.parentStack, p.current)
	p.current = node
}

// pop pops a node off the parent stack and sets it as the current node
func (p *Parser) pop() {
	if len(p.parentStack) == 0 {
		panic("Unexpected end")
	}
	p.current = p.parentStack[len(p.parentStack)-1]
	p.parentStack = p.parentStack[:len(p.parentStack)-1]
}

func (p *Parser) parseStream() string {
	stream := ""

	for {
		t, err := p.tokeniser.NextToken()

		if err != nil {
			panic(err)
		}

		if t.Type == token.KEYWORD && t.Value == "endstream" {
			p.tokeniser.UnreadToken()
			break
		}

		// TODO - fix performance, this is slow
		if t.Value != nil {
			stream += fmt.Sprintf("%v", t.Value)
		}
	}

	return stream
}
