package dpl

type trivialComponent struct {
	ComponentName string
	Data          map[string][]string
}

func (self *trivialComponent) Name() string {
	return self.ComponentName
}

func (self *trivialComponent) ValueNames() []string {
	ret := []string{}
	for key := range self.Data {
		ret = append(ret, key)
	}
	return ret
}

func (self *trivialComponent) GetValue(key string) []string {
	return self.Data[key]
}

func (self *trivialComponent) ExpandValue(key string) ([]string, error) {
	return self.GetValue(key), nil
}

func (self *trivialComponent) SetValue(string, []string) {
}

func (self *trivialComponent) EraseValue(string) {
}
