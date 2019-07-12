/**
 * @Author: huangw1
 * @Date: 2019/7/12 17:18
 */

package main

import (
	"flag"
	"github.com/huangw1/gateway/config/viper"
	"github.com/huangw1/gateway/logging/gologging"
	"github.com/huangw1/gateway/proxy"
	"github.com/huangw1/gateway/router/gin"
	"log"
	"os"
)

func main() {
	port := flag.Int("p", 0, "Port of the server")
	level := flag.String("l", "INFO", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	filename := flag.String("c", "./etc/configuration.json", "Path to the configuration filename")
	flag.Parse()

	parser := viper.New()
	cfg, err := parser.Parse(*filename)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Debug = *debug
	if *port != 0 {
		cfg.Port = *port
	}
	logger, err := gologging.NewLogger(*level, os.Stdout, "[GIN]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	routerFactory := gin.DefaultFactory(proxy.DefaultFactory(logger), logger)
	routerFactory.New().Run(cfg)
}
