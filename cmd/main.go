// Command services-notifications is the entry point for the services.notifications
// plugin binary. It provides 6 MCP tools for sending and managing notifications.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/plugin-services-notifications/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

func main() {
	builder := plugin.New("services.notifications").
		Version("0.1.0").
		Description("Desktop notification services for macOS and Linux").
		Author("Orchestra").
		Binary("services-notifications")

	tp := &internal.NotificationsPlugin{}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("services.notifications: %v", err)
	}
}
