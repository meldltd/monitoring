package spec

import (
	"errors"
)

var ConnectionFailed = errors.New("Connection failed")
var SessionFailed = errors.New("SSH Session failed")
var QueryFailed = errors.New("Query failed")
var ExpectFailed = errors.New("Expect failed")
var ExpectUndefined = errors.New("Expect is nil")
var NoCheckPerformed = errors.New("No check performed")
var NotImplemented = errors.New("Not implemented")
var WillExpire = errors.New("Will expire")
