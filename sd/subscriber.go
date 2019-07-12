/**
 * @Author: huangw1
 * @Date: 2019/7/11 10:09
 */

package sd

import "github.com/huangw1/gateway/config"

type Subscriber interface {
	Hosts() ([]string, error)
}

type FixedSubscriber []string

func (f FixedSubscriber) Hosts() ([]string, error) {
	return f, nil
}

type SubscriberFactory func(backend *config.Backend) Subscriber

func FixedSubscriberFactory(cfg *config.Backend) Subscriber {
	return FixedSubscriber(cfg.Host)
}
