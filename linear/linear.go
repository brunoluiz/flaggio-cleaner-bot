package linear

import (
	"context"
	"net/http"

	"github.com/machinebox/graphql"
)

const issueCreateMutation = `
	mutation ($input: IssueCreateInput!) {
		issueCreate(input: $input) {
			success
		}
	}
`

// Client graphql linear client
type Client struct {
	gql *graphql.Client
}

// New returns a graphql linear instance
func New(token string) *Client {
	httpClient := &http.Client{
		Transport: newAddHeaderTransport(token),
	}
	client := graphql.NewClient("https://api.linear.app/graphql", graphql.WithHTTPClient(httpClient))
	return &Client{client}
}

// IssueCreateInput graphql model for creating an issue
type IssueCreateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TeamID      string `json:"teamId"`
	ProjectID   string `json:"projectId"`
}

// IssuePayload issue payload response
type IssuePayload struct {
	Success bool `json:"success"`
}

// CreateIssue creates an issue on linear board
func (c *Client) CreateIssue(ctx context.Context, in IssueCreateInput) (IssuePayload, error) {
	gqlReq := graphql.NewRequest(issueCreateMutation)
	gqlReq.Var("input", in)

	var res IssuePayload
	if err := c.gql.Run(ctx, gqlReq, res); err != nil {
		return IssuePayload{}, err
	}

	return res, nil
}
