package client

import (
	"context"
	"github.com/alibabacloud-go/tea/dara"
)

type discoveryRequest struct{}

func (*discoveryRequest) Validate() error { return nil }
func (c *Client) GetDiscovery() (*DiscoveryResponse, error) {
	return c.GetDiscoveryWithOptions(&dara.RuntimeOptions{})
}
func (c *Client) GetDiscoveryWithOptions(runtime *dara.RuntimeOptions) (*DiscoveryResponse, error) {
	return c.GetDiscoveryWithContext(context.Background(), runtime)
}
func (c *Client) GetDiscoveryWithContext(ctx context.Context, runtime *dara.RuntimeOptions) (*DiscoveryResponse, error) {
	params := &Params{Action: dara.String("Discovery"), Pathname: dara.String("/discovery"), Method: dara.String("GET"), AuthType: dara.String("Anonymous"), ReadOnly: dara.Bool(true), Source: dara.String("platform_discovery")}
	result := &DiscoveryResponse{}
	err := c.CallTyped(ctx, params, &discoveryRequest{}, runtime, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
