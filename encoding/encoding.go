/**
 * @Author: huangw1
 * @Date: 2019/7/10 19:59
 */

package encoding

import "io"

type Decoder func(r io.Reader, v *map[string]interface{}) error
