package client

import (
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type GetUserDetailRequest struct {
	model.Extra `json:"-"`
}

func (s GetUserDetailRequest) String() string   { return dara.Prettify(s) }
func (s GetUserDetailRequest) GoString() string { return s.String() }
func (s GetUserDetailRequest) MarshalJSON() ([]byte, error) {
	type wire GetUserDetailRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetUserDetailRequest) UnmarshalJSON(data []byte) error {
	type wire GetUserDetailRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetUserDetailRequest) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetUserDetailRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetUserDetailRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type GetUserDetailResponseBodyResult struct {
	model.Extra `json:"-"`
	Id          *string `json:"id,omitempty"`
	UserId      *string `json:"user_id,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (s GetUserDetailResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s GetUserDetailResponseBodyResult) GoString() string { return s.String() }
func (s GetUserDetailResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire GetUserDetailResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetUserDetailResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire GetUserDetailResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetUserDetailResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *GetUserDetailResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetUserDetailResponseBodyResult) GetId() *string { return s.Id }
func (s *GetUserDetailResponseBodyResult) SetId(v string) *GetUserDetailResponseBodyResult {
	s.Id = &v
	return s
}
func (s *GetUserDetailResponseBodyResult) GetUserId() *string { return s.UserId }
func (s *GetUserDetailResponseBodyResult) SetUserId(v string) *GetUserDetailResponseBodyResult {
	s.UserId = &v
	return s
}
func (s *GetUserDetailResponseBodyResult) GetName() *string { return s.Name }
func (s *GetUserDetailResponseBodyResult) SetName(v string) *GetUserDetailResponseBodyResult {
	s.Name = &v
	return s
}
func (s *GetUserDetailResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type GetUserDetailResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata        `json:"ResponseMetadata,omitempty"`
	Result           *GetUserDetailResponseBodyResult `json:"Result,omitempty"`
}

func (s GetUserDetailResponseBody) String() string   { return dara.Prettify(s) }
func (s GetUserDetailResponseBody) GoString() string { return s.String() }
func (s GetUserDetailResponseBody) MarshalJSON() ([]byte, error) {
	type wire GetUserDetailResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetUserDetailResponseBody) UnmarshalJSON(data []byte) error {
	type wire GetUserDetailResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetUserDetailResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetUserDetailResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetUserDetailResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *GetUserDetailResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *GetUserDetailResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *GetUserDetailResponseBody) GetResult() *GetUserDetailResponseBodyResult { return s.Result }
func (s *GetUserDetailResponseBody) SetResult(v *GetUserDetailResponseBodyResult) *GetUserDetailResponseBody {
	s.Result = v
	return s
}
func (s *GetUserDetailResponseBody) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Result == nil {
		return openapi.NewError("InvalidResponseError", "InvalidResult", "Response omitted a valid Result.")
	}
	return nil
}

type GetUserDetailResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string         `json:"headers,omitempty"`
	StatusCode  *int32                     `json:"statusCode,omitempty"`
	Body        *GetUserDetailResponseBody `json:"body,omitempty"`
	RawBody     []byte                     `json:"-"`
	SdkContext  *openapi.ResponseContext   `json:"-"`
}

func (s GetUserDetailResponse) String() string   { return dara.Prettify(s) }
func (s GetUserDetailResponse) GoString() string { return s.String() }
func (s GetUserDetailResponse) MarshalJSON() ([]byte, error) {
	type wire GetUserDetailResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetUserDetailResponse) UnmarshalJSON(data []byte) error {
	type wire GetUserDetailResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetUserDetailResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetUserDetailResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetUserDetailResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *GetUserDetailResponse) SetHeaders(v map[string]*string) *GetUserDetailResponse {
	s.Headers = v
	return s
}
func (s *GetUserDetailResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *GetUserDetailResponse) SetStatusCode(v int32) *GetUserDetailResponse {
	s.StatusCode = &v
	return s
}
func (s *GetUserDetailResponse) GetBody() *GetUserDetailResponseBody { return s.Body }
func (s *GetUserDetailResponse) SetBody(v *GetUserDetailResponseBody) *GetUserDetailResponse {
	s.Body = v
	return s
}
func (s *GetUserDetailResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *GetUserDetailResponse) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Body == nil {
		return openapi.NewError("InvalidResponseError", "InvalidResult", "Response omitted a valid Result.")
	}
	return nil
}

var _ json.Marshaler = (*GetUserDetailRequest)(nil)
