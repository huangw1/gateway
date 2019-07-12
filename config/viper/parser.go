/**
 * @Author: huangw1
 * @Date: 2019/7/10 18:14
 */

package viper

import (
	"github.com/huangw1/gateway/config"
	"github.com/spf13/viper"
)

func New() config.Parser {
	return parser{}
}

type parser struct {
	viper *viper.Viper
}

func (p parser) Parse(filename string) (config.ServerConfig, error) {
	p.viper = viper.New()
	p.viper.SetConfigFile(filename)
	p.viper.AutomaticEnv()
	var cfg config.ServerConfig
	if err := p.viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := p.viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	if err := cfg.Init(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
