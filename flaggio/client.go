package flaggio

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

const findFlags = `
	query findFlags($search: String, $offset: Int, $limit: Int) {
		flags(search: $search, offset: $offset, limit: $limit) {
			flags { id key name description enabled updatedAt createdAt }
			total
		}
	}
`

// Client defines mongo repository
type Client struct {
	gql *graphql.Client
}

// NewFlagMongo returns a mongo repository
func New(url string) *Client {
	return &Client{gql: graphql.NewClient(url)}
}

// FindFlagsOpt Allowed GQL params on findFlag queries
type FindFlagsOpt func(req *graphql.Request)

// WithFindFlagsSearchOpt Define a string to be searched on findFlag queries
func WithFindFlagsSearchOpt(search string) FindFlagsOpt {
	return func(req *graphql.Request) {
		req.Var("search", search)
	}
}

// FindFlagsByMaxAge Find flags that are not updated for more than duration param
func (r *Client) FindFlagsByMaxAge(ctx context.Context, maxAge time.Duration, opts ...FindFlagsOpt) ([]Flag, error) {
	count := 50
	page := 0
	total := 1 // FIXME: this is just to allow one run in the loop
	ageThreshold := time.Now().Add(maxAge * -1)

	flags := []Flag{}

	for page*count < total {
		gqlReq := graphql.NewRequest(findFlags)
		gqlReq.Var("search", "")
		for _, opt := range opts {
			opt(gqlReq)
		}

		gqlReq.Var("offset", page*count)
		gqlReq.Var("limit", count)

		var res FindFlagsResponse
		if err := r.gql.Run(ctx, gqlReq, &res); err != nil {
			return flags, err
		}

		for _, f := range res.Flags.Flags {
			if f.UpdatedAt.After(ageThreshold) {
				continue
			}

			flags = append(flags, f)
		}

		// This is asking for bugs
		if res.Flags.Total == 0 {
			break
		} else if total == 1 {
			total = res.Flags.Total
		}

		page++
	}

	return flags, nil
}
