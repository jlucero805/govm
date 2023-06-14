package lexer

import "unicode"

type Token struct {
	Type   string
	Lexeme string
}

func newToken(tokenType string, lexeme string) Token {
	return Token{Type: tokenType, Lexeme: lexeme}
}

// Lexer performs our lexical analysis/scanning
type Lexer struct {
	input        []rune
	char         rune // current char under examination
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	tokens       []Token
}

func (l *Lexer) addToken(tokenType string) {
	l.tokens = append(l.tokens, Token{Type: tokenType, Lexeme: string(l.input[l.position:l.readPosition])})
	l.position = l.readPosition
}

func (l *Lexer) addId() {
	lexeme := string(l.input[l.position:l.readPosition])
	switch lexeme {
	case "fn":
		l.addToken("fn")
	case "if":
		l.addToken("if")
	case "then":
		l.addToken("then")
	case "else":
		l.addToken("else")
	case "let":
		l.addToken("let")
	case "do":
		l.addToken("do")
	case "end":
		l.addToken("end")
	default:
		l.addToken("id")
	}
}

func (l *Lexer) skip() {
	l.position = l.readPosition
}

// New creates and returns a pointer to the Lexer
func New(input string) *Lexer {
	l := &Lexer{input: []rune(input), position: 0, readPosition: 0}
	return l
}

func (l *Lexer) readChar() {
	l.char = l.input[l.readPosition]
	l.readPosition += 1
}

func Lex(input string) []Token {
	lexer := New(input)
	lexer.lexTokens()
	return lexer.tokens
}

func (l *Lexer) lexTokens() {
	for l.position < len(l.input) {
		l.lexToken()
	}
}

func (l *Lexer) lexToken() {
	l.readChar()
	switch l.char {
	case '\n':
		l.skip()
	case '\t':
		l.skip()
	case ' ':
		l.skip()
	case '(':
		l.addToken("LPAR")
	case ')':
		l.addToken("RPAR")
	case '+':
		l.addToken("+")
	case '-':
		l.addToken("-")
	case '*':
		l.addToken("*")
	case '/':
		l.addToken("/")
	case ',':
		l.addToken(",")
	case '{':
		l.addToken("{")
	case '}':
		l.addToken("}")
	case '.':
		l.addToken(".")
	case '[':
		l.addToken("[")
	case ']':
		l.addToken("]")
	case ':':
		l.addToken(":")
	case '<':
		if l.peek() == '=' {
			l.readChar()
			l.addToken("<=")
		} else {
			l.addToken("<")
		}
	case '>':
		if l.peek() == '=' {
			l.readChar()
			l.addToken(">=")
		} else {
			l.addToken(">")
		}
	case '=':
		if l.peek() == '>' {
			l.readChar()
			l.addToken("=>")
		} else if l.peek() == '=' {
			l.readChar()
			l.addToken("==")
		} else {
			l.addToken("=")
		}
	default:
		if unicode.IsLetter(l.char) {
			l.lexId()
		} else if unicode.IsNumber(l.char) {
			l.lexNumber()
		} else if l.char == '"' {
			l.skip()
			l.lexString()
		}
	}
}

func (l *Lexer) peek() rune {
	return l.input[l.readPosition]
}

func (l *Lexer) canRead() bool {
	return l.readPosition < len(l.input)
}

func (l *Lexer) lexId() {
	if l.canRead() && unicode.IsLetter(l.peek()) {
		l.readChar()
	}
	for l.canRead() && (unicode.IsLetter(l.peek()) || unicode.IsNumber(l.peek()) || l.peek() == '_') {
		l.readChar()
	}
	l.addId()
}

func (l *Lexer) lexNumber() {
	for l.canRead() && unicode.IsNumber(l.peek()) {
		l.readChar()
	}
	l.addToken("num")
}

func (l *Lexer) lexString() {
	for l.canRead() && l.peek() != '"' {
		l.readChar()
	}
	l.addToken("string")
	l.position += 1
	l.readPosition += 1 // skip the ending ""
}
