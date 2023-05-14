// embed.go sits at the root of the file because the doc file is inaccessible
// if dir to embedded is not in a subdirectory or the same directory

package api

import "embed"

var (
	// Docs are needed so that binary files can access the docs folder
	//go:embed docs/*
	Docs embed.FS
)
