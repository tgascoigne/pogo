package pogo

import (
	"fmt"
	"strings"
)

type Parsed interface {
	ptree()
	String() string
}

type NilParsedType struct{}

func (NilParsedType) ptree()         {}
func (NilParsedType) String() string { return "Nil" }

type ErrParsedType struct{}

func (ErrParsedType) ptree()         {}
func (ErrParsedType) String() string { return "Err" }

type NamedValue struct {
	Name  string
	Value Parsed
}

func (n NamedValue) ptree()         {}
func (n NamedValue) String() string { return fmt.Sprintf("(\"%v\" %v)", n.Name, n.Value) }

type Production struct {
	Ident    string
	Children Sequence
}

func (t Production) ptree() {}

func (t Production) String() string {
	return fmt.Sprintf("(%v %v)", t.Ident, t.Children.String())
}

type Item struct {
	Token      TokenType
	Value      string
	Offset     Pos
	Line       int
	Column     int
	Index      int
	SourceFile string
}

var NilItem = Item{}

func (i Item) ptree() {}

func (i Item) String() string {
	// Special case for pogo.Char style tokens
	if len(string(i.Token)) == 1 {
		return fmt.Sprintf("'%v'", i.Token)
	}

	return fmt.Sprintf("[%v %v]", i.Token, i.Value)
}

func (i Item) Start() Pos {
	return i.Offset
}

func (i Item) End() Pos {
	return i.Offset + Pos(len(i.Value))
}

type Sequence []Parsed

func (s Sequence) ptree() {}

func (s Sequence) Flatten() []Parsed {
	result := make([]Parsed, 0)
	for _, p := range s {
		switch p := p.(type) {
		case Item:
			result = append(result, p)

		case Production:
			result = append(result, p)

		case NamedValue:
			result = append(result, p)

		case Sequence:
			p = p.Flatten()
			result = append(result, p...)

		default:
			if p == NilParsed {
				result = append(result, NilParsed)
				break
			}

			panic(fmt.Sprintf("unknown parse tree node %T", p))
		}
	}
	return result
}

func (s Sequence) Filter(tokens ...TokenType) Sequence {
	seq := make(Sequence, 0)
OUTER:
	for _, p := range s {
		if item, ok := p.(Item); ok {
			for _, tok := range tokens {
				if item.Token == tok {
					continue OUTER
				}
			}
		}

		seq = append(seq, p)
	}

	return seq
}

func (s Sequence) String() string {
	children := make([]string, len(s))
	for i, c := range s {
		children[i] = c.String()
	}
	return strings.Join(children, " ")
}
