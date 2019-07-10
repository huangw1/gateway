/**
 * @Author: huangw1
 * @Date: 2019/7/11 20:03
 */

package proxy

import (
	"context"
	"github.com/huangw1/gateway/logging"
	"time"
)

func NewLoggingMiddleware(logger logging.Logger, name string) Middleware {
	return func(next ...Proxy) Proxy {
		if len(next) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			begin := time.Now()
			logger.Info(name, "calling to backend")
			logger.Debug("request", request)

			res, err := next[0](ctx, request)

			logger.Info(name, "call to backend took", time.Now().Sub(begin).String())
			if err != nil {
				logger.Warning(name, "call to backend failed:", err.Error())
			}
			if res == nil {
				logger.Warning(name, "call to backend returned a null response")
			}
			return res, err
		}
	}
}
