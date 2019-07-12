/**
 * @Author: huangw1
 * @Date: 2019/7/11 20:13
 */

package proxy

import (
	"context"
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/sd"
	"net/url"
	"time"
)

func NewRoundRobinLoadBalancedMiddleware(remote *config.Backend) Middleware {
	return NewRoundRobinLoadBalancedMiddlewareWithSubscriber(sd.FixedSubscriber(remote.Host))
}

func NewRandomLoadBalancedMiddleware(remote *config.Backend) Middleware {
	return NewRandomLoadBalancedMiddlewareWithSubscriber(sd.FixedSubscriber(remote.Host))
}

func NewRoundRobinLoadBalancedMiddlewareWithSubscriber(subscriber sd.Subscriber) Middleware {
	return newLoadBalancedMiddleware(sd.NewRoundRobinLB(subscriber))
}

func NewRandomLoadBalancedMiddlewareWithSubscriber(subscriber sd.Subscriber) Middleware {
	return newLoadBalancedMiddleware(sd.NewRandomLB(subscriber, time.Now().UnixNano()))
}

func newLoadBalancedMiddleware(balancer sd.Balancer) Middleware {
	return func(next ...Proxy) Proxy {
		if len(next) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			host, err := balancer.Host()
			if err != nil {
				return nil, err
			}
			req := request.Clone()
			buff := make([]byte, 0)
			buff = append(buff, host...)
			buff = append(buff, req.Path...)
			req.URL, err = url.Parse(string(buff))
			if err != nil {
				return nil, err
			}
			req.URL.RawQuery = req.Query.Encode()
			return next[0](ctx, &req)
		}
	}
}
