package internal

import (
	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-services-notifications/internal/tools"
)

// NotificationsPlugin registers all notification tools with the plugin builder.
type NotificationsPlugin struct {
	Sender tools.Sender
}

// RegisterTools registers all 8 notification tools with the plugin builder.
func (np *NotificationsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("notify_send",
		"Send a desktop notification with a title and body. Also pushes to connected apps via socket if enabled.",
		tools.NotifySendSchema(), tools.NotifySend(np.Sender))

	builder.RegisterTool("notify_schedule",
		"Schedule a notification to be sent at a specific time",
		tools.NotifyScheduleSchema(), tools.NotifySchedule())

	builder.RegisterTool("notify_cancel",
		"Cancel a previously scheduled notification",
		tools.NotifyCancelSchema(), tools.NotifyCancel())

	builder.RegisterTool("notify_list_pending",
		"List all pending scheduled notifications",
		tools.NotifyListPendingSchema(), tools.NotifyListPending())

	builder.RegisterTool("notify_badge",
		"Set the application dock/taskbar badge count",
		tools.NotifyBadgeSchema(), tools.NotifyBadge())

	builder.RegisterTool("notify_config",
		"Get or set notification configuration (quiet hours, channels, voice, socket push, per-event toggles)",
		tools.NotifyConfigSchema(), tools.NotifyConfig())

	builder.RegisterTool("notify_history",
		"Retrieve recent notification history, optionally filtered by channel",
		tools.NotifyHistorySchema(), tools.NotifyHistory())

	builder.RegisterTool("notify_create_channel",
		"Create or register a named notification channel (build, test, deploy, ai, reminder, system, git)",
		tools.NotifyCreateChannelSchema(), tools.NotifyCreateChannel())
}
