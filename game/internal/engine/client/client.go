package client

import "context"

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (e *Client) Run(ctx context.Context) error {
	return nil
}

func (e *Client) Shutdown(ctx context.Context) error {
	return nil
}
