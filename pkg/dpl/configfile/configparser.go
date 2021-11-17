package configfile

import (
	"gopkg.in/ini.v1"
)

func LoadRawConfig(data []byte) (*ini.File, error) {
	return ini.ShadowLoad(data)
}

func LoadProjectConfig(path string) (*ini.File, error) {
	return ini.ShadowLoad(path)
}
