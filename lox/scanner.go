package lox

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	Source  string
	Tokens  []Token
	current int
	start   int
	line    int
}

func (s Scanner) IsAtEnd() bool {
	return s.current >= len(s.Source)
}

func NewScanner(source string) *Scanner {
	var tokens []Token

	return &Scanner{
		Source:  source,
		Tokens:  tokens,
		current: 0,
		start:   0,
		line:    1,
	}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"if":     IF,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (s Scanner) ScanTokens() {
	for s.IsAtEnd() {
		s.start = s.current
		s.ScanToken()
	}
	s.Tokens = append(s.Tokens, *NewToken(EOF, "", nil, s.line))
}

func (s Scanner) advance() byte {
	out := s.Source[s.current]
	s.current++
	return out
}
func (s Scanner) ScanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.IsAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	s.addToken(IDENTIFIER)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	val, err := strconv.ParseFloat(s.Source[s.start:s.current], 32)
	if err != nil {
		error(s.line, "Not a float")
		return
	}
	s.addTokenLiteral(NUMBER, val)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.IsAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}
	if s.IsAtEnd() {
		error(s.line, "Unterminated string.")
		return
	}
	s.advance()

	value := s.Source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, value)
}

func (s *Scanner) peek() byte {
	if s.IsAtEnd() {
		return '\000'
	}
	return s.Source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return '\000'
	}
	return s.Source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.IsAtEnd() {
		return false
	}

	if s.Source[s.current] != expected {
		return false
	}
	s.current++
	return true

}

func (s *Scanner) addToken(ttype TokenType) {
	s.addTokenLiteral(ttype, nil)
}

func (s *Scanner) addTokenLiteral(ttype TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, *NewToken(ttype, text, literal, s.line))
}
