package git

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
)

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
