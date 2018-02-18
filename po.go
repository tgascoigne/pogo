package pogo

import (
	"fmt"
	"reflect"

	"gitea.local/tom/calx/errors"
)

type Parser func(p ParseState) (Parsed, ParseState)

var ErrExpected = errors.New("expected '%v' got '%v'")

// NilParsed is returned when nothing could be parsed, but no error was raised
var NilParsed = NilParsedType{}

// ErrParsed is returned when nothing could be parsed, and an error was raised
var ErrParsed = ErrParsedType{}

var EOF = Tok(TOK_EOF)

var Debug = false

// Tok is a parser which matches a given TokenType. The token is only consumed if it matches.
func Tok(token TokenType) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		item := peekTok(ps)
		if item.Token == token {
			return nextTok(ps)
		}

		ps = addError(ps, item, ErrExpected.Format(token, item))
		return ErrParsed, ps
	}
}

// Char is a shortcut for parsing single character tokens
func Char(c rune) Parser {
	return Tok(TokenType(c))
}

// Parses a named production
func Prod(ident string, p interface{}) Parser {
	VisitorTemplate.ProductionOrder = append(VisitorTemplate.ProductionOrder, ident)
	VisitorTemplate.Productions[ident] = nil

	return func(ps ParseState) (Parsed, ParseState) {
		if Debug {
			fmt.Printf("parsing prod %v at %v: %v\n", ident, ps.Pos(), ps.Lexer().At(ps))
		}

		var children Sequence
		var res Parsed
		res, ps = Do(p, ps)
		switch res := res.(type) {
		case Item:
			children = Sequence{res}

		case Production:
			children = Sequence{res}

		case Sequence:
			children = res.Flatten()

		case NilParsedType:

		case ErrParsedType:
			if Debug {
				fmt.Printf("prod %v failed at %v: %v\n", ident, ps.Pos(), ps.Lexer().At(ps))
			}
			return res, ps

		default:
			panic(fmt.Sprintf("don't know what to do with %T\n", res))
		}

		if Debug {
			fmt.Printf("prod %v at %v: %v -> %v\n", ident, ps.Pos(), ps.Lexer().At(ps), res)
		}

		return Production{
			Ident:    ident,
			Children: children,
		}, ps
	}
}

func Named(ident string, p interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		if Debug {
			fmt.Printf("parsing named %v at %v: %v\n", ident, ps.Pos(), ps.Lexer().At(ps))
		}

		var value Parsed
		value, ps = Do(p, ps)
		if value == ErrParsed {
			return value, ps
		}

		return NamedValue{
			Name:  ident,
			Value: value,
		}, ps
	}
}

// Parses a named and typed production
func TypedProd(ident string, typ interface{}, p interface{}) Parser {
	parser := Prod(ident, p)
	VisitorTemplate.Productions[ident] = reflect.TypeOf(typ)
	return parser
}

// Parses a named and typed production. Special case for interface typed productions.
// Usage: TypedProdIface("prod", (*MyIface)(nil))
func TypedProdIface(ident string, typ interface{}, p interface{}) Parser {
	parser := Prod(ident, p)
	VisitorTemplate.Productions[ident] = reflect.TypeOf(typ).Elem()
	return parser
}

// Sequence accepts 1..n parsers and expects them to succeed in their defined sequence
func Seq(p ...interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		results := make(Sequence, len(p))
		for i, p := range p {
			results[i], ps = Do(p, ps)
			// If any in the sequence fail, return error
			if results[i] == ErrParsed {
				return ErrParsed, ps
			}
		}

		return results, ps
	}
}

// Maybe executes a parser which may fail. If it fails no input is consumed and NilParsed is returned.
// Errors are not sent to the error handler.
func Maybe(p interface{}) Parser {
	return TryRecover(p, func(start, end ParseState) (Parsed, ParseState) {
		return NilParsed, start
	})
}

// Peek parses p without consuming any input.
// Errors are sent to the error handler.
func Peek(p interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		res, _ := Do(p, ps)
		return res, ps
	}
}

// MaybePeek parses p without consuming any input. If it fails, NilParsed is returned.
// Errors are not sent to the error handler.
func MaybePeek(p interface{}) Parser {
	return Peek(Maybe(p))
}

func Many(p interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		results := make(Sequence, 0)

		for {
			var res Parsed
			res, ps = Do(Maybe(p), ps)
			// Maybe() will return NilParsed when the parser fails; exit cleanly
			if res == NilParsed {
				break
			}

			results = append(results, res)
		}

		return results, ps
	}
}

func Many1(p interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		var results Sequence
		var value Parsed
		value, ps = Do(p, ps)
		results = Sequence{value}

		value, ps = Do(Many(p), ps)
		results = append(results, value.(Sequence)...)
		return results, ps
	}
}

func SepBy(p, delim interface{}) Parser {
	return doSepBy(p, delim, nil, false)
}

func SepByTerm(p, delim, terminator interface{}) Parser {
	return doSepBy(p, delim, terminator, false)
}

func SepBy1(p, delim interface{}) Parser {
	return doSepBy(p, delim, nil, true)
}

func SepBy1Term(p, delim, terminator interface{}) Parser {
	return doSepBy(p, delim, terminator, true)
}

func doSepBy(p, delim, terminator interface{}, requireOne bool) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		results := make(Sequence, 0)

		required := false
		if requireOne {
			required = true
		}

		var res Parsed
		for {
			// Parse the subject (must succeed unless required == true)
			// FIXME pretty sure this is backwards?
			if !required {
				res, ps = Do(Maybe(p), ps)
				if res == NilParsed {
					break
				}
			} else {
				res, ps = Do(p, ps)
			}

			if res == ErrParsed {
				return res, ps
			}

			results = append(results, res)
			required = false

			// Parse the delimiter (may fail the parser gracefully)
			res, ps = Do(Maybe(delim), ps)
			if res == NilParsed {
				break
			}
		}

		if terminator != nil {
			res, ps = Do(terminator, ps)
			if res == ErrParsed {
				return res, ps
			}
		}

		return results, ps
	}
}

// Or accepts 1..n parsers and attempts them one by one. The first successful match is returned.
// If no parsers succeed, then the error of the most successful (by length) parser is raised.
func Or(p ...interface{}) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		longestParse := -1
		var longestParseState ParseState
		var res Parsed

		for _, p := range p {
			res, ps = Do(TryRecover(p, func(start, end ParseState) (Parsed, ParseState) {
				thisParseLength := end.Pos() - start.Pos()
				if thisParseLength > longestParse {
					longestParseState = end
					longestParse = thisParseLength
				}

				// Backtrack to the point that we started at
				return ErrParsed, start
			}), ps)

			if res != ErrParsed {
				return res, ps
			}
		}

		return ErrParsed, longestParseState
	}
}

func Do(p interface{}, ps ParseState) (Parsed, ParseState) {
	switch p := p.(type) {
	case Parser:
		return p(ps)
	case *Parser:
		return (*p)(ps)
	default:
		panic(fmt.Sprintf("Not a parser: %T", p))
	}
}
