package binding

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	// Write Default Value
	if err := WriteDefaultValueOnTag(obj, "form");err != nil {
		return fmt.Errorf("WriteDefaultValueOnTag Error: %v",err.Error())
	}
	return decodeXML(req.Body, obj)
}

func (xmlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body), obj)
}

func decodeXML(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
