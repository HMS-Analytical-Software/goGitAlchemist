//go:build !windows

package main

// os, non-windows specific external program calls.
const (
	gitAlchemistCmd = "./gitalchemist"
	goCmd           = "go"
	gitCmd          = "git"
)
