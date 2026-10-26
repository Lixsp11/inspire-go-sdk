package client

import (
	"errors"
	"fmt"
	"github.com/alibabacloud-go/tea/dara"
)

// SDKError implements the official Dara error interfaces with native diagnostics.
type SDKError struct {
	*dara.SDKError
	RequestId  *string
	Action     string
	Attempts   int
	CauseCode  string
	retryAfter *int64
	cause      error
}

var _ dara.BaseError = (*SDKError)(nil)
var _ dara.ResponseError = (*SDKError)(nil)

// The upstream SDKError embeds a BaseError interface without initializing its GetName delegate.
type baseErrorView struct{ sdk *dara.SDKError }

func (v baseErrorView) GetName() *string { return v.sdk.Name }
func (v baseErrorView) GetCode() *string { return v.sdk.Code }
func (v baseErrorView) Error() string {
	return dara.StringValue(v.sdk.Code) + ": " + dara.StringValue(v.sdk.Message)
}

func NewError(name, code, message string) *SDKError {
	base := dara.NewSDKError(map[string]interface{}{"code": code, "message": message})
	base.Name = dara.String(name)
	base.BaseError = baseErrorView{base}
	return &SDKError{SDKError: base}
}
func (e *SDKError) Error() string {
	suffix := ""
	if e.RequestId != nil {
		suffix = " [request_id=" + *e.RequestId + "]"
	}
	return dara.StringValue(e.Code) + ": " + dara.StringValue(e.Message) + suffix
}
func (e *SDKError) GetName() *string                              { return e.Name }
func (e *SDKError) GetRetryAfter() *int64                         { return e.retryAfter }
func (e *SDKError) GetStatusCode() *int                           { return e.StatusCode }
func (e *SDKError) GetAccessDeniedDetail() map[string]interface{} { return e.AccessDeniedDetail }
func (e *SDKError) GetDescription() *string                       { return e.Description }
func (e *SDKError) GetData() map[string]interface{} {
	return map[string]interface{}{"requestId": e.RequestId, "statusCode": e.StatusCode}
}
func (e *SDKError) Unwrap() error        { return e.SDKError }
func (e *SDKError) Is(target error) bool { return e.cause != nil && errors.Is(e.cause, target) }
func (e *SDKError) String() string       { return e.Error() }
func (e *SDKError) GoString() string     { return e.Error() }
func ValidationError(message string) *SDKError {
	return NewError("ValidationError", "InvalidParameter", message)
}
func configError(message string) *SDKError {
	return NewError("ConfigurationError", "InvalidConfiguration", message)
}
func responseError(code string) *SDKError {
	return NewError("InvalidResponseError", code, "Response does not match the native contract.")
}
func localError(err error) *SDKError {
	var sdk *SDKError
	if errors.As(err, &sdk) {
		return sdk
	}
	return ValidationError("Request does not match the native contract.")
}
func uncertain(err *SDKError) *SDKError {
	out := NewError("MutationUncertainError", "MutationUncertain", "Write outcome is unknown; reconcile with platform state before another submission.")
	out.Action, out.Attempts, out.RequestId, out.StatusCode = err.Action, err.Attempts, err.RequestId, err.StatusCode
	out.CauseCode = dara.StringValue(err.Code)
	out.retryAfter = err.retryAfter
	out.cause = err.cause
	return out
}
func errorDiagnostic(err *SDKError, action string, attempt int, status int, requestId string) {
	err.Action, err.Attempts = action, attempt
	if status != 0 {
		err.StatusCode = dara.Int(status)
	}
	if requestId != "" {
		err.RequestId = dara.String(requestId)
	}
	err.Data = dara.String(fmt.Sprintf(`{"requestId":%q,"statusCode":%d}`, requestId, status))
}
