/**
 * @Author: huangw1
 * @Date: 2019/7/11 09:43
 */

package encoding

import (
	"io"
	"encoding/json"
)

func JSONDecoder(r io.Reader, v *map[string]interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}
