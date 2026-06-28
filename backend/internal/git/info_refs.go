package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
)

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
