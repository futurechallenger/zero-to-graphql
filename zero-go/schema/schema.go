package schema

import (
	"zero-go/model"

	"log"

	"github.com/graphql-go/graphql"
)

// PersonType is graphql person type
var PersonType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PersonType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: CustomScalaType,
			},
			"lastLogin": &graphql.Field{
				Type: graphql.DateTime,
			},
			"isSuperUser": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if person, ok := params.Source.(model.Person); ok == true {
						return person.IsSuperuser, nil
					}
					return nil, nil
				},
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"fullName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					firstName, fnOK := params.Args["firstName"].(string)
					lastName, lnOK := params.Args["lastName"].(string)
					if fnOK && lnOK {
						return firstName + " " + lastName, nil
					}
					return nil, nil
				},
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			// "friends": &graphql.Field{
			// 	Type:        graphql.NewList(PersonType),
			// 	Description: "People who hang out with you",
			// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 		person, ok := params.Source.(model.Person)
			// 		return nil, nil
			// 	},
			// },
		},
	},
)

// CreateSchema will create graphql schema
func CreateSchema() (*graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"allPeople": &graphql.Field{
				Type: graphql.NewList(PersonType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					ch := make(chan string, 1)
					go func() {
						defer close(ch)
						ret := GetAllPeople()
						ch <- ret
					}()

					return func() (interface{}, error) {
						r := <-ch
						return r, nil
					}, nil
				},
			},
			"friends": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Type: graphql.NewList(PersonType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					personID, ok := params.Args["id"].(int64)
					if !ok {
						log.Fatalln("Can not get person ID")
						return nil, nil
					}

					ch := make(chan string, 1)
					go func() {
						defer close(ch)
						ret := GetFriends()
						ch <- ret
					}()

					return func() (interface{}, error) {
						r := <-ch
						return r, nil
					}, nil
				},
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
