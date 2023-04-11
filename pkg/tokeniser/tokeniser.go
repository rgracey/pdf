package tokeniser

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/rgracey/pdf/pkg/token"
)

type Tokeniser interface {
	NextToken() (token.Token, error)
	UnreadToken()
}

type StreamTokeniser struct {
	r            *bufio.Reader
	readtokens   *Stack[token.Token] // All read tokens
	unreadTokens *Stack[token.Token] // Any read then unread tokens
}

func NewTokeniser(r io.Reader) Tokeniser {
	return &StreamTokeniser{
		r:            bufio.NewReader(r),
		readtokens:   NewStack[token.Token](3),
		unreadTokens: NewStack[token.Token](3),
	}
}

func (t *StreamTokeniser) NextToken() (token.Token, error) {
	if t.unreadTokens.Length() > 0 {
		tok := t.unreadTokens.Pop()
		t.readtokens.Push(tok)
		return tok, nil
	}

	tok, err := t.getToken()

	if err != nil {
		return token.Token{}, err
	}

	t.readtokens.Push(tok)
	return tok, nil
}

func (t *StreamTokeniser) UnreadToken() {
	if t.readtokens.Length() == 0 {
		panic("No tokens to unread")
	}

	t.unreadTokens.Push(t.readtokens.Pop())
}

func (t *StreamTokeniser) getToken() (token.Token, error) {
	ch, eof := t.read()

	if eof {
		return token.Token{Type: token.EOF}, nil
	}

	for isWhitespace(ch) {
		ch, _ = t.read()
	}

	// TODO - This is a bit of a hack
	// Instead, could add smart String() handling to AST PdfNodes and modify
	// the adding of children to stream nodes?
	// Edge case for handling stream bodies as they can contain characters
	// that will trip up further tokenisation (or make it hang)
	if t.readtokens.Length() > 0 &&
		t.readtokens.Top().Type == token.KEYWORD &&
		t.readtokens.Top().Value == "stream" {
		t.unread()
		return token.Token{Type: token.STREAM, Value: t.readStream()}, nil
	}

	switch {
	case ch == '<':
		if t.maybe('<') {
			return token.Token{
				Type:  token.DICT_START,
				Value: "<<",
			}, nil
		}

		return token.Token{
			Type:  token.DELIMITER,
			Value: "<",
		}, nil

		// TODO - Should be a hex string
		// return token.Token{
		// 	Type:  token.REGULAR_CHAR,
		// 	Value: t.readRegularCharacters(),
		// }, nil

	case ch == '>':
		if t.maybe('>') {
			return token.Token{
				Type:  token.DICT_END,
				Value: ">>",
			}, nil
		}

		return token.Token{
			Type:  token.DELIMITER,
			Value: ">",
		}, nil

	case ch == '{':
		return token.Token{
			Type:  token.FUNCTION_START,
			Value: '{',
		}, nil

	case ch == '}':
		return token.Token{
			Type:  token.FUNCTION_END,
			Value: '}',
		}, nil

	case ch == '[':
		return token.Token{
			Type:  token.ARRAY_START,
			Value: '[',
		}, nil

	case ch == ']':
		return token.Token{
			Type:  token.ARRAY_END,
			Value: ']',
		}, nil

	case ch == '(':
		return token.Token{
			Type:  token.STRING_LITERAL,
			Value: t.readStringLiteral(),
		}, nil

	case ch == ')':
		return token.Token{
			Type:  token.DELIMITER,
			Value: ')',
		}, nil

	case ch == '%':
		return token.Token{Type: token.COMMENT, Value: t.readComment()}, nil

	case ch == '/':
		return token.Token{Type: token.NAME, Value: t.readRegularCharacters()}, nil

	default:
		t.unread()
		tmp := t.readRegularCharacters()

		switch {
		case tmp == "true":
			return token.Token{Type: token.BOOLEAN, Value: true}, nil

		case tmp == "false":
			return token.Token{Type: token.BOOLEAN, Value: false}, nil

		case isInteger(tmp):
			num, err := strconv.ParseInt(tmp, 10, 64)

			if err != nil {
				return token.Token{}, err
			}

			return token.Token{Type: token.NUMBER_INTEGER, Value: num}, nil

		case isFloat(tmp):
			num, err := strconv.ParseFloat(tmp, 64)

			if err != nil {
				return token.Token{}, err
			}

			return token.Token{Type: token.NUMBER_FLOAT, Value: num}, nil
		}

		return token.Token{Type: token.KEYWORD, Value: tmp}, nil
	}
}

func (l *StreamTokeniser) readComment() string {
	sb := strings.Builder{}

	for {
		ch, _ := l.read()

		if ch == '\r' || ch == '\n' {
			l.read()
			c, _ := l.read()

			if ch == '\r' && c != '\n' {
				l.unread()
			}

			break
		}

		sb.WriteRune(ch)
	}

	return sb.String()
}

func (l *StreamTokeniser) readStringLiteral() string {
	sb := strings.Builder{}

	for {
		ch, _ := l.read()

		if ch == ')' {
			break
		}

		if ch == '\\' {
			ch, _ = l.read()
		}

		sb.WriteRune(ch)
	}

	return sb.String()
}

// readStream reads the stream body until it finds the endstream keyword
func (l *StreamTokeniser) readStream() string {
	sb := strings.Builder{}

	for {
		ch, _ := l.read()

		sb.WriteString(string(ch))

		if strings.HasSuffix(sb.String(), "endstream") {
			break
		}
	}

	stream := sb.String()
	return stream[:len(stream)-9] // trim "endstream"
}

func (l *StreamTokeniser) readRegularCharacters() string {
	sb := strings.Builder{}

	for {
		ch, _ := l.read()

		if isDelimiter(ch) || isWhitespace(ch) {
			l.unread()
			break
		}

		sb.WriteRune(ch)
	}

	return sb.String()
}

func (t *StreamTokeniser) unread() {
	t.r.UnreadRune()
}

func (t *StreamTokeniser) read() (rune, bool) {
	ch, _, err := t.r.ReadRune()

	if err != nil {
		return 0, true
	}

	return ch, false
}

func (t *StreamTokeniser) maybe(ch rune) bool {
	next, _ := t.read()

	if next != ch {
		t.unread()
		return false
	}

	return true
}
