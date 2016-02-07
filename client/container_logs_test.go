package client

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/engine-api/types"

	"golang.org/x/net/context"
)

func ExampleClient_ContainerLogs_withTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _ := NewEnvClient()
	reader, err := client.ContainerLogs(ctx, types.ContainerLogsOptions{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
