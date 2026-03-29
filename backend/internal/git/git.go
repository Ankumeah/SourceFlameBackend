package git

import (
  "github.com/go-git/go-git/v6"

  "os"
  "io"
  "os/exec"
  "fmt"
  "context"
  "time"
)

const base_path = "/srv/git/"

func real_path(repo_id uint64) string {
  return fmt.Sprintf("%v%v", base_path, repo_id)
}

func Create_Repo(repo_id uint64, private bool) error {
  _, err := git.PlainInit(real_path(repo_id), true)
  if err == git.ErrTargetDirNotEmpty {
    return Error_Repository_Exists
  }

  return err
}

func Delete_Repo(repo_id uint64) error {
  return os.RemoveAll(real_path(repo_id))
}

func Info_Refs(c context.Context, repo_id uint64, service string, writer io.Writer) error {
  ctx, cancel := context.WithTimeout(c, time.Second * 10)
  defer cancel()

	if service != "git-receive-pack" && service != "git-upload-pack" {
		return Error_Unsupported_Service
	}

	pktLine := fmt.Sprintf("# service=%s\n", service)
	fmt.Fprintf(writer, "%04x%s0000", len(pktLine)+4, pktLine)

	cmd := exec.CommandContext(ctx, service, "--stateless-rpc", "--advertise-refs", real_path(repo_id))
	cmd.Stdout = writer
	cmd.Stderr = nil

  err := cmd.Run()
  if err == context.DeadlineExceeded {
    return Error_Timeout
  } else {
    return err
  }
}

func Receive_Pack(c context.Context, repo_id uint64, reader io.Reader, writer io.Writer) error {
  ctx, cancel := context.WithTimeout(c, time.Second * 60)
  defer cancel()

	cmd := exec.CommandContext(ctx, "git-receive-pack", "--stateless-rpc", real_path(repo_id))
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = nil

  err := cmd.Run()
  if err == context.DeadlineExceeded {
    return Error_Timeout
  } else {
  return err
  }
}
