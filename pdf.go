package pdf

import (
	"io"
	"os"

	"github.com/rgracey/pdf/pkg/ast"
	"github.com/rgracey/pdf/pkg/parser"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

func ParseStream(r io.Reader) (ast.PdfNode, error) {
	tokeniser := tokeniser.NewTokeniser(r)
	parser := parser.NewParser(tokeniser)
	return parser.Parse(), nil
}

func ParseFile(filename string) (ast.PdfNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ast.NewRootNode(), err
	}
	defer file.Close()

	return ParseStream(file)
}
