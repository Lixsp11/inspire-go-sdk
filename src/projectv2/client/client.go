package client

import (
	"context"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/alibabacloud-go/tea/dara"
)

type Client struct{ *openapi.Client }

func NewClient(config *openapi.Config) (*Client, error) {
	client := &Client{}
	err := client.Init(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func (c *Client) Init(config *openapi.Config) error {
	if c.Client != nil {
		return openapi.NewError("ConfigurationError", "AlreadyInitialized", "Client is already initialized.")
	}
	var err error
	c.Client, err = openapi.NewClient(config)
	return err
}
func (c *Client) ListProjects(request *ListProjectsRequest) (*ListProjectsResponse, error) {
	return c.ListProjectsWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) ListProjectsWithOptions(request *ListProjectsRequest, runtime *dara.RuntimeOptions) (*ListProjectsResponse, error) {
	return c.ListProjectsWithContext(context.Background(), request, runtime)
}
func (c *Client) ListProjectsWithContext(ctx context.Context, request *ListProjectsRequest, runtime *dara.RuntimeOptions) (*ListProjectsResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("ListProjects"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/project"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(true), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/interactiveModeling")}
	result := &ListProjectsResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
