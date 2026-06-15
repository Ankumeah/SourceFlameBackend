package database

import "errors"

var Error_limit_too_big = errors.New("Limit too big")

var Safe_Error = errors.New("Error")

var Error_Invalid = errors.Join(errors.New("Invalid"), Safe_Error)
var Error_invalid_user = errors.Join(errors.New("Invalid user"), Error_Invalid)
var Error_invalid_repo = errors.Join(errors.New("Invalid repo"), Error_Invalid)
var Error_invalid_pat = errors.Join(errors.New("Invalid pat"), Error_Invalid)

var Error_Exists = errors.Join(errors.New("Exists"), Safe_Error)
var Error_user_exists = errors.Join(errors.New("User exists"), Error_Exists)
var Error_repo_exists = errors.Join(errors.New("Repo exists"), Error_Exists)
var Error_pat_exists = errors.Join(errors.New("PAT exists"), Error_Exists)
