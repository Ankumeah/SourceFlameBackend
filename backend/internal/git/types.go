package git

import "path"

const max_blob_size = 2097152

var base_path string

func init() {
	// This breaks on windows,
	// Im sorry if you are running this on bare
	// metal on windows but if you are then you are
	// kinda asking for it.
	// The fact that you are reading this means you
	// can probably change this on your own.
	// Good luck and apologies from my side
	base_path = path.Join("/srv", "git")
}

type File struct {
	File_name string `json:"file_name" binding:"required"`
	Dir       bool   `json:"dir" binding:"required"`
}

type author struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
type Commit struct {
	Author    author `json:"author" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
	Timestamp uint64 `json:"timestamp" binding:"required"`
}
