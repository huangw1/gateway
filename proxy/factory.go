/**
 * @Author: huangw1
 * @Date: 2019/7/12 11:14
 */

package proxy

import (
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/logging"
)

type Factory interface {
	New(cfg *config.EndpointConfig) (Proxy, error)
}

func DefaultFactory(logger logging.Logger) Factory {
	return NewDefaultFactory(httpProxy, logger)
}

func NewDefaultFactory(backendFactory BackendFactory, logger logging.Logger) Factory {
	return defaultFactory{backendFactory, logger}
}

type defaultFactory struct {
	backendFactory BackendFactory
	logger         logging.Logger
}

func (f defaultFactory) New(cfg *config.EndpointConfig) (p Proxy, err error) {
	switch len(cfg.Backend) {
	case 0:
		err = ErrNoBackends
	case 1:
		p, err = f.newSingle(cfg)
	default:
		p, err = f.newMulti(cfg)
	}
	return
}

func (f defaultFactory) newSingle(cfg *config.EndpointConfig) (p Proxy, err error) {
	p = f.backendFactory(cfg.Backend[0])
	p = NewRoundRobinLoadBalancedMiddleware(cfg.Backend[0])(p)
	if cfg.Backend[0].ConcurrentCalls > 1 {
		p = NewConcurrentMiddleware(cfg.Backend[0])(p)
	}
	p = NewRequestBuilderMiddleware(cfg.Backend[0])(p)
	return
}

func (f defaultFactory) newMulti(cfg *config.EndpointConfig) (p Proxy, err error) {
	backendProxy := make([]Proxy, len(cfg.Backend))

	for i, backend := range cfg.Backend {
		backendProxy[i] = f.backendFactory(backend)
		backendProxy[i] = NewRoundRobinLoadBalancedMiddleware(backend)(backendProxy[i])
		if backend.ConcurrentCalls > 1 {
			backendProxy[i] = NewConcurrentMiddleware(backend)(backendProxy[i])
		}
		backendProxy[i] = NewRequestBuilderMiddleware(backend)(backendProxy[i])
	}
	// todo merge
	return
}
