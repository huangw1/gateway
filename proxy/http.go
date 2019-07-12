/**
 * @Author: huangw1
 * @Date: 2019/7/11 11:32
 */

package proxy

import (
	"context"
	"errors"
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/encoding"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

type HTTPClientFactory func(ctx context.Context) *http.Client

func NewHTTPClient(_ context.Context) *http.Client {
	return http.DefaultClient
}

func DefaultHTTPProxy(remote *config.Backend) Proxy {
	return NewHTTPProxy(remote, NewHTTPClient, remote.Decoder)
}

func NewRequestBuilderMiddleware(remote *config.Backend) Middleware {
	return func(next ...Proxy) Proxy {
		if len(next) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			r := request.Clone()
			r.GeneratePath(remote.URLPattern)
			r.Method = remote.Method
			return next[0](ctx, &r)
		}
	}
}

func NewHTTPProxy(remote *config.Backend, clientFactory HTTPClientFactory, decode encoding.Decoder) Proxy {
	formatter := NewEntityFormatter(remote.Target, remote.Whitelist, remote.Blacklist, remote.Group, remote.Mapping)
	return func(ctx context.Context, request *Request) (*Response, error) {
		req, err := http.NewRequest(request.Method, request.URL.String(), request.Body)
		if err != nil {
			return nil, err
		}
		req.Header = request.Headers
		res, err := clientFactory(ctx).Do(req.WithContext(ctx))
		defer res.Body.Close()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
			return nil, ErrInvalidStatusCode
		}
		var data map[string]interface{}
		err = decode(res.Body, &data)
		if err != nil {
			return nil, err
		}
		r := formatter.Format(Response{Data: data, IsComplete: true})
		return &r, nil
	}
}
