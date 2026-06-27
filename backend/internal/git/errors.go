package git

import (
	"errors"
	"fmt"
)

var Error_Repository_Exists = errors.New("Repo alreday exists")
var Error_Blob_Too_Large = errors.New("Blob too large")
var Error_Path_Too_Deep = errors.New("Path too deep")

var Error_Not_Found = errors.New("Not found")
var Error_Blob_Not_Found = fmt.Errorf("Blob not found: %w", Error_Not_Found)
var Error_Path_Not_Found = fmt.Errorf("Path not found: %w", Error_Not_Found)
var Error_Inavlid_Commit_Hash = fmt.Errorf("Invalid commit hash: %w", Error_Not_Found)
