// Code generated by Kitex v0.12.3. DO NOT EDIT.

package brokerservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	api "github.com/magicnana999/im/api/kitex_gen/api"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Deliver(ctx context.Context, Req *api.DeliverRequest, callOptions ...callopt.Option) (r *api.DeliverReply, err error)
}

// NewClient creates a client for the cmd_service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kBrokerServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the cmd_service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kBrokerServiceClient struct {
	*kClient
}

func (p *kBrokerServiceClient) Deliver(ctx context.Context, Req *api.DeliverRequest, callOptions ...callopt.Option) (r *api.DeliverReply, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Deliver(ctx, Req)
}
