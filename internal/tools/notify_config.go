package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/globaldb"
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
			"ai_push_enabled": map[string]any{
				"type":        "boolean",
				"description": "Whether AI agent push notifications are enabled. Only used with action=set.",
			},
			"ai_voice_enabled": map[string]any{
				"type":        "boolean",
				"description": "Whether AI agent voice (TTS) alerts are enabled. Only used with action=set.",
			},
			"socket_push_enabled": map[string]any{
				"type":        "boolean",
				"description": "Whether to push notifications to connected desktop/mobile apps via TCP EventBus. Only used with action=set.",
			},
			"voice_name": map[string]any{
				"type":        "string",
				"description": "TTS voice name (e.g. 'Samantha'). Only used with action=set.",
			},
			"voice_speed": map[string]any{
				"type":        "string",
				"description": "TTS speed in words per minute. Only used with action=set.",
			},
			"voice_volume": map[string]any{
				"type":        "string",
				"description": "TTS volume from 0.0 to 1.0. Only used with action=set.",
			},
			"event_overrides": map[string]any{
				"type":        "object",
				"description": `Per-event-type overrides as a map: {"Notification": {"push": true, "voice": false}}. Only used with action=set.`,
			},
		},
		"required": []any{"action"},
	})
	return s
}

// configBool reads a globaldb config key and returns its boolean value.
// If the key is not set, returns the provided default.
func configBool(key string, defaultVal bool) bool {
	v := globaldb.GetConfig(key)
	if v == "" {
		return defaultVal
	}
	return v == "true"
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
			return handleGetConfig()
		case "set":
			return handleSetConfig(req.Arguments)
		default:
			return helpers.ErrorResult("validation_error", fmt.Sprintf("unknown action %q: must be 'get' or 'set'", action)), nil
		}
	}
}

func handleGetConfig() (*pluginv1.ToolResponse, error) {
	aiPush := configBool("notify.ai_push_enabled", true)
	aiVoice := configBool("notify.ai_voice_enabled", true)
	socketPush := configBool("notify.socket_push_enabled", true)
	qhs := globaldb.GetConfig("notify.quiet_hours_start")
	qhe := globaldb.GetConfig("notify.quiet_hours_end")
	dc := globaldb.GetConfig("notify.default_channel")
	voiceName := globaldb.GetConfig("notify.voice_name")
	voiceSpeed := globaldb.GetConfig("notify.voice_speed")
	voiceVolume := globaldb.GetConfig("notify.voice_volume")

	if dc == "" {
		dc = "system"
	}

	resultMap := map[string]any{
		"enabled":              true,
		"default_channel":      dc,
		"quiet_hours_start":    orDefault(qhs, "(not set)"),
		"quiet_hours_end":      orDefault(qhe, "(not set)"),
		"ai_push_enabled":      aiPush,
		"ai_voice_enabled":     aiVoice,
		"socket_push_enabled":  socketPush,
		"voice_name":           orDefault(voiceName, "(not set)"),
		"voice_speed":          orDefault(voiceSpeed, "(not set)"),
		"voice_volume":         orDefault(voiceVolume, "(not set)"),
	}

	// Collect per-event-type overrides.
	eventOverrides := globaldb.GetConfigPrefix("notify.event.")
	if len(eventOverrides) > 0 {
		overrides := make(map[string]map[string]string)
		for key, val := range eventOverrides {
			// Key format: notify.event.<Type>.<push|voice>
			rest := strings.TrimPrefix(key, "notify.event.")
			parts := strings.SplitN(rest, ".", 2)
			if len(parts) != 2 {
				continue
			}
			eventType, field := parts[0], parts[1]
			if overrides[eventType] == nil {
				overrides[eventType] = make(map[string]string)
			}
			overrides[eventType][field] = val
		}
		resultMap["event_overrides"] = overrides
	}

	result, err := helpers.JSONResult(resultMap)
	if err != nil {
		return helpers.TextResult(fmt.Sprintf("%v", resultMap)), nil
	}
	return result, nil
}

func handleSetConfig(args *structpb.Struct) (*pluginv1.ToolResponse, error) {
	var parts []string

	// String fields.
	for _, sf := range []struct{ param, key string }{
		{"quiet_hours_start", "notify.quiet_hours_start"},
		{"quiet_hours_end", "notify.quiet_hours_end"},
		{"default_channel", "notify.default_channel"},
		{"voice_name", "notify.voice_name"},
		{"voice_speed", "notify.voice_speed"},
		{"voice_volume", "notify.voice_volume"},
	} {
		if v := helpers.GetString(args, sf.param); v != "" {
			_ = globaldb.SetConfig(sf.key, v)
			parts = append(parts, fmt.Sprintf("%s=%s", sf.param, v))
		}
	}

	// Boolean fields.
	if args != nil {
		for _, bf := range []struct{ param, key string }{
			{"ai_push_enabled", "notify.ai_push_enabled"},
			{"ai_voice_enabled", "notify.ai_voice_enabled"},
			{"socket_push_enabled", "notify.socket_push_enabled"},
			{"enabled", "notify.enabled"},
		} {
			if _, ok := args.Fields[bf.param]; ok {
				bv := helpers.GetBool(args, bf.param)
				_ = globaldb.SetConfig(bf.key, fmt.Sprintf("%v", bv))
				parts = append(parts, fmt.Sprintf("%s=%v", bf.param, bv))
			}
		}

		// Per-event-type overrides.
		if eo, ok := args.Fields["event_overrides"]; ok && eo != nil {
			if sv, ok := eo.Kind.(*structpb.Value_StructValue); ok && sv.StructValue != nil {
				for eventType, val := range sv.StructValue.Fields {
					innerStruct, ok := val.Kind.(*structpb.Value_StructValue)
					if !ok || innerStruct.StructValue == nil {
						continue
					}
					for field, fv := range innerStruct.StructValue.Fields {
						key := fmt.Sprintf("notify.event.%s.%s", eventType, field)
						var strVal string
						switch v := fv.Kind.(type) {
						case *structpb.Value_BoolValue:
							strVal = fmt.Sprintf("%v", v.BoolValue)
						case *structpb.Value_StringValue:
							strVal = v.StringValue
						default:
							b, _ := json.Marshal(fv.AsInterface())
							strVal = string(b)
						}
						_ = globaldb.SetConfig(key, strVal)
						parts = append(parts, fmt.Sprintf("event.%s.%s=%s", eventType, field, strVal))
					}
				}
			}
		}
	}

	if len(parts) == 0 {
		return helpers.ErrorResult("validation_error", "at least one config field must be provided for action=set"), nil
	}
	return helpers.TextResult(fmt.Sprintf("Config updated: %s", strings.Join(parts, ", "))), nil
}

func orDefault(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
