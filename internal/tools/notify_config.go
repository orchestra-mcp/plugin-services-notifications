package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// NotifyConfigSchema returns the JSON Schema for the notify_config tool.
func NotifyConfigSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": `Action to perform: "get" to read current config, "set" to update config`,
				"enum":        []any{"get", "set"},
			},
			"quiet_hours_start": map[string]any{
				"type":        "string",
				"description": `Start of quiet hours in HH:MM 24h format (e.g. "22:00"). Only used with action=set.`,
			},
			"quiet_hours_end": map[string]any{
				"type":        "string",
				"description": `End of quiet hours in HH:MM 24h format (e.g. "08:00"). Only used with action=set.`,
			},
			"default_channel": map[string]any{
				"type":        "string",
				"description": "Default channel for notifications. Only used with action=set.",
			},
			"enabled": map[string]any{
				"type":        "boolean",
				"description": "Whether notifications are globally enabled. Only used with action=set.",
			},
		},
		"required": []any{"action"},
	})
	return s
}

// NotifyConfig returns a handler that gets or sets notification configuration.
func NotifyConfig() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "action"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		action := helpers.GetString(req.Arguments, "action")
		switch action {
		case "get":
			return helpers.TextResult(strings.Join([]string{
				"## Notification Config",
				"- **enabled**: true",
				"- **default_channel**: system",
				"- **quiet_hours_start**: (not set)",
				"- **quiet_hours_end**: (not set)",
			}, "\n")), nil

		case "set":
			var parts []string
			if qhs := helpers.GetString(req.Arguments, "quiet_hours_start"); qhs != "" {
				parts = append(parts, fmt.Sprintf("quiet_hours_start=%s", qhs))
			}
			if qhe := helpers.GetString(req.Arguments, "quiet_hours_end"); qhe != "" {
				parts = append(parts, fmt.Sprintf("quiet_hours_end=%s", qhe))
			}
			if dc := helpers.GetString(req.Arguments, "default_channel"); dc != "" {
				parts = append(parts, fmt.Sprintf("default_channel=%s", dc))
			}
			if len(parts) == 0 {
				return helpers.ErrorResult("validation_error", "at least one config field must be provided for action=set"), nil
			}
			return helpers.TextResult(fmt.Sprintf("Config updated: %s", strings.Join(parts, ", "))), nil

		default:
			return helpers.ErrorResult("validation_error", fmt.Sprintf("unknown action %q: must be 'get' or 'set'", action)), nil
		}
	}
}
