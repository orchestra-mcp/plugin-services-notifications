package notify

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//go:embed notification-icon.png
var notificationIcon []byte

const bundleID = "com.orchestramcp.orchestra"

var iconOnce sync.Once

// ensureIcon writes the embedded notification icon to ~/.orchestra/notification-icon.png
// if it doesn't already exist. Returns the path to the icon file.
func ensureIcon() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	iconPath := filepath.Join(home, ".orchestra", "notification-icon.png")
	iconOnce.Do(func() {
		if _, err := os.Stat(iconPath); err == nil {
			return // Already exists.
		}
		_ = os.MkdirAll(filepath.Dir(iconPath), 0755)
		_ = os.WriteFile(iconPath, notificationIcon, 0644)
	})
	return iconPath
}

// Send delivers a desktop notification with the Orchestra icon.
// On macOS, clicking the notification opens the Orchestra app.
// On Linux, the Orchestra icon is shown via notify-send.
func Send(ctx context.Context, title, body string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return sendDarwin(ctx, title, body)
	default:
		return sendLinux(ctx, title, body)
	}
}

// sendDarwin sends a notification via AppleScript.
// Uses "tell application id" so the notification shows the Orchestra icon
// and clicking it opens the app. Falls back to plain osascript if the app isn't installed.
func sendDarwin(ctx context.Context, title, body string) (string, error) {
	// Primary: use Orchestra app as notification source (shows app icon, click opens app).
	script := fmt.Sprintf(`
try
	tell application id %q
		display notification %q with title %q
	end tell
on error
	display notification %q with title %q
end try`, bundleID, body, title, body, title)

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// sendLinux sends a notification via notify-send with the Orchestra icon.
func sendLinux(ctx context.Context, title, body string) (string, error) {
	iconPath := ensureIcon()
	args := []string{"--app-name=Orchestra"}
	if iconPath != "" {
		args = append(args, "--icon="+iconPath)
	}
	args = append(args, title, body)

	cmd := exec.CommandContext(ctx, "notify-send", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
