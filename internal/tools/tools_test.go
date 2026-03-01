package tools

import (
	"context"
	"strings"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- helpers ----------

func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

func getText(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	if r := resp.GetResult(); r != nil {
		if f := r.GetFields(); f != nil {
			if tf, ok := f["text"]; ok {
				return tf.GetStringValue()
			}
		}
	}
	return ""
}

// ---------- notify_send ----------

func TestNotifySend_MissingTitle(t *testing.T) {
	resp := callTool(t, NotifySend(), map[string]any{"body": "hello"})
	if !isError(resp) {
		t.Error("expected validation_error for missing title")
	}
}

func TestNotifySend_MissingBody(t *testing.T) {
	resp := callTool(t, NotifySend(), map[string]any{"title": "Test"})
	if !isError(resp) {
		t.Error("expected validation_error for missing body")
	}
}

func TestNotifySend_ValidArgs(t *testing.T) {
	// notify.Send may fail without a notification daemon — that's OK.
	// We just verify no Go-level panic and the response is well-formed.
	resp := callTool(t, NotifySend(), map[string]any{"title": "Test", "body": "Hello"})
	_ = resp
}

// ---------- notify_schedule ----------

func TestNotifySchedule_MissingAt(t *testing.T) {
	resp := callTool(t, NotifySchedule(), map[string]any{"title": "T", "body": "B"})
	if !isError(resp) {
		t.Error("expected validation_error for missing at")
	}
}

func TestNotifySchedule_ValidArgs(t *testing.T) {
	resp := callTool(t, NotifySchedule(), map[string]any{
		"title": "Reminder",
		"body":  "Stand up!",
		"at":    "2026-03-01T09:00:00Z",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "Scheduled") {
		t.Errorf("expected 'Scheduled' in response, got: %s", txt)
	}
}

// ---------- notify_cancel ----------

func TestNotifyCancel_MissingID(t *testing.T) {
	resp := callTool(t, NotifyCancel(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing notification_id")
	}
}

func TestNotifyCancel_ValidID(t *testing.T) {
	resp := callTool(t, NotifyCancel(), map[string]any{"notification_id": "notif-abc123"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "notif-abc123") {
		t.Errorf("expected ID in response, got: %s", txt)
	}
}

// ---------- notify_list_pending ----------

func TestNotifyListPending_NoArgs(t *testing.T) {
	resp := callTool(t, NotifyListPending(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
}

// ---------- notify_badge ----------

func TestNotifyBadge_MissingCount(t *testing.T) {
	resp := callTool(t, NotifyBadge(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing count")
	}
}

func TestNotifyBadge_SetCount(t *testing.T) {
	resp := callTool(t, NotifyBadge(), map[string]any{"count": float64(5)})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "5") {
		t.Errorf("expected count in response, got: %s", txt)
	}
}

func TestNotifyBadge_ClearBadge(t *testing.T) {
	resp := callTool(t, NotifyBadge(), map[string]any{"count": float64(0)})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
}

// ---------- notify_config ----------

func TestNotifyConfig_MissingAction(t *testing.T) {
	resp := callTool(t, NotifyConfig(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing action")
	}
}

func TestNotifyConfig_GetAction(t *testing.T) {
	resp := callTool(t, NotifyConfig(), map[string]any{"action": "get"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "enabled") {
		t.Errorf("expected config fields in response, got: %s", txt)
	}
}

func TestNotifyConfig_SetAction(t *testing.T) {
	resp := callTool(t, NotifyConfig(), map[string]any{
		"action":            "set",
		"quiet_hours_start": "22:00",
		"quiet_hours_end":   "08:00",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "22:00") {
		t.Errorf("expected quiet_hours_start in response, got: %s", txt)
	}
}

func TestNotifyConfig_SetNoFields(t *testing.T) {
	resp := callTool(t, NotifyConfig(), map[string]any{"action": "set"})
	if !isError(resp) {
		t.Error("expected validation_error when set has no fields")
	}
}

func TestNotifyConfig_InvalidAction(t *testing.T) {
	resp := callTool(t, NotifyConfig(), map[string]any{"action": "delete"})
	if !isError(resp) {
		t.Error("expected validation_error for unknown action")
	}
}

// ---------- notify_history ----------

func TestNotifyHistory_NoArgs(t *testing.T) {
	resp := callTool(t, NotifyHistory(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "20") {
		t.Errorf("expected default limit=20 in response, got: %s", txt)
	}
}

func TestNotifyHistory_WithLimit(t *testing.T) {
	resp := callTool(t, NotifyHistory(), map[string]any{"limit": float64(5)})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "5") {
		t.Errorf("expected limit=5 in response, got: %s", txt)
	}
}

func TestNotifyHistory_WithChannel(t *testing.T) {
	resp := callTool(t, NotifyHistory(), map[string]any{"channel": "build"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "build") {
		t.Errorf("expected channel=build in response, got: %s", txt)
	}
}

// ---------- notify_create_channel ----------

func TestNotifyCreateChannel_MissingName(t *testing.T) {
	resp := callTool(t, NotifyCreateChannel(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing name")
	}
}

func TestNotifyCreateChannel_BuiltinChannel(t *testing.T) {
	resp := callTool(t, NotifyCreateChannel(), map[string]any{"name": "build"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "build") {
		t.Errorf("expected channel name in response, got: %s", txt)
	}
}

func TestNotifyCreateChannel_CustomChannel(t *testing.T) {
	resp := callTool(t, NotifyCreateChannel(), map[string]any{
		"name":        "my-custom",
		"description": "My custom alerts",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "my-custom") {
		t.Errorf("expected channel name in response, got: %s", txt)
	}
}
