package servicesnotifications

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-services-notifications/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Sender is the interface for sending requests to the orchestrator.
// It is satisfied by the in-process router.
type Sender interface {
	Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error)
}

// Register adds all notification tools to the builder. The sender is used for
// socket push delivery to connected desktop/mobile apps via EventBus.
func Register(builder *plugin.PluginBuilder, sender Sender) {
	np := &internal.NotificationsPlugin{Sender: sender}
	np.RegisterTools(builder)
}
