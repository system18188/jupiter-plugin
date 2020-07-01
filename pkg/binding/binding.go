package binding

import (
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON               = "application/json"
	MIMEHTML               = "text/html"
	MIMEXML                = "application/xml"
	MIMEXML2               = "text/xml"
	MIMEPlain              = "text/plain"
	MIMEPOSTForm           = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm  = "multipart/form-data"
	MIMEPROTOBUF           = "application/x-protobuf"
	MIMEMSGPACK            = "application/x-msgpack"
	MIMEMSGPACK2           = "application/msgpack"
	MIMEYAML               = "application/x-yaml"
	MIMEJPEN               = "image/jpeg"
	MIMEGIF                = "image/gif"
	MIMEPNG                = "image/png"
	MIMEDOC                = "application/msword"
	MIMEPDF                = "application/pdf"
	MIMEDOCX               = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MIMExlsx               = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MIMExls                = "application/vnd.ms-excel"
	MIMEpptx               = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MIME7z                 = "application/x-7z-compressed"
	MIMErar                = "application/x-rar-compressed"
	MIMEZip                = "application/zip"
	HEADER_ContentType     = "Content-Type"
	HEADER_ContentEncoding = "Content-Encoding"

	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)

var MIMETypes = map[string]string{
	MIMEJSON: ".json",
	MIMEHTML: ".html",
	MIMEXML:  ".xml",
	MIMEXML2: ".xml",
	MIMEYAML: ".yaml",
	MIMEJPEN: ".jpg",
	MIMEGIF:  ".gif",
	MIMEPNG:  ".png",
	MIMEPDF:  ".pdf",
	MIME7z:   ".7z",
	MIMErar:  ".rar",
	MIMEDOC:  ".doc",
	MIMEDOCX: ".docx",
	MIMEpptx: ".pptx",
	MIMExlsx: ".xlsx",
	MIMExls:  ".xls",
	MIMEZip:  ".zip",
}

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar with Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

// BindingUri adds BindUri method to Binding. BindUri is similar with Bind,
// but it read the Params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

// BindingArray adds BindArray method to Binding. BindArray is similar with Bind,
type BindingArray interface {
	Name() string
	BindArray(map[string]interface{}, interface{}) error
}


// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v9
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() interface{}

	// AddValidation 添加一个新的验证方法
	AddValidation(tag string, fn func(fl validator.FieldLevel) bool) error

	// AddTranslation 添加一个新的验证提示
	AddTranslation(tag string, addfn func(ut ut.Translator) error, translatorFn func(ut ut.Translator, fe validator.FieldError) string) error

	// AddValidateType 添加一个新的类型
	AddValidateType(CustomTypeFunc func(field reflect.Value) interface{}, types ...interface{})
}

// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v9
// under the hood.
var Validator StructValidator = &defaultValidator{}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Array         = arrayBinding{}
	RegisterFormType = NewBindType()
)

func NewArray() BindingArray {
	return Array
}

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	}
	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	case MIMEYAML:
		return YAML
	default: 	//case MIMEPOSTForm, MIMEMultipartPOSTForm:
		return Form
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
