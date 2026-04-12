package database

import "errors"

var Error_invalid_user = errors.New("Invalid user")
var Error_invalid_repo = errors.New("Invalid repo")
var Error_invalid_pat = errors.New("Invalid repo")
var Error_limit_too_big = errors.New("Limit too big")
