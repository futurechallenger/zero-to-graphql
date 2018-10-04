package schema

import (
	"github.com/graphql-go/graphql"
)

// CreateSchema will create graphql schema
func CreateSchema() (*graphql.Schema, error) {
	PersonType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "PersonType",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: CustomScalaType,
				},
				"password": &graphql.Field{
					Type: graphql.String,
				},
				"lastLogin": &graphql.Field{
					Type: graphql.DateTime,
				},
				"isSuperUser": &graphql.Field{
					Type: graphql.Boolean,
				},
				"firstName": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"allPeople": &graphql.Field{
				Type: graphql.NewList(PersonType),
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

	if err != nil {
		return nil, err
	}

	return &schema, nil
}
