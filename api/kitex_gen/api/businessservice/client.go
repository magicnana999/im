// Code generated by Kitex v0.12.3. DO NOT EDIT.

package businessservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	api "github.com/magicnana999/im/api/kitex_gen/api"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Login(ctx context.Context, Req *api.LoginRequest, callOptions ...callopt.Option) (r *api.LoginReply, err error)
	Logout(ctx context.Context, Req *api.LogoutRequest, callOptions ...callopt.Option) (r *api.LogoutReply, err error)
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
	return &kBusinessServiceClient{
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

type kBusinessServiceClient struct {
	*kClient
}

func (p *kBusinessServiceClient) Login(ctx context.Context, Req *api.LoginRequest, callOptions ...callopt.Option) (r *api.LoginReply, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Login(ctx, Req)
}

func (p *kBusinessServiceClient) Logout(ctx context.Context, Req *api.LogoutRequest, callOptions ...callopt.Option) (r *api.LogoutReply, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Logout(ctx, Req)
}
