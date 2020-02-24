package github

import (
	"context"
)

// queryHandler is called on every iteration of paginationQuery to populate a slice of query results
// queryHandler checks to see whether there are more records given that GitHub has a limit of 100 records per query
type queryHandler func(queryResult interface{}) (bool, error)

// paginationQuery populates a destination slice with the appropriately typed query results
// GitHub has a limit of 100 records so we must perform pagination
func paginationQuery(
	ctx context.Context,
	ghc ghGraphQLClient,
	queryResult interface{},
	queryPopulationVariables map[string]interface{},
	pageLimit int,
	qh queryHandler,
) error {
	for i := 0; i < pageLimit; i++ {
		err := ghc.Query(ctx, queryResult, queryPopulationVariables)
		if err != nil {
			return err
		}

		hasMoreResults, err := qh(queryResult)
		if nil != err {
			return err
		}

		if !hasMoreResults {
			return nil
		}
	}
	return nil
}

// commit contains information pertaining to a commit
type commit struct {
	URL string
}
