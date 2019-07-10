/**
 * @Author: huangw1
 * @Date: 2019/7/11 11:19
 */

package proxy

import (
	"context"
	"errors"
	"github.com/huangw1/gateway/config"
)

var (
	ErrNoBackends       = errors.New("all endpoints must have at least one backend")
	ErrTooManyBackends  = errors.New("too many backends for this proxy")
	ErrTooManyProxies   = errors.New("too many proxies for this proxy middleware")
	ErrNotEnoughProxies = errors.New("not enough proxies for this endpoint")
)

type Response struct {
	Data       map[string]interface{}
	IsComplete bool
}

type Proxy func(ctx context.Context, request *Request) (*Response, error)

type BackendFactory func(remote *config.Backend) Proxy

type Middleware func(next ... Proxy) Proxy

func EmptyMiddleware(next ...Proxy) Proxy {
	if len(next) > 1 {
		panic(ErrTooManyProxies)
	}
	return next[0]
}
