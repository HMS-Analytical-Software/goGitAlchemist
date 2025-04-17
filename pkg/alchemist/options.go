package alchemist

// Options defines the options provided to the execution of the commands.
type Options struct {
	TaskDir       string // directory where the definition is located
	RepoDir       string // directory where the bare repository is located
	CfgDir        string // directory of the configuration definitions
	Verbose       bool   // verbose logging
	Test          bool   // test mode
	ExecuteSpells int    // execute only the first # steps

	// set during processing
	taskName       string // name of the task to execute
	cloneTo        string // directory of the repository clone
	numberOfSpells int    // number of steps
	currentSpell   int    // number of the current step (1-based)
}
