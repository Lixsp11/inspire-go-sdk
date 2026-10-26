package client

import (
	"context"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/alibabacloud-go/tea/dara"
	"net/url"
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
func (c *Client) ListJobs(request *ListJobsRequest) (*ListJobsResponse, error) {
	return c.ListJobsWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) ListJobsWithOptions(request *ListJobsRequest, runtime *dara.RuntimeOptions) (*ListJobsResponse, error) {
	return c.ListJobsWithContext(context.Background(), request, runtime)
}
func (c *Client) ListJobsWithContext(ctx context.Context, request *ListJobsRequest, runtime *dara.RuntimeOptions) (*ListJobsResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("ListJobs"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/train"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(true), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/distributedTraining")}
	result := &ListJobsResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) GetJob(request *GetJobRequest) (*GetJobResponse, error) {
	return c.GetJobWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) GetJobWithOptions(request *GetJobRequest, runtime *dara.RuntimeOptions) (*GetJobResponse, error) {
	return c.GetJobWithContext(context.Background(), request, runtime)
}
func (c *Client) GetJobWithContext(ctx context.Context, request *GetJobRequest, runtime *dara.RuntimeOptions) (*GetJobResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("GetJob"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/train"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(true), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/distributedTraining")}
	p.Referer = dara.String("/jobs/distributedTrainingDetail/" + url.PathEscape(dara.StringValue(request.JobId)))
	result := &GetJobResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) CreateJobConsole(request *CreateJobConsoleRequest) (*CreateJobConsoleResponse, error) {
	return c.CreateJobConsoleWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) CreateJobConsoleWithOptions(request *CreateJobConsoleRequest, runtime *dara.RuntimeOptions) (*CreateJobConsoleResponse, error) {
	return c.CreateJobConsoleWithContext(context.Background(), request, runtime)
}
func (c *Client) CreateJobConsoleWithContext(ctx context.Context, request *CreateJobConsoleRequest, runtime *dara.RuntimeOptions) (*CreateJobConsoleResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("CreateJobConsole"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/train"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(false), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/distributedTraining")}
	result := &CreateJobConsoleResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) StopJob(request *StopJobRequest) (*StopJobResponse, error) {
	return c.StopJobWithOptions(request, &dara.RuntimeOptions{})
}
func (c *Client) StopJobWithOptions(request *StopJobRequest, runtime *dara.RuntimeOptions) (*StopJobResponse, error) {
	return c.StopJobWithContext(context.Background(), request, runtime)
}
func (c *Client) StopJobWithContext(ctx context.Context, request *StopJobRequest, runtime *dara.RuntimeOptions) (*StopJobResponse, error) {
	if c == nil || c.Client == nil || request == nil {
		return nil, openapi.ValidationError("Initialized client and request are required.")
	}
	p := &openapi.Params{Action: dara.String("StopJob"), Version: dara.String("v2"), Protocol: dara.String("HTTPS"), Pathname: dara.String("/api/v2/train"), Method: dara.String("POST"), AuthType: dara.String("SSO"), Style: dara.String("ROA"), ReqBodyType: dara.String("json"), BodyType: dara.String("json"), ReadOnly: dara.Bool(false), Source: dara.String("inspire_reference"), Referer: dara.String("/jobs/distributedTraining")}
	result := &StopJobResponse{}
	if err := c.Client.CallTyped(ctx, p, request, runtime, result); err != nil {
		return nil, err
	}
	return result, nil
}
