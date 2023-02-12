package generator

import "errors"

var (
	errUnknownKind     = errors.New("unknown kind")
	errUnknownFileType = errors.New("unknown file type")
	errArgsIsNil       = errors.New("args is nil")
	errArgsUnavailable = errors.New("args converting is unavailable")
)
