package encoding

import (
	"encoding/xml"
	"io"
)

const XML = "xml"

/*
["a","b"] to {"collection":["a","b"]}
Backend config adds:
"mapping": {
	"collection": "list"
}
*/
func NewXMLDecoder(isCollection bool) Decoder {
	if isCollection {
		return XMLCollectionDecoder
	}
	return XMLDecoder
}

func XMLDecoder(r io.Reader, v *map[string]interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}

func XMLCollectionDecoder(r io.Reader, v *map[string]interface{}) error {
	var collection []interface{}
	d := xml.NewDecoder(r)
	err := d.Decode(&collection)
	if err != nil {
		return err
	}
	*v = map[string]interface{}{"collection": collection}
	return nil
}
