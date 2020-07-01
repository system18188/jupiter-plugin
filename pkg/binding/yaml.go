package binding

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (yamlBinding) Bind(req *http.Request, obj interface{}) error {
	// Write Default Value
	if err := WriteDefaultValueOnTag(obj, "form");err != nil {
		return fmt.Errorf("WriteDefaultValueOnTag Error: %v",err.Error())
	}
	return decodeYAML(req.Body, obj)
}

func (yamlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeYAML(bytes.NewReader(body), obj)
}

func decodeYAML(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
