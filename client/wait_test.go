package client

import (
	"log"
	"time"

	"golang.org/x/net/context"
)

func ExampleClient_ContainerWait_withTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _ := NewEnvClient()
	_, err := client.ContainerWait(ctx, "container_id")
	if err != nil {
		log.Fatal(err)
	}
}
