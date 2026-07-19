package git

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
)

func ReceivePack(ctx context.Context, repoId uint64, reader io.Reader, writer io.Writer) error {
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "git-receive-pack", "--stateless-rpc", realPath(repoId))
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error while receiving pack: %v\n%v", err.Error(), stderr.String())
	}

	return err
}
