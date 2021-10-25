package reconfigure

type Flags struct {
	Append    bool
	Overrides []string
	Profiles  []string
}

func DoReconfigure(flags Flags, args []string) {

}
