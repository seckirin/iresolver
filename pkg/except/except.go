package except

import "errors"

// ErrServerUnreachable is returned when a server cannot be reached.
var ErrServerUnreachable = errors.New("target server is unreachable")

// ErrFileUnreachable is returned when a file cannot be accessed.
var ErrFileUnreachable = errors.New("target file is unreachable")
