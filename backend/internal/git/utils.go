package git

import (
	"fmt"
	"path"
)

func real_path(repo_id uint64) string {
	return path.Join(base_path, fmt.Sprintf("%v", repo_id))
}
