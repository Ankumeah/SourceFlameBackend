package database

import (
	"errors"
	"fmt"
)

var ErrLimitTooLarge = errors.New("Limit too big")

var ErrSafe = errors.New("Error")

var ErrInvalid = fmt.Errorf("%w: Invalid", ErrSafe)
var ErrInvalidUser = fmt.Errorf("%w: Invalid user", ErrInvalid)
var ErrInvalidRepo = fmt.Errorf("%w: Invalid repo", ErrInvalid)
var ErrInvalidPat = fmt.Errorf("%w: Invalid pat", ErrInvalid)

var ErrExists = fmt.Errorf("%w: Exists", ErrSafe)
var ErrUserExists = fmt.Errorf("%w: User exists", ErrExists)
var ErrRepoExists = fmt.Errorf("%w: Repo exists", ErrExists)
var ErrPatExists = fmt.Errorf("%w: PAT exists", ErrExists)
