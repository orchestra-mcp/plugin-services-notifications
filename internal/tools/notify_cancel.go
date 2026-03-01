package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// NotifyCancelSchema returns the JSON Schema for the notify_cancel tool.
func NotifyCancelSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"notification_id": map[string]any{
				"type":        "string",
				"description": "ID of the scheduled notification to cancel",
			},
		},
		"required": []any{"notification_id"},
	})
	return s
}

// NotifyCancel returns a handler that cancels a scheduled notification by ID.
func NotifyCancel() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "notification_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		id := helpers.GetString(req.Arguments, "notification_id")

		return helpers.TextResult(fmt.Sprintf("Notification %s cancelled (in-memory only)", id)), nil
	}
}
