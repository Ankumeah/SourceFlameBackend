package git

import (
	"fmt"
	"path"
)

func realPath(repoId uint64) string {
	return path.Join(basePath, fmt.Sprintf("%v", repoId))
}
