package git

import "errors"

var Error_Repository_Exists = errors.New("Repo alreday exists")
var Error_Inavlid_Commit_Hash = errors.New("Invalid commit hash")
var Error_Commit_Not_Found = errors.New("Commit not found")
var Error_Blob_Not_Found = errors.New("Blob not found")
var Error_Blob_Too_Large = errors.New("Blob too large")
