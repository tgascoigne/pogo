package json

import (
	"io"

	"github.com/tgascoigne/pogo"
)

// nex stuff
type yySymType struct {
	tok pogo.Item
}

func (l *Lexer) Error(err string) {
	panic(err)
}

// Bridges nex -> pogo.Lexer
type LexerImpl struct {
	buffer []pogo.Item
}

func NewBufferedLexer(r io.Reader) pogo.LexerIface {
	nex := NewLexer(r)
	lexer := &LexerImpl{
		buffer: []pogo.Item{},
	}
	lexer.bufferItems(nex)
	return lexer
}

func (b *LexerImpl) bufferItems(nex *Lexer) {
	for {
		next := yySymType{}
		ret := nex.Lex(&next)
		b.buffer = append(b.buffer, next.tok)
		if ret == 0 {
			break
		}
	}
}

func (b *LexerImpl) At(ps pogo.ParseState) pogo.Item {
	if ps.Pos() >= len(b.buffer) {
		return pogo.Item{
			Token: pogo.TOK_EOF,
			Value: "EOF",
		}
	}
	return b.buffer[ps.Pos()]
}
