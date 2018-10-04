package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type CustomID struct {
	value int64
}

func (self *CustomID) Value() int64 {
	return self.value
}

func NewCustomID(v int64) *CustomID {
	return &CustomID{value: v}
}

func serialize(value interface{}) interface{} {
	switch value := value.(type) {
	case CustomID:
		return value.Value()
	case *CustomID:
		return *value.Value()
	default:
		return nil
	}
}

func parseValue(value interface{}) interface{} {
	switch value := value.(type) {
	case int64:
		return NewCustomID(value)
	case *int64:
		return NewCustomID(*value)
	default:
		return nil
	}
}

func parseLiteral(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.IntValue:
		return NewCustomID(valueAST.value)
	default:
		return nil
	}
}

var CustomScalaType = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "CustomScalaType",
	Description:  "Custom scala type for `GraphQLIDType`",
	Serialize:    serialize,
	ParseValue:   parseValue,
	ParseLiteral: parseLiteral,
})
