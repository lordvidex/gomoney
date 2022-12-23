// embed.go sits at the root of the file because the doc file is inaccessible to
// if it is not in a subdirectory or the same directory

package api

import "embed"

var (
	//go:embed docs/*
	Docs embed.FS
)
