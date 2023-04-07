package tokeniser

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
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

func isEol(ch rune) bool {
	return ch == '\n' || ch == '\r'
}
