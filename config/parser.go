/**
 * @Author: huangw1
 * @Date: 2019/7/10 18:14
 */

package config

type Parser interface {
	Parse(filename string) (ServerConfig, error)
}
