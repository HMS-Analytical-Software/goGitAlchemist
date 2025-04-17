//go:build !windows

package alchemist

// use the linux git command unless we compile for windows.
const gitCmd = linuxGitCmd
