/**
 * @Author: huangw1
 * @Date: 2019/7/10 18:14
 */

package config

import (
	"github.com/huangw1/gateway/encoding"
	"time"
)

type ServerConfig struct {
	Endpoints []*EndpointConfig `mapstructure:"endpoints"`
	Timeout   time.Duration     `mapstructure:"timeout"`
	CacheTTL  time.Duration     `mapstructure:"cache_ttl"`
	Host      []string          `mapstructure:"host"`
	Port      int               `mapstructure:"port"`
	Version   int               `mapstructure:"version"`
	Debug     bool
}

type EndpointConfig struct {
	Endpoint        string        `mapstructure:"endpoint"`
	Method          string        `mapstructure:"method"`
	Backend         []*Backend    `mapstructure:"backend"`
	ConcurrentCalls int           `mapstructure:"concurrent_calls"`
	Timeout         time.Duration `mapstructure:"timeout"`
	CacheTTL        time.Duration `mapstructure:"cache_ttl"`
	QueryString     []string      `mapstructure:"querystring_params"`
}

type Backend struct {
	Group           string            `mapstructure:"group"`
	Method          string            `mapstructure:"method"`
	Host            []string          `mapstructure:"host"`
	URLPattern      string            `mapstructure:"url_pattern"`
	Blacklist       []string          `mapstructure:"blacklist"`
	Whitelist       []string          `mapstructure:"whitelist"`
	Mapping         map[string]string `mapstructure:"mapping"`
	Encoding        string            `mapstructure:"encoding"`
	Target          string            `mapstructure:"target"`
	URLKeys         []string
	ConcurrentCalls int
	Timeout         time.Duration
	Decoder         encoding.Decoder
}

func (s ServerConfig) Init() error {

}
