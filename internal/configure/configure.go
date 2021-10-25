package configure

type Flags struct {
	BuildDir         string
	BuildDirBasename string
	ConfigFile       string
	Overrides        []string
	Profiles         []string
	RootDir          string
}

func DoConfigure(flags Flags, args []string) {
}
