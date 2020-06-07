package scanner

import (
	"strconv"

	"github.com/christiansakai/glox/token"
)

type Scanner struct {
	source string
	tokens []token.Token

	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens(onError func(int, string)) []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken(onError)
	}

	s.tokens = append(s.tokens, token.Token{
		Type:    token.EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})

	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken(onError func(int, string)) {
	ch := s.advance()

	switch ch {
	case '(':
		s.addToken(token.LEFT_PAREN)

	case ')':
		s.addToken(token.RIGHT_PAREN)

	case '{':
		s.addToken(token.LEFT_BRACE)

	case '}':
		s.addToken(token.RIGHT_BRACE)

	case ',':
		s.addToken(token.COMMA)

	case '.':
		s.addToken(token.DOT)

	case '-':
		s.addToken(token.MINUS)

	case '+':
		s.addToken(token.PLUS)

	case ';':
		s.addToken(token.SEMICOLON)

	case '*':
		s.addToken(token.STAR)

	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}

	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}

	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}

	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}

	// Ignore whitespaces
	case ' ':
	case '\r':
	case '\t':

	case '\n':
		s.line += 1

	case '"':
		s.string(onError)

	default:
		if s.isDigit(ch) {
			s.number()
		} else if s.isAlpha(ch) {
			s.identifier()
		} else {
			onError(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) advance() byte {
	s.current += 1
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType token.Type, literal ...interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current += 1
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}

	return s.source[s.current]
}

func (s *Scanner) string(onError func(int, string)) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
			s.advance()
		}
	}

	// Unterminated string
	if s.isAtEnd() {
		onError(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *Scanner) isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	num, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(token.NUMBER, num)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}

	return s.source[s.current+1]
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// See if the identifier is a reserved word
	text := s.source[s.start:s.current]

	tokenType, ok := keywords[text]
	if !ok {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType)
}

func (s *Scanner) isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (s *Scanner) isAlphaNumeric(ch byte) bool {
	return s.isAlpha(ch) || s.isDigit(ch)
}
