package schema

import (
	"context"
	"encoding/json"
	"strconv"

	"zero-go/model"
	"zero-go/util"

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
func (c *Client) GetPerson(personID int64) (*model.Person, error) {
	ret := GetPerson(personID)

	var person []model.Person
	err := json.Unmarshal([]byte(ret), &person)

	if err == nil && len(person) > 0 {
		return &person[0], nil
	}

	return nil, err
}

// GetFriends will request a person's all friends
func (c *Client) GetFriends(personIDs []int64) ([]*model.Person, error) {
	// ret := GetFriends(personID)

	// // Deserialize
	// var personList []model.Person
	// err := json.Unmarshal([]byte(ret), &personList)

	// return personList, err

	util.ULog("Request friend start")
	var personList []*model.Person

	for _, personID := range personIDs {
		p, err := c.GetPerson(personID)
		if err != nil {
			util.HandleError(err)
			continue
		}

		personList = append(personList, p)
	}

	return personList, nil
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
	util.ULog(keys)

	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	personList, err := keys[0].(*ResolverKey).client().GetAllPeople()
	if err != nil {
		util.HandleError(err)
		return handleError(err)
	}

	var results []*dataloader.Result

	result := &dataloader.Result{
		Data:  personList,
		Error: nil,
	}
	results = append(results, result)

	return results
}

func friendBatchedFunc(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	util.ULog("start")
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	// personID, err := strconv.ParseInt(keys[0].(*ResolverKey).String(), 10, 64)
	// if err != nil {
	// 	util.HandleError(err)
	// 	return handleError(err)
	// }

	var personIDs []int64

	for _, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err != nil {
			util.ULog(err)
			return handleError(err)
		}

		personIDs = append(personIDs, id)
	}

	personList, err := keys[0].(*ResolverKey).client().GetFriends(personIDs)
	if err != nil {
		util.HandleError(err)
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
	util.ULog("start", keys)
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		results = append(results, &dataloader.Result{Error: err})
		return results
	}

	personID, err := strconv.ParseInt(keys[0].(*ResolverKey).String(), 10, 64)
	if err != nil {
		util.HandleError(err)
		return handleError(err)
	}

	person, err := keys[0].(*ResolverKey).client().GetPerson(personID)
	if err != nil {
		util.HandleError(err)
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

// BatchedLoaders are loaders map
var BatchedLoaders = map[string]*dataloader.Loader{
	"allPeopleLoader": dataloader.NewBatchedLoader(allPeopleBatchedFunc, dataloader.WithCache(LoaderCache)),
	"friendsLoader":   dataloader.NewBatchedLoader(friendBatchedFunc, dataloader.WithCache(LoaderCache)),
	"personLoader":    dataloader.NewBatchedLoader(personBatchedFunc, dataloader.WithCache(LoaderCache)),
}
