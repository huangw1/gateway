/**
 * @Author: huangw1
 * @Date: 2019/7/12 14:37
 */

package router

import (
	"context"
	"github.com/huangw1/gateway/config"
)

type Router interface {
	Run(cfg config.ServerConfig)
}

type Factory interface {
	New() Router
	NewWithContext(ctx context.Context) Router
}
