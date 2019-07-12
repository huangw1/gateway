/**
 * @Author: huangw1
 * @Date: 2019/7/11 09:43
 */

package encoding

import (
	"encoding/json"
	"io"
)

const JSON = "json"

/*
["a","b"] to {"collection":["a","b"]}
Backend config adds:
"mapping": {
	"collection": "list"
}
*/
func NewJSONDecoder(isCollection bool) Decoder {
	if isCollection {
		return JSONCollectionDecoder
	}
	return JSONDecoder
}

func JSONDecoder(r io.Reader, v *map[string]interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}

func JSONCollectionDecoder(r io.Reader, v *map[string]interface{}) error {
	var collection []interface{}
	d := json.NewDecoder(r)
	d.UseNumber()
	err := d.Decode(&collection)
	if err != nil {
		return err
	}
	*v = map[string]interface{}{"collection": collection}
	return nil
}
