package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-services-notifications/internal/notify"
	"google.golang.org/protobuf/types/known/structpb"
)

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
func NotifySend() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "title", "body"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		title := helpers.GetString(req.Arguments, "title")
		body := helpers.GetString(req.Arguments, "body")

		_, err := notify.Send(ctx, title, body)
		if err != nil {
			return helpers.ErrorResult("notify_error", fmt.Sprintf("Failed to send notification: %v", err)), nil
		}

		return helpers.TextResult(fmt.Sprintf("Notification sent: %s", title)), nil
	}
}
