package generator

import "errors"

var (
	errUnknownKind     = errors.New("unknown kind")
	errUnknownFileType = errors.New("unknown file type")
	errArgsUnavailable = errors.New("args converting is unavailable")
)
