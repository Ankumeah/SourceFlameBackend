package git

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
)

func UploadPack(ctx context.Context, repoId uint64, reader io.Reader, writer io.Writer) error {
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "git-upload-pack", "--stateless-rpc", realPath(repoId))
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error while uploading pack: %v\n%v", err.Error(), stderr.String())
	}

	return err
}
