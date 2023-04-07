package tokeniser

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rgracey/pdf/pkg/token"
)

const (
	eof = rune(0)
)

type Tokeniser interface {
	NextToken() (token.Token, error)
	UnreadToken()
}

type StreamTokeniser struct {
	r            *bufio.Reader
	readtokens   []token.Token // All read tokens
	unreadTokens []token.Token // Any read then unread tokens
}

func NewTokeniser(r io.Reader) Tokeniser {
	return &StreamTokeniser{
		r: bufio.NewReader(r),
	}
}

func (t *StreamTokeniser) NextToken() (token.Token, error) {
	if len(t.unreadTokens) > 0 {
		tok := t.unreadTokens[0]
		t.unreadTokens = t.unreadTokens[1:]
		t.readtokens = append(t.readtokens, tok)
		return tok, nil
	}

	tok, err := t.getToken()

	if err != nil {
		return token.Token{}, err
	}

	t.readtokens = append(t.readtokens, tok)
	return tok, nil
}

func (t *StreamTokeniser) UnreadToken() {
	if len(t.readtokens) == 0 {
		panic("No tokens to unread")
	}

	t.unreadTokens = append([]token.Token{t.readtokens[len(t.readtokens)-1]}, t.unreadTokens...)
	t.readtokens = t.readtokens[:len(t.readtokens)-1]
}

func (t *StreamTokeniser) getToken() (token.Token, error) {
	ch := t.read()

	for isWhitespace(ch) {
		ch = t.read()
	}

	switch {
	case ch == eof:
		return token.Token{Type: token.EOF}, nil

	case ch == '<':
		if t.maybe('<') {
			return token.Token{Type: token.DICT_START}, nil
		}

		// TODO - Should be a hex string
		return token.Token{
			Type:  token.REGULAR_CHAR,
			Value: t.readRegularCharacters(),
		}, nil

	case ch == '>':
		if t.maybe('>') {
			return token.Token{Type: token.DICT_END}, nil
		}

		// TODO - is this correct?
		return token.Token{
			Type:  token.KEYWORD,
			Value: t.readRegularCharacters(),
		}, nil

	case ch == '[':
		return token.Token{Type: token.ARRAY_START}, nil

	case ch == ']':
		return token.Token{Type: token.ARRAY_END}, nil

	case ch == '(':
		return token.Token{
			Type:  token.STRING_LITERAL,
			Value: t.readStringLiteral(),
		}, nil

	case ch == '%':
		return token.Token{Type: token.COMMENT, Value: t.readComment()}, nil

	case ch == '/':
		return token.Token{Type: token.NAME, Value: t.readRegularCharacters()}, nil

	// TODO - more rigorous number parsing
	case isDigit(ch):
		t.unread()
		return token.Token{Type: token.NUMBER, Value: t.readNumber()}, nil

	case isLetter(ch):
		t.unread()
		return token.Token{Type: token.KEYWORD, Value: t.readRegularCharacters()}, nil
	}

	return token.Token{}, fmt.Errorf("unexpected character: %c", ch)
}

func (l *StreamTokeniser) readComment() string {
	var comment string

	for {
		ch := l.read()

		if ch == '\r' || ch == '\n' {
			l.read()

			if ch == '\r' && l.read() != '\n' {
				l.unread()
			}

			break
		}

		comment += string(ch)
	}

	return comment
}

func (l *StreamTokeniser) readStringLiteral() string {
	var literal string

	for {
		ch := l.read()

		if ch == ')' {
			break
		}

		if ch == '\\' {
			ch = l.read()
		}

		literal += string(ch)
	}

	return literal
}

func (l *StreamTokeniser) readRegularCharacters() string {
	var characters string

	for {
		ch := l.read()

		if isDelimiter(ch) || isWhitespace(ch) {
			l.unread()
			break
		}

		characters += string(ch)
	}

	return characters
}

func (l *StreamTokeniser) readNumber() string {
	var number string

	for {
		ch := l.read()

		if !isDigit(ch) {
			l.unread()
			break
		}

		number += string(ch)
	}

	return number
}

func (t *StreamTokeniser) unread() {
	t.r.UnreadRune()
}

func (t *StreamTokeniser) read() rune {
	ch, _, err := t.r.ReadRune()

	if err != nil {
		return eof
	}

	return ch
}

func (t *StreamTokeniser) maybe(ch rune) bool {
	next := t.read()

	if next != ch {
		t.unread()
		return false
	}

	return true
}
