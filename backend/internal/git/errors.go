package git

import (
	"errors"
	"fmt"
)

var ErrRepositoryExists = errors.New("Repo alreday exists")
var ErrBlobTooLarge = errors.New("Blob too large")
var ErrPathTooDeep = errors.New("Path too deep")

var ErrNotFound = errors.New("Not found")
var ErrBlobNotFound = fmt.Errorf("Blob not found: %w", ErrNotFound)
var ErrPathNotFound = fmt.Errorf("Path not found: %w", ErrNotFound)
var ErrBranchNotFound = fmt.Errorf("Branch not found: %w", ErrNotFound)
var ErrInvalidCommitHash = fmt.Errorf("Invalid commit hash: %w", ErrNotFound)
