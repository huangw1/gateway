/**
 * @Author: huangw1
 * @Date: 2019/7/12 14:39
 */

package gin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/logging"
	"github.com/huangw1/gateway/proxy"
	"github.com/huangw1/gateway/router"
	"net/http"
)

type Config struct {
	Engine         *gin.Engine
	Middlewares    []gin.HandlerFunc
	HandlerFactory HandlerFactory
	ProxyFactory   proxy.Factory
	Logger         logging.Logger
}

type factory struct {
	cfg Config
}

func (f factory) New() router.Router {
	return ginRouter{f.cfg, context.Background()}
}

func (f factory) NewWithContext(ctx context.Context) router.Router {
	return ginRouter{f.cfg, ctx}
}

func NewFactory(cfg Config) router.Factory {
	return factory{cfg}
}

func DefaultFactory(f proxy.Factory, logger logging.Logger) router.Factory {
	return factory{Config{
		Engine:         gin.Default(),
		Middlewares:    []gin.HandlerFunc{gin.Recovery()},
		HandlerFactory: EndpointHandler,
		ProxyFactory:   f,
		Logger:         logger,
	}}
}

type ginRouter struct {
	cfg Config
	ctx context.Context
}

func (r ginRouter) Run(cfg config.ServerConfig) {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.cfg.Logger.Info("Debug enabled")
	}
	r.cfg.Engine.RedirectTrailingSlash = true
	r.cfg.Engine.RedirectFixedPath = true
	r.cfg.Engine.HandleMethodNotAllowed = true
	r.cfg.Engine.Use(r.cfg.Middlewares...)

	if cfg.Debug {
		r.registerDebugEndpoints()
	}
	r.registerEndpoints(cfg.Endpoints)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r.cfg.Engine,
	}

	go func() {
		r.cfg.Logger.Critical(srv.ListenAndServe())
	}()

	<-r.ctx.Done()
	r.cfg.Logger.Error(srv.Shutdown(context.Background()))
}

func (r ginRouter) registerEndpoints(endpoints []*config.EndpointConfig) {
	for _, c := range endpoints {
		proxyStack, err := r.cfg.ProxyFactory.New(c)
		if err != nil {
			r.cfg.Logger.Error("Calling the proxy.Factory", err.Error())
			continue
		}
		r.registerEndpoint(c.Method, c.Endpoint, r.cfg.HandlerFactory(c, proxyStack), len(c.Backend))
	}
}

func (r ginRouter) registerEndpoint(method, path string, handler gin.HandlerFunc, count int) {
	if method != http.MethodGet && count > 1 {
		r.cfg.Logger.Error(method, "endpoints must have a single backend! Ignoring", path)
		return
	}
	switch method {
	case "GET":
		r.cfg.Engine.GET(path, handler)
	case "POST":
		r.cfg.Engine.POST(path, handler)
	case "PUT":
		r.cfg.Engine.PUT(path, handler)
	case "PATCH":
		r.cfg.Engine.PATCH(path, handler)
	case "DELETE":
		r.cfg.Engine.DELETE(path, handler)
	default:
		r.cfg.Logger.Error("Unsupported method", method)
	}
}

func (r ginRouter) registerDebugEndpoints() {
	handler := DebugHandler(r.cfg.Logger)
	r.cfg.Engine.GET("/__debug/*param", handler)
	r.cfg.Engine.POST("/__debug/*param", handler)
	r.cfg.Engine.PUT("/__debug/*param", handler)
}
