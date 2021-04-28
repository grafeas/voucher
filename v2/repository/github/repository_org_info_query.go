package github

// repositoryOrgInfoQuery is the GraphQL query for retrieving GitHub repository and organizational info
type repositoryOrgInfoQuery struct {
	Resource struct {
		Repository struct {
			Owner struct {
				Typename     string `graphql:"__typename"`
				Organization struct {
					ID   string
					Name string
					URL  string
				} `graphql:"... on Organization"`
			}
		} `graphql:"... on Repository"`
	} `graphql:"resource(url: $url)"`
}
