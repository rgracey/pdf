package tokeniser

import (
	"bufio"
	"io"

	"github.com/rgracey/pdf/pkg/token"
)

const (
	eof                = rune(0)
	carriageReturn     = rune(13)
	lineFeed           = rune(10)
	percent            = rune(37)
	leftCurlyBracket   = rune(123)
	rightCurlyBracket  = rune(125)
	leftSquareBracket  = rune(91)
	rightSquareBracket = rune(93)
	leftParenthesis    = rune(40)
	rightParenthesis   = rune(41)
	lessThan           = rune(60)
	greaterThan        = rune(62)
	slash              = rune(47)
)

type Tokeniser interface {
	PeekToken() token.Token
	NextToken() token.Token
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

func (t *StreamTokeniser) PeekToken() token.Token {
	if len(t.unreadTokens) > 0 {
		return t.unreadTokens[0]
	}

	t.unreadTokens = append(t.unreadTokens, t.getToken())
	return t.unreadTokens[0]
}

func (t *StreamTokeniser) NextToken() token.Token {
	if len(t.unreadTokens) > 0 {
		tok := t.unreadTokens[0]
		t.unreadTokens = t.unreadTokens[1:]
		t.readtokens = append(t.readtokens, tok)
		return tok
	}

	tok := t.getToken()
	t.readtokens = append(t.readtokens, tok)
	return tok
}

func (t *StreamTokeniser) UnreadToken() {
	if len(t.readtokens) == 0 {
		panic("No tokens to unread")
	}

	t.unreadTokens = append([]token.Token{t.readtokens[len(t.readtokens)-1]}, t.unreadTokens...)
	t.readtokens = t.readtokens[:len(t.readtokens)-1]
}

func (t *StreamTokeniser) getToken() token.Token {
	for isWhitespace(t.peek()) {
		t.read()
	}

	ch := t.read()

	switch {
	case ch == eof:
		return token.Token{Type: token.EOF}

	case ch == percent:
		return token.Token{Type: token.COMMENT, Value: t.readComment()}

	case ch == leftCurlyBracket:
		return token.Token{Type: token.FUNCTION_START}

	case ch == rightCurlyBracket:
		return token.Token{Type: token.FUNCTION_END}

	case ch == leftSquareBracket:
		return token.Token{Type: token.ARRAY_START}

	case ch == rightSquareBracket:
		return token.Token{Type: token.ARRAY_END}

	case ch == leftParenthesis:
		return token.Token{Type: token.STRING_START}

	case ch == rightParenthesis:
		return token.Token{Type: token.STRING_END}

	case ch == slash:
		return token.Token{Type: token.NAME, Value: t.readName()}

	case isDigit(ch):
		t.unread()
		return token.Token{Type: token.NUMBER, Value: t.readNumber()}

	case isLetter(ch):
		t.unread()
		return token.Token{Type: token.KEYWORD, Value: t.readLetters()}

	case ch == lessThan:
		if t.peek() == lessThan {
			t.read()
			return token.Token{Type: token.DICT_START}
		}

	case ch == greaterThan:
		if t.peek() == greaterThan {
			t.read()
			return token.Token{Type: token.DICT_END}
		}
	}

	return token.Token{Type: token.UNKNOWN, Value: string(ch)}
}

func (l *StreamTokeniser) readName() string {
	var name string

	for {
		ch := l.peek()

		if !isLetter(ch) && !isDigit(ch) {
			break
		}

		name += string(ch)
		l.read()
	}

	return name
}

func (l *StreamTokeniser) readComment() string {
	var comment string

	for {
		ch := l.peek()

		if ch == carriageReturn || ch == lineFeed {
			l.read()

			if ch == carriageReturn && l.peek() == lineFeed {
				l.read()
			}

			break
		}

		comment += string(ch)
		l.read()
	}

	return comment
}

func (l *StreamTokeniser) readLetters() string {
	var word string

	for {
		ch := l.peek()

		if !isLetter(ch) {
			break
		}

		word += string(ch)
		l.read()
	}

	return word
}

func (l *StreamTokeniser) readNumber() string {
	var number string

	for {
		ch := l.peek()

		if !isDigit(ch) {
			break
		}

		number += string(ch)
		l.read()
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

func (t *StreamTokeniser) peek() rune {
	ch := t.read()
	t.unread()

	return ch
}
