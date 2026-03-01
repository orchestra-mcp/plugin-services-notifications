package notify

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Send delivers a desktop notification with the given title and body.
// On macOS it uses AppleScript via `osascript`; on all other platforms it uses
// `notify-send`.
func Send(ctx context.Context, title, body string) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, body, title)
		cmd = exec.CommandContext(ctx, "osascript", "-e", script)
	default:
		cmd = exec.CommandContext(ctx, "notify-send", title, body)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
