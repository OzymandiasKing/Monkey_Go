package lexer

import (
	"bufio"
	"bytes"
	"io"
	"monkey/token"
)

type Lexer struct {
	input    *bufio.Reader
	ch       byte
	fileName string
	line     int
	column   int
}

func New(input io.Reader, fileName string) *Lexer {
	l := &Lexer{
		input:    bufio.NewReader(input),
		fileName: fileName,
		line:     1,
		column:   0,
	}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	pos := token.Position{
		Filename: l.fileName,
		Line:     l.line,
		Column:   l.column,
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Pos: pos}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '<':
		tok = l.newToken(token.LT, l.ch)
	case '>':
		tok = l.newToken(token.GT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Pos: pos}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Pos = pos
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			tok.Pos = pos
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Pos: token.Position{
			Filename: l.fileName,
			Line:     l.line,
			Column:   l.column,
		},
	}
}

func (l *Lexer) readChar() {
	var err error
	l.ch, err = l.input.ReadByte()
	if err != nil {
		l.ch = 0
	}
	if l.ch == '\n' {
		l.line += 1
		l.column = 0
	} else {
		l.column += 1
	}
}

func (l *Lexer) peekChar() byte {
	bytes, err := l.input.Peek(1)
	if err != nil {
		return 0
	}
	return bytes[0]
}

// readIdentifier 读取标识符（变量名、函数名等）
// 修改意图：修复字符读取错位问题
// 原实现使用 Peek(1) 检查下一个字符后调用 readChar()，导致循环结束时 l.ch 指向标识符的最后一个字符
// 新实现先通过 peekChar() 预读判断，再调用 readChar() 读取，最后额外调用一次 readChar() 确保 l.ch 指向下一个待处理字符
// 这样保证词法分析器在返回标识符后能正确定位到下一个字符位置
func (l *Lexer) readIdentifier() string {
	var buf bytes.Buffer
	buf.WriteByte(l.ch)
	for isLetter(l.peekChar()) {
		l.readChar()
		buf.WriteByte(l.ch)
	}
	l.readChar() // 移动到下一个字符（关键修复：确保循环结束后 l.ch 指向标识符之后的字符）
	return buf.String()
}

// readNumber 读取数字字面量
// 修改意图：修复字符读取错位问题（同 readIdentifier）
// 原实现在循环中使用 Peek(1) 检查后读取，导致函数返回时 l.ch 停留在数字的最后一位
// 这会导致 NextToken() 末尾的 readChar() 跳过下一个有效字符
// 新实现确保读完数字后，l.ch 正确指向数字之后的下一个待处理字符
func (l *Lexer) readNumber() string {
	var buf bytes.Buffer
	buf.WriteByte(l.ch)
	for isDigit(l.peekChar()) {
		l.readChar()
		buf.WriteByte(l.ch)
	}
	l.readChar() // 移动到下一个字符（关键修复：确保循环结束后 l.ch 指向数字之后的字符）
	return buf.String()
}

func (l *Lexer) readString() string {
	var buf bytes.Buffer
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		buf.WriteByte(l.ch)
	}
	return buf.String()
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
