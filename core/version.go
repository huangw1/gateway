/**
 * @Author: huangw1
 * @Date: 2019/7/12 16:51
 */

package core

import "fmt"

const Gateway = "X-GATEWAY"

var (
	GatewayVersion     = "undefined"
	GatewayHeaderValue = fmt.Sprintf("Version %s", GatewayVersion)
	GatewayUserAgent   = fmt.Sprintf("Gateway Version %s", GatewayVersion)
)
