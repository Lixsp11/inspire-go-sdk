package cli_test

import (
	"encoding/json"
	"errors"
	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"strings"
	"testing"
)

func TestAgentErrorContract(t *testing.T) {
	cause := errors.New("private diagnostic value")
	err := cli.NewAgentError(cli.AgentErrorEnvelope{Message: "Safe message.", ErrorCode: "Test", Recovery: cli.AgentErrorRecovery{Action: "inspect", Hint: "Inspect the input."}}, cause).WithSchemaVersion().WithExitCode(1)
	data, e := json.Marshal(err.Envelope())
	if e != nil || !errors.Is(err, cause) || cli.ExitCode(err) != 1 || strings.Contains(string(data), cause.Error()) || !strings.Contains(string(data), `"schemaVersion":"v1"`) {
		t.Fatalf("invalid agent error: %s %v", data, e)
	}
	if cli.NewAgentError(cli.AgentErrorEnvelope{}, cause) != nil {
		t.Fatal("incomplete recovery envelope was accepted")
	}
}
