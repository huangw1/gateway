/**
 * @Author: huangw1
 * @Date: 2019/7/12 11:14
 */

package proxy

import (
	"github.com/huangw1/gateway/config"
	"github.com/huangw1/gateway/logging"
	"github.com/huangw1/gateway/sd"
)

type Factory interface {
	New(cfg *config.EndpointConfig) (Proxy, error)
}

func DefaultFactory(logger logging.Logger) Factory {
	return NewDefaultFactory(DefaultHTTPProxy, logger)
}

func NewDefaultFactory(backendFactory BackendFactory, logger logging.Logger) Factory {
	return NewDefaultFactoryWithSubscriber(backendFactory, logger, sd.FixedSubscriberFactory)
}

func DefaultFactoryWithSubscriber(logger logging.Logger, sF sd.SubscriberFactory) Factory {
	return NewDefaultFactoryWithSubscriber(DefaultHTTPProxy, logger, sF)
}

func NewDefaultFactoryWithSubscriber(backendFactory BackendFactory, logger logging.Logger, sF sd.SubscriberFactory) Factory {
	return defaultFactory{backendFactory, logger, sF}
}

type defaultFactory struct {
	backendFactory    BackendFactory
	logger            logging.Logger
	subscriberFactory sd.SubscriberFactory
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
	return f.newStack(cfg.Backend[0]), nil
}

func (f defaultFactory) newMulti(cfg *config.EndpointConfig) (p Proxy, err error) {
	backendProxy := make([]Proxy, len(cfg.Backend))

	for i, backend := range cfg.Backend {
		backendProxy[i] = f.newStack(backend)
	}
	p = NewMergeDataMiddleware(cfg)(backendProxy...)
	return
}

func (f defaultFactory) newStack(backend *config.Backend) (p Proxy) {
	p = f.backendFactory(backend)
	p = NewRoundRobinLoadBalancedMiddlewareWithSubscriber(f.subscriberFactory(backend))(p)
	if backend.ConcurrentCalls > 1 {
		p = NewConcurrentMiddleware(backend)(p)
	}
	p = NewLoggingMiddleware(f.logger, "[GATEWAY]")(p)
	p = NewRequestBuilderMiddleware(backend)(p)
	return
}
