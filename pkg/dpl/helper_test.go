package dpl

type trivialComponent struct {
	ComponentName string
	Data          map[string][]string
}

func (tc *trivialComponent) Name() string {
	return tc.ComponentName
}

func (tc *trivialComponent) KeyNames() []string {
	ret := []string{}
	for key := range tc.Data {
		ret = append(ret, key)
	}
	return ret
}

func (tc *trivialComponent) GetValues(key string) []string {
	return tc.Data[key]
}

func (tc *trivialComponent) ExpandValues(key string) ([]string, error) {
	return tc.GetValues(key), nil
}

func (tc *trivialComponent) SetValues(string, []string) {
}

func (tc *trivialComponent) EraseKey(string) {
}

func (tc *trivialComponent) GetSourceDir() string {
	return ""
}

func (tc *trivialComponent) GetWorkDir() string {
	return ""
}
