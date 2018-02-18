package json

import (
	"strconv"

	"github.com/tgascoigne/pogo"
)

// astBuilder implements jsonVisitor in order to walk the parse tree
// and generate a map from the json structure.
type astBuilder struct {
}

func (v *astBuilder) visitIntLiteral(delegate jsonVisitor, items []pogo.Parsed) int {
	i, err := strconv.Atoi(items[0].(pogo.Item).Value)
	if err != nil {
		panic(err)
	}

	return i
}

func (v *astBuilder) visitStringLiteral(delegate jsonVisitor, items []pogo.Parsed) string {
	str := items[0].(pogo.Item).Value
	return str[1 : len(str)-1]
}

func (v *astBuilder) visitBooleanLiteral(delegate jsonVisitor, items []pogo.Parsed) bool {
	b, err := strconv.ParseBool(items[0].(pogo.Item).Value)
	if err != nil {
		panic(err)
	}

	return b
}

func (v *astBuilder) visitList(delegate jsonVisitor, items []pogo.Parsed) []interface{} {
	values := allValues(delegate, items)

	result := make([]interface{}, len(values))
	for i, value := range values {
		result[i] = value.Value
	}
	return result
}

func (v *astBuilder) visitDict(delegate jsonVisitor, items []pogo.Parsed) map[interface{}]interface{} {
	fields := allFields(delegate, items)

	result := make(map[interface{}]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}

func (v *astBuilder) visitField(delegate jsonVisitor, items []pogo.Parsed) dictField {
	key := acceptStringLiteral(delegate, items[0])
	value := acceptValue(delegate, items[2])
	return dictField{key, value.Value}
}

func (v *astBuilder) visitValue(delegate jsonVisitor, items []pogo.Parsed) jsonValue {
	return jsonValue{accept(delegate, items[0])}
}
