/**
 * @Author: huangw1
 * @Date: 2019/7/12 09:38
 */

package proxy

import (
	"context"
	"errors"
	"github.com/huangw1/gateway/config"
	"time"
)

func NewConcurrentMiddleware(remote *config.Backend) Middleware {
	serviceTimeout := time.Duration(remote.Timeout.Nanoseconds()) * time.Millisecond
	return func(next ...Proxy) Proxy {
		if len(next) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			localCtx, cancel := context.WithTimeout(ctx, serviceTimeout)
			results := make(chan *Response, remote.ConcurrentCalls)
			failed := make(chan error, remote.ConcurrentCalls)
			for i := 0; i < remote.ConcurrentCalls; i++ {
				go processConcurrentCall(localCtx, request, next[0], results, failed)
			}

			var response *Response
			var err error
			for i := 0; i < remote.ConcurrentCalls; i++ {
				select {
				case response = <-results:
					if response != nil && response.IsComplete {
						cancel()
						return response, nil
					}
				case err = <-failed:
				case <-ctx.Done():
				}
			}
			cancel()
			return response, err
		}
	}
}

var ErrNullResult = errors.New("invalid response")

func processConcurrentCall(ctx context.Context, request *Request, next Proxy, results chan<- *Response, failed chan<- error) {
	localCtx, cancel := context.WithCancel(ctx)
	res, err := next(localCtx, request)
	if err != nil {
		failed <- err
		cancel()
		return
	}
	if res == nil {
		failed <- ErrNullResult
		cancel()
		return
	}
	select {
	case results <- res:
	case <-ctx.Done():
		failed <- ctx.Err()
	}
	cancel()
}
