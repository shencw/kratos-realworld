// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.1.4

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type AccountHTTPServer interface {
	Accounts(context.Context, *AccountsRequest) (*AccountsReply, error)
}

func RegisterAccountHTTPServer(s *http.Server, srv AccountHTTPServer) {
	r := s.Route("/")
	r.GET("/api/accounts", _Account_Accounts0_HTTP_Handler(srv))
}

func _Account_Accounts0_HTTP_Handler(srv AccountHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AccountsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/account.v1.Account/Accounts")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Accounts(ctx, req.(*AccountsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AccountsReply)
		return ctx.Result(200, reply)
	}
}

type AccountHTTPClient interface {
	Accounts(ctx context.Context, req *AccountsRequest, opts ...http.CallOption) (rsp *AccountsReply, err error)
}

type AccountHTTPClientImpl struct {
	cc *http.Client
}

func NewAccountHTTPClient(client *http.Client) AccountHTTPClient {
	return &AccountHTTPClientImpl{client}
}

func (c *AccountHTTPClientImpl) Accounts(ctx context.Context, in *AccountsRequest, opts ...http.CallOption) (*AccountsReply, error) {
	var out AccountsReply
	pattern := "/api/accounts"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/account.v1.Account/Accounts"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}