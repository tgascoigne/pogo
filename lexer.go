package pogo

import "go/token"

type Pos token.Pos

type TokenType string

const TOK_EOF = TokenType("")

type LexerIface interface {
	At(state ParseState) Item
}
