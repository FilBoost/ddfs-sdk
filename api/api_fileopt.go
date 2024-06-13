package api

import (
	"io"
)

type RemoteFileOpt interface {
	Fetch() (io.ReadCloser, uint64, error)
	FetchWithConfirm() (io.ReadCloser, uint64, error)
	GetFetchUrl() (string, error)
	Revert() error
	Confirm(string) error
}

type RemoteFile interface {
	io.ReadCloser
	Size() (uint64, error)
}
