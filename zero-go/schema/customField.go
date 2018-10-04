package schema

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// CustomID hold a int64 value in it
type CustomID struct {
	value int64
}

// Value return int64 value
func (c *CustomID) Value() int64 {
	return c.value
}

// NewCustomID return a new `CustomID`
func NewCustomID(v int64) *CustomID {
	return &CustomID{value: v}
}

func serialize(value interface{}) interface{} {
	switch value := value.(type) {
	case CustomID:
		return strconv.FormatInt(value.Value(), 10)
	case *CustomID:
		v := *value
		return strconv.FormatInt(v.Value(), 10)
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
		if int64Value, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
			return NewCustomID(int64Value)
		}
	default:
		return nil
	}
	return nil
}

// CustomScalaType is a representation of `GraphQLIDType`
var CustomScalaType = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "CustomScalaType",
	Description:  "Custom scala type for `GraphQLIDType`",
	Serialize:    serialize,
	ParseValue:   parseValue,
	ParseLiteral: parseLiteral,
})
