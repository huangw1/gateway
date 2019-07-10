/**
 * @Author: huangw1
 * @Date: 2019/7/11 10:09
 */

package sd

type Subscriber interface {
	Hosts() ([]string, error)
}

type FixedSubscriber []string

func (f FixedSubscriber) Hosts() ([]string, error) {
	return f, nil
}
