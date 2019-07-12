/**
 * @Author: huangw1
 * @Date: 2019/7/11 10:10
 */

package sd

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

var ErrNoHosts = errors.New("no hosts available")

type Balancer interface {
	Host() (string, error)
}

func NewRoundRobinLB(subscriber Subscriber) Balancer {
	return roundRobinLB{
		subscriber: subscriber,
		counter:    0,
	}
}

type roundRobinLB struct {
	subscriber Subscriber
	counter    uint64
}

func (r roundRobinLB) Host() (string, error) {
	hosts, err := r.subscriber.Hosts()
	if err != nil {
		return "", err
	}
	if len(hosts) <= 0 {
		return "", ErrNoHosts
	}
	offset := (atomic.AddUint64(&r.counter, 1) - 1) % uint64(len(hosts))
	return hosts[offset], nil
}

func NewRandomLB(subscriber Subscriber, seed int64) Balancer {
	return randomLB{
		subscriber: subscriber,
		rand:       rand.New(rand.NewSource(seed)),
	}
}

type randomLB struct {
	subscriber Subscriber
	rand       *rand.Rand
}

func (r randomLB) Host() (string, error) {
	hosts, err := r.subscriber.Hosts()
	if err != nil {
		return "", err
	}
	if len(hosts) <= 0 {
		return "", ErrNoHosts
	}
	return hosts[r.rand.Intn(len(hosts))], nil
}
