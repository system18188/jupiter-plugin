package binding

type arrayBinding struct{}

func (arrayBinding) Name() string {
	return "array"
}

func (arrayBinding) BindArray(m map[string]interface{}, obj interface{}) error {
	if err := mapArray(obj, m); err != nil {
		return err
	}
	return validate(obj)
}
