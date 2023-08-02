package storage

import "errors"

// use this to abort scans without giving an error
var EOD = errors.New("EOD")
