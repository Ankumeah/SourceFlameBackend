package git

import (
  "path"
)

const max_blob_size = 2097152

var base_path string

func init() {
  base_path = path.Join("srv", "git")
}

type File struct {
  File_name string `json:"file_name" binding:"required"`
  Dir bool `json:"dir" binding:"required"`
}
