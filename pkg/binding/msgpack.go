package binding

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ugorji/go/codec"
)

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (msgpackBinding) Bind(req *http.Request, obj interface{}) error {
	// Write Default Value
	if err := WriteDefaultValueOnTag(obj, "form");err != nil {
		return fmt.Errorf("WriteDefaultValueOnTag Error: %v",err.Error())
	}
	return decodeMsgPack(req.Body, obj)
}

func decodeMsgPack(r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return validate(obj)
}
