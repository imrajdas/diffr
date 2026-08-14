package diffr

var (
	Port             int
	Address          string
	Stdout           bool
	JSON             bool
	NoOpen           bool
	PatchFile        string
	Exclude          []string
	IgnoreFile       string
	NoDefaultExclude bool
	NoGitignore      bool
	IgnoreWhitespace bool
	IgnoreCase       bool
	Context          int
)
