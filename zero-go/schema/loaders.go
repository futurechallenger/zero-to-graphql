package schema

import (
	"context"
	"strconv"

	dataloader "gopkg.in/nicksrandall/dataloader.v5"
)

func peopleBacheFunc(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		reslts = append(results, &dataloader.Result{Error: err})
		return results
	}

	var personIDs []int64
	for _, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err != nil {
			handleError(err)
		}
		personIDs = append(personIDs, id)
	}
}

// PeopleLoader used to batch people requests
var PeopleLoader = dataloader.NewBatchedLoader(peopleBacheFunc, dataloader.WithCache(LoaderCache))
