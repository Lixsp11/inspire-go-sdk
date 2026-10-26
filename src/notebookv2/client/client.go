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
func (c *Client) ListNotebooks(request *ListNotebooksRequest) (*ListNotebooksResponse, error) {
	return c.ListNotebooksWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) ListNotebooksWithOptions(request *ListNotebooksRequest, runtime *dara.RuntimeOptions) (*ListNotebooksResponse, error) {
	return c.ListNotebooksWithContext(context.Background(), request, runtime)
}
func (c *Client) ListNotebooksWithContext(ctx context.Context, request *ListNotebooksRequest, runtime *dara.RuntimeOptions) (*ListNotebooksResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("ListNotebooks"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/notebook"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(true), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/interactiveModeling")}
	result := &ListNotebooksResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) GetNotebook(request *GetNotebookRequest) (*GetNotebookResponse, error) {
	return c.GetNotebookWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) GetNotebookWithOptions(request *GetNotebookRequest, runtime *dara.RuntimeOptions) (*GetNotebookResponse, error) {
	return c.GetNotebookWithContext(context.Background(), request, runtime)
}
func (c *Client) GetNotebookWithContext(ctx context.Context, request *GetNotebookRequest, runtime *dara.RuntimeOptions) (*GetNotebookResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("GetNotebook"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/notebook"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(true), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/interactiveModeling")}
	result := &GetNotebookResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
