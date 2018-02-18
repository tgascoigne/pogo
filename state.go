package pogo

type ParseState interface {
	Lexer() LexerIface
	Pos() int
	Errors() []ParseError
	Clone() ParseState

	base() *BaseParseState
}

type ParseError struct {
	Token Item
	Error error
}

type BaseParseState struct {
	pos    int
	errors []ParseError
}

func (b *BaseParseState) branch() string {
	return ""
}

func (b *BaseParseState) base() *BaseParseState {
	return b
}

func (b *BaseParseState) Errors() []ParseError {
	return b.errors
}

func (b *BaseParseState) Pos() int {
	return b.pos
}

func (b BaseParseState) CloneBase() BaseParseState {
	new := BaseParseState{
		pos:    b.pos,
		errors: make([]ParseError, len(b.errors)),
	}
	copy(new.errors, b.errors)
	return new
}

func nextTok(ps ParseState) (Item, ParseState) {
	newState := ps.Clone()
	newState.base().pos++
	return ps.Lexer().At(ps), newState
}

func peekTok(ps ParseState) Item {
	return ps.Lexer().At(ps)
}

func lastTok(ps ParseState) Item {
	ps = ps.Clone()
	ps.base().pos--
	return ps.Lexer().At(ps)
}

func addError(ps ParseState, tok Item, err error) ParseState {
	ps = ps.Clone()
	base := ps.base()
	base.errors = append(base.errors, ParseError{tok, err})
	return ps
}

func addErrors(source ParseState, dest ParseState) ParseState {
	for _, e := range source.Errors() {
		dest = addError(dest, e.Token, e.Error)
	}
	return dest
}
