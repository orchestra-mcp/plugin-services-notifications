package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// NotifyHistorySchema returns the JSON Schema for the notify_history tool.
func NotifyHistorySchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"limit": map[string]any{
				"type":        "integer",
				"description": "Maximum number of history entries to return (default 20)",
			},
			"channel": map[string]any{
				"type":        "string",
				"description": "Filter history by channel name",
			},
		},
	})
	return s
}

// NotifyHistory returns a handler that retrieves recent notification history.
// In the current implementation, notifications are not persisted across process
// restarts, so this always returns an empty history with an informational note.
func NotifyHistory() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		limit := helpers.GetInt(req.Arguments, "limit")
		if limit <= 0 {
			limit = 20
		}
		channel := helpers.GetString(req.Arguments, "channel")

		msg := fmt.Sprintf("Notification history (limit=%d", limit)
		if channel != "" {
			msg += fmt.Sprintf(", channel=%s", channel)
		}
		msg += "): no history available (in-memory only, not persisted across restarts)"

		return helpers.TextResult(msg), nil
	}
}
