package pdf

import (
	"io"
	"os"

	"github.com/rgracey/pdf/pkg/document"
	"github.com/rgracey/pdf/pkg/parser"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

func ParseStream(r io.Reader) (document.Pdf, error) {
	tokeniser := tokeniser.NewTokeniser(r)
	parser := parser.NewParser(tokeniser)
	doc := parser.Parse()

	return doc, nil
}

func ParseFile(filename string) (document.Pdf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return document.Pdf{}, err
	}
	defer file.Close()

	return ParseStream(file)
}
