package git

import (
  "context"
  "bytes"
  "os/exec"
  "io"
  "log"
)

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
