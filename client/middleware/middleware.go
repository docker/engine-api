package middleware

import "github.com/docker/engine-api/client/transport"

// Middleware defines a function interface that the client
// middlewares must implement.
type Middleware func(sender transport.Sender) transport.Sender
