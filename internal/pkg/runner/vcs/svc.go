package vcs

import "io"

type Options struct {
	Directory     string
	RepositoryURL string
	Branch        string
	Revision      string
	Progress      io.Writer
}

type VCS interface {
	Clone(options Options) error
}
