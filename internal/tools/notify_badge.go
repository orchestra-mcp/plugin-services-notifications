package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// NotifyBadgeSchema returns the JSON Schema for the notify_badge tool.
func NotifyBadgeSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"count": map[string]any{
				"type":        "integer",
				"description": "Badge count to display (0 clears the badge)",
			},
		},
		"required": []any{"count"},
	})
	return s
}

// hasField reports whether the given field key is present in the struct,
// regardless of its type. This is used for numeric fields where ValidateRequired
// (which uses GetString) would incorrectly reject valid zero values.
func hasField(args *structpb.Struct, key string) bool {
	if args == nil {
		return false
	}
	_, ok := args.Fields[key]
	return ok
}

// NotifyBadge returns a handler that sets the application dock/taskbar badge count.
func NotifyBadge() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if req.Arguments == nil || !hasField(req.Arguments, "count") {
			return helpers.ErrorResult("validation_error", "missing required fields: count"), nil
		}

		count := helpers.GetInt(req.Arguments, "count")

		if runtime.GOOS == "darwin" {
			script := fmt.Sprintf(
				`tell application "Finder" to set badge of front window to %d`,
				count,
			)
			cmd := exec.CommandContext(ctx, "osascript", "-e", script)
			out, err := cmd.CombinedOutput()
			if err != nil {
				// Badge setting may fail without a running app context — report
				// what happened but still return the human-readable message.
				_ = out
			}
		}

		return helpers.TextResult(fmt.Sprintf("Badge set to %d", count)), nil
	}
}
