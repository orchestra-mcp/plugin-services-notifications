package servicesnotifications

import (
	"github.com/orchestra-mcp/plugin-services-notifications/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all notification tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	np := &internal.NotificationsPlugin{}
	np.RegisterTools(builder)
}
