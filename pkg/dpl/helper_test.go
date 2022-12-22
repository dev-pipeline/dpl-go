package dpl

type trivialComponent struct {
	ComponentName string
	Data          map[string][]string
}

func (tc *trivialComponent) Name() string {
	return tc.ComponentName
}

func (tc *trivialComponent) ValueNames() []string {
	ret := []string{}
	for key := range tc.Data {
		ret = append(ret, key)
	}
	return ret
}

func (tc *trivialComponent) GetValue(key string) []string {
	return tc.Data[key]
}

func (tc *trivialComponent) ExpandValue(key string) ([]string, error) {
	return tc.GetValue(key), nil
}

func (tc *trivialComponent) SetValue(string, []string) {
}

func (tc *trivialComponent) EraseValue(string) {
}

func (tc *trivialComponent) GetSourceDir() string {
	return ""
}

func (tc *trivialComponent) GetWorkDir() string {
	return ""
}
