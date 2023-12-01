package storage

import "errors"


// TODO improve error handling...
var ConnectionErr = errors.New("storage not ready, no connection established")
