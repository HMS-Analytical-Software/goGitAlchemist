// Package alchemist contains all the elements to do alchemistry.
package alchemist

// emails provides dummy email addresses for git.
var email = map[string]string{
	"red":       "richard@pw-compa.ny",
	"blue":      "betty@pw-compa.ny",
	"green":     "garry@pw-compa.ny",
	"api":       "api@pw-compa.ny",
	"blacklist": "blacklist@pw-compa.ny",
	"config":    "config@pw-compa.ny",
}

// authors provides dummy author names for git.
var author = map[string]string{
	"red":       "Richard Red",
	"blue":      "Betty Blue",
	"green":     "Garry Green",
	"api":       "Alissa Api",
	"blacklist": "Benjamin Blacklist",
	"config":    "Carry Config",
}

// getAuthor tries to get the author with mail address.
// if is not known, the name is returned as is.
func getAuthor(name string) string {
	author, ok := author[name]
	if !ok {
		return name
	}
	return author + " <" + email[name] + ">"
}

// defaultUser is used when a new repo is initialzeed.
const defaultUser = "red"

// FormulaFileName defines the name of gitalchemist formula files.
const FormulaFileName = "gitalchemist.yaml"

// dirMode defines the access mode for new directories.
const dirMode = 0755

// defaultBranch is the git default branch name.
const defaultBranch = "main"

// git commands on linux and windows.
const (
	linuxGitCmd   = "git"
	windowsGitCmd = "git.exe"
)
