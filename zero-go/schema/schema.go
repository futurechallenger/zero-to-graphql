package schema

import (
	"errors"
	"fmt"
	"sync"
	"zero-go/model"
	"zero-go/util"

	"github.com/graphql-go/graphql"
	dataloader "gopkg.in/nicksrandall/dataloader.v5"
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
					person, ok := params.Source.(*model.Person)
					if !ok {
						return "", errors.New("No person found for resolving `fullName`")
					}
					return person.FirstName + " " + person.LastName, nil
				},
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func init() {
	PersonType.AddFieldConfig("friends", &graphql.Field{
		Type:        graphql.NewList(PersonType),
		Description: "People who hang out with you",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// person, ok := params.Source.(model.Person)
			// var resolveRet ResolveRet
			// var friends []model.Person

			// if !ok {
			// 	log.Fatalln("Can not get person ID")
			// 	return nil, nil
			// }

			// ch := make(chan *ResolveRet)
			// go func() {
			// 	defer close(ch)
			// 	ret := GetFriends(person.ID)

			// 	json.Unmarshal([]byte(ret), &friends)

			// 	resolveRet.data = friends
			// 	resolveRet.err = nil
			// 	ch <- &resolveRet
			// }()

			// return func() interface{} {
			// 	r := <-ch
			// 	return r.data
			// }, nil

			var (
				person, _ = params.Source.(model.Person)
				v         = params.Context.Value
				loaders   = v(util.LoadersKey).(map[string]*dataloader.Loader)
				c         = v(util.ClientKey).(*Client)
				thunks    []dataloader.Thunk
				wg        sync.WaitGroup
			)

			for _, personID := range person.Friends {
				key := NewResolverKey(fmt.Sprintf("%d", personID), c)
				thunk := loaders["personLoader"].Load(params.Context, key)
				thunks = append(thunks, thunk)
			}

			type result struct {
				friends []*model.Person
				errs    []error
			}
			// ch := make(chan *result, 1)
			ch := make(chan *result)

			go func() {
				var friends []*model.Person
				var errs []error
				for _, thunk := range thunks {
					wg.Add(1)
					go func(t dataloader.Thunk) {
						defer wg.Done()

						r, err := t()
						if err != nil {
							errs = append(errs, err)
							return
						}

						p := r.(*model.Person)
						friends = append(friends, p)
					}(thunk)
				}

				wg.Wait()
				ch <- &result{friends: friends, errs: errs}
			}()

			return func() interface{} {
				ret := <-ch
				errs := ret.errs
				friends := ret.friends
				if len(errs) > 0 {
					util.HandleError(errs[len(errs)-1])
					return errs
				}

				return friends
			}, nil
		},
	})
}

// CreateSchema will create graphql schema
func CreateSchema() (*graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"allPeople": &graphql.Field{
				Type: graphql.NewList(PersonType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// var result ResolveRet
					// ch := make(chan *ResolveRet)
					// go func() {
					// 	defer close(ch)
					// 	ret := GetAllPeople()

					// 	fmt.Printf("resolve all prople %v", ret)

					// 	// Deserialize
					// 	var personList []model.Person
					// 	err := json.Unmarshal([]byte(ret), &personList)
					// 	if err != nil {
					// 		result.data = nil
					// 		result.err = err
					// 	} else {
					// 		result.data = personList
					// 		result.err = nil
					// 	}
					// 	ch <- &result
					// }()

					// return func() interface{} {
					// 	r := <-ch
					// 	return r.data
					// }, nil

					var (
						v       = params.Context.Value
						c       = v(util.ClientKey).(*Client)
						loaders = v(util.LoadersKey).(map[string]*dataloader.Loader)
						key     = NewResolverKey("__all__", c)
					)

					fmt.Printf("resolve p.Context.Value %+v\n", v("client"))

					thunk := loaders["allPeopleLoader"].Load(params.Context, key)

					return func() interface{} {
						ret, _ := thunk()
						return ret
					}, nil
				},
			},
			// "friends": &graphql.Field{
			// 	Args: graphql.FieldConfigArgument{
			// 		"id": &graphql.ArgumentConfig{
			// 			Type: graphql.Int,
			// 		},
			// 	},
			// 	Type: graphql.NewList(PersonType),
			// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 		personID, ok := params.Args["id"].(int64)
			// 		if !ok {
			// 			log.Fatalln("Can not get person ID")
			// 			return nil, nil
			// 		}

			// 		ch := make(chan string, 1)
			// 		go func() {
			// 			defer close(ch)
			// 			ret := GetFriends(personID)
			// 			ch <- ret
			// 		}()

			// 		return func() interface{} {
			// 			r := <-ch
			// 			return r
			// 		}, nil
			// 	},
			// },
			"person": &graphql.Field{
				Type: PersonType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					ret := &model.Person{
						FirstName: "Bruce",
						LastName:  "Lee",
					}
					return func() interface{} {
						return ret
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
