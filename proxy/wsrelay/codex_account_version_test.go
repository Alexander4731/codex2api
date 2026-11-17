package wsrelay

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/codex2api/auth"
	"github.com/codex2api/proxy"
)

func TestPinnedAccountVersionFollowsSyncOnWebsocketHandshake(t *testing.T) {
	old := proxy.CurrentRuntimeSettings()
	t.Cleanup(func() { proxy.ApplyRuntimeSettings(old) })
	t.Setenv("CODEX_DISABLE_GENERATED_ATTESTATION", "true")
	t.Setenv("CODEX_WS_SEND_USER_AGENT", "true")
	ua := "codex-tui/{{codex_cli_version}} (Windows 10.0.19045; x86_64) unknown (codex-tui; {{codex_cli_version}})"
	account := &auth.Account{DBID: 1, AccountID: "workspace", CustomHeaders: map[string]string{"User-Agent": ua, "Version": proxy.CodexCLIVersionPlaceholder, "Originator": "codex-tui"}}
	for _, version := range []string{"0.169.0", "0.170.0"} {
		settings := old
		settings.ClientCompatMode, settings.CodexSyncedCLIVersion = proxy.ClientCompatModeForce, version
		proxy.ApplyRuntimeSettings(settings)
		headers, err := NewExecutor().prepareWebsocketHeadersChecked(context.Background(), "token", account, "workspace", "session", "key", nil, http.Header{"User-Agent": {"caller/1"}}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		if headers.Get("User-Agent") != strings.ReplaceAll(ua, proxy.CodexCLIVersionPlaceholder, version) || headers.Get("Version") != version || headers.Get("Originator") != "codex-tui" {
			t.Fatalf("wrong handshake identity: %v", headers)
		}
	}
	if account.GetCustomHeaders()["User-Agent"] != ua {
		t.Fatal("stored template mutated")
	}
}
