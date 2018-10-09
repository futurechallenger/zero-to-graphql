package schema

import (
	"context"
	"encoding/json"
	"strconv"

	"zero-go/model"

	dataloader "gopkg.in/nicksrandall/dataloader.v5"
)

// Client reprensents http requests
type Client struct {
}

// GetAllPeople request all people
func (c *Client) GetAllPeople() ([]model.Person, error) {
	ret := GetAllPeople()

	// Deserialize
	var personList []model.Person
	err := json.Unmarshal([]byte(ret), &personList)

	return personList, err
}

// GetPerson will request just a person
func (c *Client) GetPerson(personID int64) (model.Person, error) {
	ret := GetPerson(personID)

	var person model.Person
	err := json.Unmarshal([]byte(ret), &person)

	return person, err
}

// GetFriends will request a person's all friends
func (c *Client) GetFriends(personID int64) ([]model.Person, error) {
	ret := GetFriends(personID)

	// Deserialize
	var personList []model.Person
	err := json.Unmarshal([]byte(ret), &personList)

	return personList, err
}

// ResolverKey is a key for batched function
type ResolverKey struct {
	Key    string
	Client *Client
}

// NewResolverKey method create a `ResolverKey` instance
func NewResolverKey(key string, client *Client) *ResolverKey {
	return &ResolverKey{Key: key, Client: client}
}

func (rk *ResolverKey) client() *Client {
	return rk.Client
}

// String return resolver key
func (rk *ResolverKey) String() string {
	return rk.Key
}

// Raw return key's raw value, here it still string but typed interface{}
func (rk *ResolverKey) Raw() interface{} {
	return rk.Key
}

func allPeopleBatchedFunc(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	personList, err := keys[0].(*ResolverKey).client().GetAllPeople()
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	for _, p := range personList {
		result := &dataloader.Result{
			Data:  p,
			Error: nil,
		}
		results = append(results, result)
	}

	return results
}

func friendBatchedFunc(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	personID, err := strconv.ParseInt(keys[0].(*ResolverKey).String(), 10, 64)
	if err != nil {
		return handleError(err)
	}

	personList, err := keys[0].(*ResolverKey).client().GetFriends(personID)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	for _, p := range personList {
		result := &dataloader.Result{
			Data:  p,
			Error: nil,
		}
		results = append(results, result)
	}

	return results
}

func personBatchedFunc(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	personID, err := strconv.ParseInt(keys[0].(*ResolverKey).String(), 10, 64)
	if err != nil {
		return handleError(err)
	}

	person, err := keys[0].(*ResolverKey).client().GetPerson(personID)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := &dataloader.Result{
		Data:  person,
		Error: nil,
	}
	results = append(results, result)

	return results
}

// AllPeopleLoader used to batch people requests
var AllPeopleLoader = dataloader.NewBatchedLoader(allPeopleBatchedFunc, dataloader.WithCache(LoaderCache))

// FriendsLoader loads all friends of a person
var FriendsLoader = dataloader.NewBatchedLoader(friendBatchedFunc, dataloader.WithCache(LoaderCache))

// PersonLoader loads a person
var PersonLoader = dataloader.NewBatchedLoader(personBatchedFunc, dataloader.WithCache(LoaderCache))
