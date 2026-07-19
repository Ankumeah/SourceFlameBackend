package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

	"context"
)

type userDriver interface {
	AddUser(ctx context.Context, username string, passwordHash *hash.Hash) (uint64, error)
	GetId(ctx context.Context, username string) (uint64, error)
	DeleteUser(ctx context.Context, userId uint64) error
	GetHash(ctx context.Context, userId uint64) (*hash.Hash, error)
	Info(ctx context.Context, userId uint64) (*UserInfo, error)
}

type gitDriver interface {
	CreateRepo(ctx context.Context, ownerId uint64, repoName string) (uint64, error)
	GetId(ctx context.Context, ownerId uint64, repoName string) (uint64, error)
	DeleteRepo(ctx context.Context, repoId uint64) error
	GetRepos(ctx context.Context, userId uint64, limit uint8, offset uint64) ([]string, error)
	Info(ctx context.Context, repoId uint64) (*RepoInfo, error)
}

type patDriver interface {
	AddPAT(ctx context.Context, ownerId uint64, hash string, patName string) (uint64, error)
	ValidatePAT(ctx context.Context, ownerId uint64, hash string) (uint64, error)
	GetId(ctx context.Context, ownerId uint64, patName string) (uint64, error)
	DeletePAT(ctx context.Context, patId uint64) error
	GetPATs(ctx context.Context, userId uint64) ([]string, error)
	Info(ctx context.Context, patId uint64) (*PATInfo, error)
	UpdateUse(ctx context.Context, patId uint64) error
}

type UserDb struct{ db userDriver }
type GitDb struct{ db gitDriver }
type PATDb struct{ db patDriver }

type UserInfo struct {
	Creation uint64 `json:"creation"`
}

type RepoInfo struct {
	Creation uint64 `json:"creation"`
	Stars    uint64 `json:"stars"`
	Owner    string `json:"owner"`
}

type PATInfo struct {
	Name     string  `json:"name"`
	Creation uint64  `json:"creation"`
	LastUsed *uint64 `json:"last_used"`
}
