package client

import (
	"encoding/json"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type ResponseMetadata struct {
	model.Extra `json:"-"`
	RequestId   *string `json:"RequestId,omitempty"`
	Action      *string `json:"Action,omitempty"`
	Service     *string `json:"Service,omitempty"`
}

func (s ResponseMetadata) String() string   { return dara.Prettify(s) }
func (s ResponseMetadata) GoString() string { return s.String() }
func (s ResponseMetadata) MarshalJSON() ([]byte, error) {
	type wire ResponseMetadata
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ResponseMetadata) UnmarshalJSON(data []byte) error {
	type wire ResponseMetadata
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ResponseMetadata) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ResponseMetadata) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ResponseMetadata) GetRequestId() *string                   { return s.RequestId }
func (s *ResponseMetadata) SetRequestId(v string) *ResponseMetadata { s.RequestId = &v; return s }
func (s *ResponseMetadata) GetAction() *string                      { return s.Action }
func (s *ResponseMetadata) SetAction(v string) *ResponseMetadata    { s.Action = &v; return s }
func (s *ResponseMetadata) GetService() *string                     { return s.Service }
func (s *ResponseMetadata) SetService(v string) *ResponseMetadata   { s.Service = &v; return s }
func (s *ResponseMetadata) Validate() error {
	if s == nil {
		return ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type DiscoveryResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *ResponseMetadata      `json:"ResponseMetadata,omitempty"`
	Result           map[string]interface{} `json:"Result,omitempty"`
}

func (s DiscoveryResponseBody) String() string   { return dara.Prettify(s) }
func (s DiscoveryResponseBody) GoString() string { return s.String() }
func (s DiscoveryResponseBody) MarshalJSON() ([]byte, error) {
	type wire DiscoveryResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *DiscoveryResponseBody) UnmarshalJSON(data []byte) error {
	type wire DiscoveryResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *DiscoveryResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *DiscoveryResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *DiscoveryResponseBody) GetResponseMetadata() *ResponseMetadata { return s.ResponseMetadata }
func (s *DiscoveryResponseBody) SetResponseMetadata(v *ResponseMetadata) *DiscoveryResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *DiscoveryResponseBody) GetResult() map[string]interface{} { return s.Result }
func (s *DiscoveryResponseBody) SetResult(v map[string]interface{}) *DiscoveryResponseBody {
	s.Result = v
	return s
}
func (s *DiscoveryResponseBody) Validate() error {
	if s == nil {
		return ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Result == nil {
		return responseError("InvalidResult")
	}
	return nil
}

type DiscoveryResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string     `json:"headers,omitempty"`
	StatusCode  *int32                 `json:"statusCode,omitempty"`
	Body        *DiscoveryResponseBody `json:"body,omitempty"`
	RawBody     []byte                 `json:"-"`
	SdkContext  *ResponseContext       `json:"-"`
}

func (s DiscoveryResponse) String() string   { return dara.Prettify(s) }
func (s DiscoveryResponse) GoString() string { return s.String() }
func (s DiscoveryResponse) MarshalJSON() ([]byte, error) {
	type wire DiscoveryResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *DiscoveryResponse) UnmarshalJSON(data []byte) error {
	type wire DiscoveryResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *DiscoveryResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *DiscoveryResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *DiscoveryResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *DiscoveryResponse) SetHeaders(v map[string]*string) *DiscoveryResponse {
	s.Headers = v
	return s
}
func (s *DiscoveryResponse) GetStatusCode() *int32                    { return s.StatusCode }
func (s *DiscoveryResponse) SetStatusCode(v int32) *DiscoveryResponse { s.StatusCode = &v; return s }
func (s *DiscoveryResponse) GetBody() *DiscoveryResponseBody          { return s.Body }
func (s *DiscoveryResponse) SetBody(v *DiscoveryResponseBody) *DiscoveryResponse {
	s.Body = v
	return s
}
func (s *DiscoveryResponse) SetDiagnostics(raw []byte, ctx *ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *DiscoveryResponse) Validate() error {
	if s == nil {
		return ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Body == nil {
		return responseError("InvalidResult")
	}
	return nil
}

var _ json.Marshaler = (*ResponseMetadata)(nil)
