package github

import (
	"context"
	"fmt"
	"testing"

	utils "github.com/mennanov/fieldmask-utils"
	"github.com/stretchr/testify/assert"
)

type queryHandlerFunc func(query interface{}, variables map[string]interface{}) error

func createHandler(input interface{}, mask []string) queryHandlerFunc {
	return func(query interface{}, variables map[string]interface{}) error {
		mask, err := utils.MaskFromPaths(mask, func(s string) string { return s })
		if err != nil {
			return err
		}

		err = utils.StructToStruct(mask, input, query)
		if err != nil {
			return err
		}

		return nil
	}
}

type mockGitHubGraphQLClient struct {
	HandlerFunc queryHandlerFunc
}

func (m *mockGitHubGraphQLClient) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	return m.HandlerFunc(query, variables)
}

func TestPaginationQuery(t *testing.T) {
	paginationTests := []struct {
		testName                 string
		queryResult              interface{}
		queryPopulationVariables map[string]interface{}
		pageLimit                int
		resultHandler            queryHandler
		errorExpected            error
	}{
		{
			testName:                 "Test infinite result with finite page limit",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return true, nil
			},
			errorExpected: nil,
		},
		{
			testName:                 "Test error result",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return true, fmt.Errorf("cannot paginate anymore")
			},
			errorExpected: fmt.Errorf("cannot paginate anymore"),
		},
		{
			testName:                 "Test false result",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return false, nil
			},
			errorExpected: nil,
		},
	}

	for _, test := range paginationTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = func(query interface{}, variables map[string]interface{}) error {
				return nil
			}

			err := paginationQuery(
				context.Background(),
				c,
				test.queryResult,
				test.queryPopulationVariables,
				test.pageLimit,
				test.resultHandler,
			)

			if test.errorExpected != nil {
				assert.Equal(t, test.errorExpected, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
