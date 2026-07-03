package mpesa

import (
	"context"
)

type Client struct {
	baseURL string
	apiKey  string
}

func (c *Client) PullTransactions(ctx context.Context) ([]APITransaction, error) {
	// Make a HTTP Request to the service provider
	return nil, nil
}
