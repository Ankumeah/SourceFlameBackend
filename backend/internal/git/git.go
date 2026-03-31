package git

import (
  "github.com/go-git/go-git/v6"

  "os"
  "io"
  "os/exec"
  "fmt"
  "context"
  "log"
  "bytes"
)

const base_path = "/srv/git/"

func real_path(repo_id uint64) string {
  return fmt.Sprintf("%v%v", base_path, repo_id)
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
    log.Printf("Error while unloading pack: %v\n%v", err.Error(), string(stderr.Bytes()))
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
