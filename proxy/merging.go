/**
 * @Author: huangw1
 * @Date: 2019/7/12 11:59
 */

package proxy

import (
	"context"
	"github.com/huangw1/gateway/config"
	"time"
)

func NewMergeDataMiddleware(cfg *config.EndpointConfig) Middleware {
	totalBackends := len(cfg.Backend)
	if totalBackends == 0 {
		panic(ErrNoBackends)
	}
	if totalBackends == 1 {
		return EmptyMiddleware
	}
	serviceTimeout := time.Duration(cfg.Timeout.Nanoseconds()) * time.Millisecond
	return func(next ...Proxy) Proxy {
		if len(next) != totalBackends {
			panic(ErrNotEnoughProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			localCtx, cancel := context.WithTimeout(ctx, serviceTimeout)
			results := make(chan *Response, totalBackends)
			failed := make(chan error, totalBackends)
			for _, n := range next {
				go requestPart(localCtx, request, n, results, failed)
			}
			var err error
			isEmpty := true
			responses := make([]*Response, totalBackends)
			for i := 0; i < totalBackends; i++ {
				select {
				case err = <-failed:
				case responses[i] = <-results:
					isEmpty = false
				}
			}
			if isEmpty {
				cancel()
				return &Response{map[string]interface{}{}, false}, err
			}
			result := combineData(responses, totalBackends)
			cancel()
			return result, err
		}
	}
}

func requestPart(ctx context.Context, request *Request, next Proxy, results chan *Response, failed chan error) {
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

func combineData(responses []*Response, total int) *Response {
	composedData := make(map[string]interface{})
	isComplete := total == len(responses)
	for _, part := range responses {
		if part != nil && part.IsComplete {
			for k, v := range part.Data {
				composedData[k] = v
			}
			isComplete = isComplete && part.IsComplete
		} else {
			isComplete = false
		}
	}

	return &Response{composedData, isComplete}
}
