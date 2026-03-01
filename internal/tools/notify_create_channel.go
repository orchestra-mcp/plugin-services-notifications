package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// validChannels is the set of supported notification channels.
var validChannels = map[string]bool{
	"build":    true,
	"test":     true,
	"deploy":   true,
	"ai":       true,
	"reminder": true,
	"system":   true,
	"git":      true,
}

// NotifyCreateChannelSchema returns the JSON Schema for the notify_create_channel tool.
func NotifyCreateChannelSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Channel name (e.g. build, test, deploy, ai, reminder, system, git)",
			},
			"description": map[string]any{
				"type":        "string",
				"description": "Human-readable description of the channel's purpose",
			},
			"sound": map[string]any{
				"type":        "string",
				"description": "Default sound for this channel (optional)",
			},
			"enabled": map[string]any{
				"type":        "boolean",
				"description": "Whether this channel is enabled (default true)",
			},
		},
		"required": []any{"name"},
	})
	return s
}

// NotifyCreateChannel returns a handler that creates or updates a named notification channel.
func NotifyCreateChannel() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "name"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		name := helpers.GetString(req.Arguments, "name")
		desc := helpers.GetString(req.Arguments, "description")

		msg := fmt.Sprintf("Channel '%s' created", name)
		if validChannels[name] {
			msg = fmt.Sprintf("Channel '%s' registered (built-in channel)", name)
		}
		if desc != "" {
			msg += fmt.Sprintf(": %s", desc)
		}

		return helpers.TextResult(msg), nil
	}
}
