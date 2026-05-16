package database

import (
  "github.com/Ankumeah/SourceFlameBackend/internal/hash"

  "context"
)

type user_driver interface {
  Add_User(ctx context.Context, username string, password_hash *hash.Hash) (uint64, error)
  Get_Id(ctx context.Context, username string) (uint64, error)
  Delete_User(ctx context.Context, user_id uint64) error
  Get_Hash(ctx context.Context, user_id uint64) (*hash.Hash, error)
  Info(ctx context.Context, user_id uint64) (*User_Info, error)
}

type git_driver interface {
  Create_Repo(ctx context.Context, owner_id uint64, repo_name string, private bool) (uint64, error)
  Get_Id(ctx context.Context, owner_id uint64, repo_name string) (uint64, error)
  Delete_Repo(ctx context.Context, repo_id uint64) error
  Get_Repos(ctx context.Context, user_id uint64, all bool, limit uint8, offset uint64) ([]string, error)
  Info(ctx context.Context, repo_id uint64) (*Repo_Info, error)
}

type pat_driver interface {
  Add_PAT(ctx context.Context, owner_id uint64, hash string, pat_name string) (uint64, error)
  Validate_PAT(ctx context.Context, owner_id uint64, hash string) (uint64, error)
  Get_Id(ctx context.Context, owner_id uint64, pat_name string) (uint64, error)
  Delete_PAT(ctx context.Context, pat_id uint64) error
  Get_PATs(ctx context.Context, user_id uint64) ([]string, error)
  Info(ctx context.Context, pat_id uint64) (*PAT_Info, error)
}

type User_db struct { db user_driver }
type Git_db struct { db git_driver }
type PAT_db struct { db pat_driver }

type User_Info struct {
  Creation uint64 `json:"creation"`
}

type Repo_Info struct {
  Creation uint64 `json:"creation"`
  Stars uint64 `json:"stars"`
  Private bool `json:"private"`
  Owner string `json:"owner"`
}

type PAT_Info struct {
  Name string `json:"name"`
  Creation uint64 `json:"creation"`
  Last_Used *uint64 `json:"last_used"`
}

type Connection_Config struct {
  Username string
  Password string
  Hostname string
  Port string
  Db_name string
  Db_config string
}
