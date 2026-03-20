package tools

import (
	"context"
	"fmt"
	"log/slog"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/globaldb"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-services-notifications/internal/notify"
	"google.golang.org/protobuf/types/known/structpb"
)

// Sender is the interface for sending requests to the orchestrator.
type Sender interface {
	Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error)
}

// NotifySendSchema returns the JSON Schema for the notify_send tool.
func NotifySendSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "Notification title",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Notification body text",
			},
			"sound": map[string]any{
				"type":        "string",
				"description": "Optional sound name to play with the notification",
			},
		},
		"required": []any{"title", "body"},
	})
	return s
}

// NotifySend returns a handler that sends a desktop notification.
// If sender is non-nil and socket push is enabled, it also publishes
// to the EventBus topic "notifications" for connected apps.
func NotifySend(sender Sender) func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "title", "body"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		title := helpers.GetString(req.Arguments, "title")
		body := helpers.GetString(req.Arguments, "body")
		sound := helpers.GetString(req.Arguments, "sound")

		// Send OS-level notification.
		_, err := notify.Send(ctx, title, body)
		if err != nil {
			return helpers.ErrorResult("notify_error", fmt.Sprintf("Failed to send notification: %v", err)), nil
		}

		// Also push to connected apps via EventBus if enabled.
		if sender != nil && globaldb.GetConfig("notify.socket_push_enabled") != "false" {
			payload, pErr := structpb.NewStruct(map[string]any{
				"type":  "notification",
				"title": title,
				"body":  body,
				"sound": sound,
			})
			if pErr == nil {
				publishReq := &pluginv1.PluginRequest{
					RequestId: helpers.NewUUID(),
					Request: &pluginv1.PluginRequest_Publish{
						Publish: &pluginv1.Publish{
							Topic:        "notifications",
							EventType:    "push",
							Payload:      payload,
							SourcePlugin: "services.notifications",
						},
					},
				}
				if _, sErr := sender.Send(ctx, publishReq); sErr != nil {
					slog.Warn("notify_send: failed to send socket push", "error", sErr)
				}
			}
		}

		return helpers.TextResult(fmt.Sprintf("Notification sent: %s", title)), nil
	}
}
