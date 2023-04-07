package tokeniser

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDelimiter returns true if the given character is a delimiter as per the PDF
// spec.
//
// Delimiters:
//   - 0x28 (left parenthesis) (
//   - 0x29 (right parenthesis) )
//   - 0x3C (less than) <
//   - 0x3E (greater than) >
//   - 0x5B (left square bracket) [
//   - 0x5D (right square bracket) ]
//   - 0x7B (left brace) {
//   - 0x7D (right brace) }
//   - 0x2F (forward slash) /
//   - 0x25 (percent sign) %
func isDelimiter(ch rune) bool {
	switch ch {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}

	return false
}

// isWhitespace returns true if the given character is a whitespace character as
// per the PDF spec.
//
// Whitespace characters:
//   - 0x00 (null)
//   - 0x09 (horizontal tab) \t
//   - 0x0A (line feed) \n
//   - 0x0C (form feed) \f
//   - 0x0D (carriage return) \r
//   - 0x20 (space)
func isWhitespace(ch rune) bool {
	switch ch {
	case '\x00', '\t', '\n', '\f', '\r', ' ':
		return true
	}

	return false
}
