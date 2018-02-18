package json

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/tgascoigne/pogo"
)

//go:generate nex -e json.nex
//go:generate go run ./cmd/json-gen/main.go -generate-visitor -path . -pkg json

type ParserImpl struct {
	pogo.BaseParseState
	lexer pogo.LexerIface
}

func newParser(r io.Reader) *ParserImpl {
	return &ParserImpl{
		lexer: NewBufferedLexer(r),
	}
}

func (p *ParserImpl) HandleError(tok pogo.Item, err error) {
	fmt.Printf("error reported at %v: %v\n", tok, err)
}

func (p *ParserImpl) Clone() pogo.ParseState {
	return &ParserImpl{
		BaseParseState: p.BaseParseState.CloneBase(),
		lexer:          p.lexer,
	}
}

func (p *ParserImpl) Lexer() pogo.LexerIface {
	return p.lexer
}

func Parse(r io.Reader) {
	parser := newParser(r)
	tree, state := pogo.Do(pogo.Root(&value), parser)
	for _, err := range state.Errors() {
		parser.HandleError(err.Token, err.Error)
	}
	ast := acceptValue(new(astBuilder), tree)
	fmt.Printf("parse tree %v\n", tree)
	spew.Dump(ast)
}
