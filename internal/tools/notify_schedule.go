package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// NotifyScheduleSchema returns the JSON Schema for the notify_schedule tool.
func NotifyScheduleSchema() *structpb.Struct {
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
			"at": map[string]any{
				"type":        "string",
				"description": "ISO 8601 datetime string for when to send the notification",
			},
			"sound": map[string]any{
				"type":        "string",
				"description": "Optional sound name to play with the notification",
			},
		},
		"required": []any{"title", "body", "at"},
	})
	return s
}

// NotifySchedule returns a handler that acknowledges a scheduled notification request.
func NotifySchedule() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "title", "body", "at"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		title := helpers.GetString(req.Arguments, "title")
		at := helpers.GetString(req.Arguments, "at")

		return helpers.TextResult(
			fmt.Sprintf(
				"Scheduled notification '%s' for %s. Note: persistence across restarts not supported.",
				title, at,
			),
		), nil
	}
}
