package pdf

import (
	"io"
	"os"

	"github.com/rgracey/pdf/pkg/parser"
	"github.com/rgracey/pdf/pkg/tokeniser"
)

func ParseStream(r io.Reader) error {
	tokeniser := tokeniser.NewTokeniser(r)
	parser := parser.NewParser(tokeniser)
	parser.Parse()

	return nil
}

func ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return ParseStream(file)
}
