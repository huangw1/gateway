/**
 * @Author: huangw1
 * @Date: 2019/7/12 14:56
 */

package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/core"
	"github.com/huangw1/gateway/proxy"
	"net/http"
	"strings"
	"time"
)

var ErrInternalError = errors.New("internal server error")

type HandlerFactory func(configuration *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc

func EndpointHandler(configuration *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
	endpointTimeout := time.Duration(configuration.Timeout) * time.Millisecond
	return func(c *gin.Context) {
		requestCtx, cancel := context.WithTimeout(c, endpointTimeout)
		c.Header(core.GatewayHeaderValue, core.GatewayVersion)

		res, err := proxy(requestCtx, NewRequest(c, configuration.QueryString))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			cancel()
			return
		}

		select {
		case <-requestCtx.Done():
			c.AbortWithError(http.StatusInternalServerError, ErrInternalError)
			cancel()
		default:
		}

		if configuration.CacheTTL.Seconds() != 0 && res != nil && res.IsComplete {
			c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(configuration.CacheTTL.Seconds())))
			c.JSON(http.StatusOK, res.Data)
			cancel()
			return
		}
		if res != nil {
			c.JSON(http.StatusOK, res.Data)
		} else {
			c.JSON(http.StatusOK, gin.H{})
		}
		cancel()
	}
}

var (
	headersToSend        = []string{"Content-Type"}
	userAgentHeaderValue = []string{core.GatewayUserAgent}
)

func NewRequest(c *gin.Context, queryString []string) *proxy.Request {
	params := make(map[string]string, len(c.Params))
	for _, param := range c.Params {
		params[strings.Title(param.Key)] = param.Value
	}
	headers := make(map[string][]string, len(headersToSend)+2)
	headers["X-Forwarded-For"] = []string{c.ClientIP()}
	headers["User-Agent"] = userAgentHeaderValue
	for _, k := range headersToSend {
		if v, ok := c.Request.Header[k]; ok {
			headers[k] = v
		}
	}

	query := make(map[string][]string, len(queryString))
	for i := range queryString {
		if v := c.Request.URL.Query().Get(queryString[i]); v != "" {
			query[queryString[i]] = []string{v}
		}
	}

	return &proxy.Request{
		Method:  c.Request.Method,
		Query:   query,
		Body:    c.Request.Body,
		Params:  params,
		Headers: headers,
	}
}
