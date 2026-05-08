package database

import "errors"

var Error_limit_too_big = errors.New("Limit too big")

var Error_Invalid = errors.New("Invalid")
var Error_invalid_user = errors.Join(errors.New("Invalid user"), Error_Invalid)
var Error_invalid_repo = errors.Join(errors.New("Invalid repo"), Error_Invalid)
var Error_invalid_pat = errors.Join(errors.New("Invalid pat"), Error_Invalid)
