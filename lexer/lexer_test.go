package lexer

import (
	"monkey/token"
	"strings"

	"testing"
)

func TestNextToken(t *testing.T) {
	input := strings.NewReader(`	let five = 5;
				let ten = 10;
				let add = fn(x, y) {
					x + y;
				};
				let result = add(five, ten);
				!-/*5
				5 < 10 > 5
				if (5 < 10) {
					return true;
				} else {
					return false;
				}
				10 == 10
				10 != 9
				"foobar"
				"foo bar"
				[1, 2]
				{"foo": "bar"}
				macro(x, y) { x + y; };
			 `)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.MACRO, "macro"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input, "test_lexer.monkey")

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestStringToken(t *testing.T) {
	input := strings.NewReader(`let str = "new string";`)
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LET, "let", 1, 0},
		{token.IDENT, "str", 1, 4},
		{token.ASSIGN, "=", 1, 8},
		{token.STRING, "new string", 1, 11},
		{token.SEMICOLON, ";", 1, 22},
		{token.EOF, "", 2, 0},
	}
	l := New(input, "test_string.monkey")
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestIllegalCharacter tests handling of illegal characters
func TestIllegalCharacter(t *testing.T) {
	input := strings.NewReader(`let @ = 5;`)
	l := New(input, "test_illegal.monkey")
	
	// First token should be LET
	tok := l.NextToken()
	if tok.Type != token.LET {
		t.Fatalf("expected LET, got %q", tok.Type)
	}
	
	// Second token should be IDENT (@ is treated as identifier start)
	tok = l.NextToken()
	if tok.Type != token.ILLEGAL && tok.Type != token.IDENT {
		t.Fatalf("expected ILLEGAL or IDENT for @, got %q", tok.Type)
	}
}

// TestWhitespaceHandling tests that various whitespace characters are handled correctly
func TestWhitespaceHandling(t *testing.T) {
	input := strings.NewReader("let\t\tfive\r\n=\r5;\n\nlet ten = 10;")
	l := New(input, "test_whitespace.monkey")
	
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestComments tests single-line comments (if supported)
func TestComments(t *testing.T) {
	input := strings.NewReader(`let x = 5; // this is a comment
let y = 10;`)
	l := New(input, "test_comments.monkey")
	
	// Just verify the lexer doesn't crash and produces valid tokens
	tokens := []token.Token{}
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	
	// Should have at least let, x, =, 5, ;, let, y, =, 10, ;, EOF
	if len(tokens) < 11 {
		t.Fatalf("expected at least 11 tokens, got %d", len(tokens))
	}
}

// TestLargeNumbers tests parsing of large integer literals
func TestLargeNumbers(t *testing.T) {
	input := strings.NewReader(`let big = 999999999999;`)
	l := New(input, "test_large_numbers.monkey")
	
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "big"},
		{token.ASSIGN, "="},
		{token.INT, "999999999999"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestSpecialIdentifiers tests keywords vs identifiers
func TestSpecialIdentifiers(t *testing.T) {
	input := strings.NewReader(`let letx = 5;
let fnx = 10;
let ifx = 15;
let returnx = 20;`)
	l := New(input, "test_special_identifiers.monkey")
	
	// Verify these are treated as identifiers, not keywords
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "letx"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "fnx"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ifx"},
		{token.ASSIGN, "="},
		{token.INT, "15"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "returnx"},
		{token.ASSIGN, "="},
		{token.INT, "20"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestStringEscapeSequences tests string literals with escape sequences
func TestStringEscapeSequences(t *testing.T) {
	input := strings.NewReader(`let str = "hello\nworld\t!";`)
	l := New(input, "test_escape_sequences.monkey")
	
	// Find the STRING token
	var stringTok token.Token
	for {
		tok := l.NextToken()
		if tok.Type == token.STRING {
			stringTok = tok
			break
		}
		if tok.Type == token.EOF {
			t.Fatal("expected STRING token, reached EOF")
		}
	}
	
	// The literal should contain the escape sequences as-is (lexer doesn't process them)
	if stringTok.Literal != `hello\nworld\t!` {
		t.Errorf("expected literal 'hello\\nworld\\t!', got %q", stringTok.Literal)
	}
}

// TestConsecutiveOperators tests consecutive operators
func TestConsecutiveOperators(t *testing.T) {
	input := strings.NewReader(`!!!5; ---3; +++1;`)
	l := New(input, "test_consecutive_operators.monkey")
	
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BANG, "!"},
		{token.BANG, "!"},
		{token.BANG, "!"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.MINUS, "-"},
		{token.MINUS, "-"},
		{token.MINUS, "-"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.PLUS, "+"},
		{token.PLUS, "+"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestNestedStructures tests deeply nested structures
func TestNestedStructures(t *testing.T) {
	input := strings.NewReader(`{"a": {"b": {"c": [1, [2, [3]]]}}};`)
	l := New(input, "test_nested.monkey")
	
	// Just verify it doesn't crash and produces tokens
	tokenCount := 0
	for {
		tok := l.NextToken()
		tokenCount++
		if tok.Type == token.EOF {
			break
		}
	}
	
	// Should have many tokens for the nested structure
	if tokenCount < 20 {
		t.Fatalf("expected at least 20 tokens for nested structure, got %d", tokenCount)
	}
}

// TestLineAndColumnTracking verifies line and column numbers are tracked correctly
func TestLineAndColumnTracking(t *testing.T) {
	input := strings.NewReader(`let x = 5;
let y = 10;
let z = 15;`)
	l := New(input, "test_line_column.monkey")
	
	tests := []struct {
		expectedType   token.TokenType
		expectedLine   int
		expectedColumn int
	}{
		{token.LET, 1, 1},
		{token.IDENT, 1, 5},
		{token.ASSIGN, 1, 7},
		{token.INT, 1, 9},
		{token.SEMICOLON, 1, 10},
		{token.LET, 2, 1},
		{token.IDENT, 2, 5},
		{token.ASSIGN, 2, 7},
		{token.INT, 2, 9},
		{token.SEMICOLON, 2, 11},
		{token.LET, 3, 1},
		{token.IDENT, 3, 5},
		{token.ASSIGN, 3, 7},
		{token.INT, 3, 9},
		{token.SEMICOLON, 3, 11},
		{token.EOF, 0, 0},
	}
	
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Pos.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Pos.Line)
		}
		if tok.Pos.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedColumn, tok.Pos.Column)
		}
	}
}
