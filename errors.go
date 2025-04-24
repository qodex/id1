package id1

import "errors"

var ErrNotFound = errors.New("not found")

var ErrExists = errors.New("item exists")

var ErrLimitExceeded = errors.New("limit exceeded")
