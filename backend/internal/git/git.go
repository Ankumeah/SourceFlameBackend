package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

const max_blob_size = 2097152

var base_path string

func init() {
  base_path = path.Join("srv", "git")
}

func real_path(repo_id uint64) string {
  return path.Join(base_path, fmt.Sprintf("%v", repo_id))
}

func Create_Repo(repo_id uint64, private bool) error {
  _, err := git.PlainInit(real_path(repo_id), true)
  if err == git.ErrTargetDirNotEmpty {
    return Error_Repository_Exists
  } else if err != nil {
    return err
  }

  return nil
}

func Delete_Repo(repo_id uint64) error {
  return os.RemoveAll(real_path(repo_id))
}

func Info_Refs(ctx context.Context, repo_id uint64, service string, writer io.Writer) error {
  var buffer bytes.Buffer
  var stderr bytes.Buffer

  pktLine := "# service=" + service + "\n"
  fmt.Fprintf(&buffer, "%04x%v0000", len(pktLine)+4, pktLine)

  cmd := exec.CommandContext(ctx, service, "--stateless-rpc", "--advertise-refs", real_path(repo_id))
  cmd.Stdout = &buffer
  cmd.Stderr = &stderr

  err := cmd.Run()
  if err != nil {
    log.Printf("Error while providing info refs: %v\n%v", err.Error(), string(stderr.Bytes()))
    return err
  }

  writer.Write(buffer.Bytes())
  return nil
}

func Upload_Pack(ctx context.Context, repo_id uint64, reader io.Reader, writer io.Writer) error {
  var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "git-upload-pack", "--stateless-rpc", real_path(repo_id))
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = &stderr

  err := cmd.Run()
  if err != nil {
    log.Printf("Error while uploading pack: %v\n%v", err.Error(), string(stderr.Bytes()))
  }

  return err
}

func Receive_Pack(ctx context.Context, repo_id uint64, reader io.Reader, writer io.Writer) error {
  var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "git-receive-pack", "--stateless-rpc", real_path(repo_id))
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = &stderr

  err := cmd.Run()
  if err != nil {
    log.Printf("Error while receiveing pack: %v\n%v", err.Error(), string(stderr.Bytes()))
  }

  return err
}

func Get_Glob(repo_id uint64, commit_hash string, path string) (string, error) {
  repo, err := git.PlainOpen(real_path(repo_id))
  if err != nil {
    return "", err
  }

  hash := plumbing.NewHash(commit_hash)
  commit, err := repo.CommitObject(hash)
  if err == plumbing.ErrObjectNotFound {
    return "", Error_Inavlid_Commit_Hash
  } else if err != nil {
    log.Printf("Error while getting commit object: %v\n", err.Error())
    return "", err
  }

  tree, err := commit.Tree()
  if err != nil {
    log.Printf("Error while getting commit worktree: %v\n", err.Error())
    return "", err
  }

  file, err := tree.File(path)
  if err == object.ErrFileNotFound || err == object.ErrDirectoryNotFound {
    return "", Error_Blob_Not_Found
  } else if err != nil {
    log.Printf("Error while opening blob file: %v\n", err.Error())
    return "", err
  }

  if file.Size >= max_blob_size {
    return "", Error_Blob_Too_Large
  }

  return file.Contents()
}
