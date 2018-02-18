package json

import . "github.com/tgascoigne/pogo"

var value, intLit, stringLit, booleanLit, list, dict, field Parser

func init() {
	intLit = TypedProd("IntLiteral", int(0),
		Tok(INT_LITERAL))

	stringLit = TypedProd("StringLiteral", string(""),
		Tok(STRING_LITERAL))

	booleanLit = TypedProd("BooleanLiteral", bool(false),
		Tok(BOOL_LITERAL))

	list = TypedProd("List", []interface{}{},
		Seq(Char('['),
			SepByTerm(&value, Char(','),
				Char(']'))))

	dict = TypedProd("Dict", map[interface{}]interface{}{},
		Seq(Char('{'),
			SepByTerm(&field, Char(','),
				Char('}'))))

	field = TypedProd("Field", dictField{},
		Seq(&stringLit, Char(':'), &value))

	value = TypedProd("Value", jsonValue{},
		Or(&intLit, &stringLit, &booleanLit, &list, &dict))
}
