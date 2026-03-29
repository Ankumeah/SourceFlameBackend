package git

import "errors"

var Error_Repository_Exists = errors.New("Repo alreday exists")
var Error_Unsupported_Service = errors.New("Unsupported Service")
var Error_Timeout = errors.New("Timeout")
