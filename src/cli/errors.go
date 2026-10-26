package cli

import (
	"errors"
	"strings"
)

type AgentErrorRecovery struct {
	Action  string `json:"action"`
	Command string `json:"command,omitempty"`
	Hint    string `json:"hint"`
}
type AgentErrorEnvelope struct {
	SchemaVersion string             `json:"schemaVersion"`
	Message       string             `json:"message"`
	ErrorCode     string             `json:"error_code,omitempty"`
	StatusCode    int                `json:"status_code,omitempty"`
	RequestId     string             `json:"request_id,omitempty"`
	Recovery      AgentErrorRecovery `json:"recovery"`
}
type AgentError struct {
	envelope AgentErrorEnvelope
	exitCode int
	cause    error
}

func NewAgentError(envelope AgentErrorEnvelope, cause error) *AgentError {
	if strings.TrimSpace(envelope.Message) == "" || strings.TrimSpace(envelope.Recovery.Action) == "" || strings.TrimSpace(envelope.Recovery.Hint) == "" {
		return nil
	}
	return &AgentError{envelope: envelope, cause: cause}
}
func (e *AgentError) WithExitCode(code int) *AgentError {
	copy := *e
	copy.exitCode = code
	return &copy
}
func (e *AgentError) WithSchemaVersion() *AgentError {
	copy := *e
	copy.envelope.SchemaVersion = "v1"
	return &copy
}
func (e *AgentError) Unwrap() error                { return e.cause }
func (e *AgentError) Error() string                { return e.envelope.Message }
func (e *AgentError) Envelope() AgentErrorEnvelope { return e.envelope }
func (e *AgentError) ExitCode() int {
	if e.exitCode != 0 {
		return e.exitCode
	}
	return 2
}
func UsageError(code, message string) *AgentError {
	return NewAgentError(AgentErrorEnvelope{Message: message, ErrorCode: code, Recovery: AgentErrorRecovery{Action: "help", Command: "inspire --help --cli-output json", Hint: "Inspect command help and correct the input."}}, nil).WithSchemaVersion().WithExitCode(2)
}
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var target *AgentError
	if errors.As(err, &target) {
		return target.ExitCode()
	}
	return 1
}
func ErrorEnvelope(err error) AgentErrorEnvelope {
	var target *AgentError
	if errors.As(err, &target) {
		return target.Envelope()
	}
	return AgentErrorEnvelope{SchemaVersion: "v1", Message: "Command failed.", ErrorCode: "ExecutionError", Recovery: AgentErrorRecovery{Action: "inspect", Hint: "Check the command and execution environment."}}
}
