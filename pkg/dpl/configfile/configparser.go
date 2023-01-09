package configfile

import (
	"gopkg.in/ini.v1"
)

var (
	options ini.LoadOptions = ini.LoadOptions{
		IgnoreInlineComment: true,
		AllowShadows:        true,
	}
)

func LoadRawConfig(data []byte) (*ini.File, error) {
	return ini.LoadSources(options, data)
}

func LoadProjectConfig(path string) (*ini.File, error) {
	return ini.LoadSources(options, path)
}
