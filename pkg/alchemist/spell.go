package alchemist

// spellHint provides the information for casting some spells.
// It can be used to create lists of arguments and process them
// in a loop.
type spellHint struct {
	dir  string
	args []string
}
