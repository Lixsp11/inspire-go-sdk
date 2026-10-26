package client

import (
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type ListProjectsRequestFilter struct {
	model.Extra `json:"-"`
	WorkspaceId *string `json:"workspace_id,omitempty"`
	CheckAdmin  *bool   `json:"check_admin,omitempty"`
}

func (s ListProjectsRequestFilter) String() string   { return dara.Prettify(s) }
func (s ListProjectsRequestFilter) GoString() string { return s.String() }
func (s ListProjectsRequestFilter) MarshalJSON() ([]byte, error) {
	type wire ListProjectsRequestFilter
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListProjectsRequestFilter) UnmarshalJSON(data []byte) error {
	type wire ListProjectsRequestFilter
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListProjectsRequestFilter) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListProjectsRequestFilter) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListProjectsRequestFilter) GetWorkspaceId() *string { return s.WorkspaceId }
func (s *ListProjectsRequestFilter) SetWorkspaceId(v string) *ListProjectsRequestFilter {
	s.WorkspaceId = &v
	return s
}
func (s *ListProjectsRequestFilter) GetCheckAdmin() *bool { return s.CheckAdmin }
func (s *ListProjectsRequestFilter) SetCheckAdmin(v bool) *ListProjectsRequestFilter {
	s.CheckAdmin = &v
	return s
}
func (s *ListProjectsRequestFilter) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListProjectsRequest struct {
	model.Extra `json:"-"`
	Page        *int32                     `json:"page,omitempty"`
	PageSize    *int32                     `json:"page_size,omitempty"`
	Filter      *ListProjectsRequestFilter `json:"filter,omitempty"`
}

func (s ListProjectsRequest) String() string   { return dara.Prettify(s) }
func (s ListProjectsRequest) GoString() string { return s.String() }
func (s ListProjectsRequest) MarshalJSON() ([]byte, error) {
	type wire ListProjectsRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListProjectsRequest) UnmarshalJSON(data []byte) error {
	type wire ListProjectsRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListProjectsRequest) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListProjectsRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListProjectsRequest) GetPage() *int32                          { return s.Page }
func (s *ListProjectsRequest) SetPage(v int32) *ListProjectsRequest     { s.Page = &v; return s }
func (s *ListProjectsRequest) GetPageSize() *int32                      { return s.PageSize }
func (s *ListProjectsRequest) SetPageSize(v int32) *ListProjectsRequest { s.PageSize = &v; return s }
func (s *ListProjectsRequest) GetFilter() *ListProjectsRequestFilter    { return s.Filter }
func (s *ListProjectsRequest) SetFilter(v *ListProjectsRequestFilter) *ListProjectsRequest {
	s.Filter = v
	return s
}
func (s *ListProjectsRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Page(s.Page, s.PageSize, true, 0)
}

type Project struct {
	model.Extra `json:"-"`
	Id          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (s Project) String() string   { return dara.Prettify(s) }
func (s Project) GoString() string { return s.String() }
func (s Project) MarshalJSON() ([]byte, error) {
	type wire Project
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *Project) UnmarshalJSON(data []byte) error {
	type wire Project
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *Project) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *Project) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *Project) GetId() *string                             { return s.Id }
func (s *Project) SetId(v string) *Project                    { s.Id = &v; return s }
func (s *Project) GetName() *string                           { return s.Name }
func (s *Project) SetName(v string) *Project                  { s.Name = &v; return s }
func (s *Project) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListProjectsResponseBodyResult struct {
	model.Extra `json:"-"`
	Items       []*Project  `json:"items,omitempty"`
	Total       interface{} `json:"total,omitempty"`
}

func (s ListProjectsResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s ListProjectsResponseBodyResult) GoString() string { return s.String() }
func (s ListProjectsResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire ListProjectsResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListProjectsResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire ListProjectsResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListProjectsResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListProjectsResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListProjectsResponseBodyResult) GetItems() []*Project { return s.Items }
func (s *ListProjectsResponseBodyResult) SetItems(v []*Project) *ListProjectsResponseBodyResult {
	s.Items = v
	return s
}
func (s *ListProjectsResponseBodyResult) GetTotal() interface{} { return s.Total }
func (s *ListProjectsResponseBodyResult) SetTotal(v interface{}) *ListProjectsResponseBodyResult {
	s.Total = v
	return s
}
func (s *ListProjectsResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Items == nil {
		return openapi.NewError("InvalidResponseError", "InvalidCollection", "Response omitted its collection.")
	}
	return nil
}

type ListProjectsResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata       `json:"ResponseMetadata,omitempty"`
	Result           *ListProjectsResponseBodyResult `json:"Result,omitempty"`
}

func (s ListProjectsResponseBody) String() string   { return dara.Prettify(s) }
func (s ListProjectsResponseBody) GoString() string { return s.String() }
func (s ListProjectsResponseBody) MarshalJSON() ([]byte, error) {
	type wire ListProjectsResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListProjectsResponseBody) UnmarshalJSON(data []byte) error {
	type wire ListProjectsResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListProjectsResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListProjectsResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListProjectsResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *ListProjectsResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *ListProjectsResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *ListProjectsResponseBody) GetResult() *ListProjectsResponseBodyResult { return s.Result }
func (s *ListProjectsResponseBody) SetResult(v *ListProjectsResponseBodyResult) *ListProjectsResponseBody {
	s.Result = v
	return s
}
func (s *ListProjectsResponseBody) Validate() error {
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

type ListProjectsResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string        `json:"headers,omitempty"`
	StatusCode  *int32                    `json:"statusCode,omitempty"`
	Body        *ListProjectsResponseBody `json:"body,omitempty"`
	RawBody     []byte                    `json:"-"`
	SdkContext  *openapi.ResponseContext  `json:"-"`
}

func (s ListProjectsResponse) String() string   { return dara.Prettify(s) }
func (s ListProjectsResponse) GoString() string { return s.String() }
func (s ListProjectsResponse) MarshalJSON() ([]byte, error) {
	type wire ListProjectsResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListProjectsResponse) UnmarshalJSON(data []byte) error {
	type wire ListProjectsResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListProjectsResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListProjectsResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListProjectsResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *ListProjectsResponse) SetHeaders(v map[string]*string) *ListProjectsResponse {
	s.Headers = v
	return s
}
func (s *ListProjectsResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *ListProjectsResponse) SetStatusCode(v int32) *ListProjectsResponse {
	s.StatusCode = &v
	return s
}
func (s *ListProjectsResponse) GetBody() *ListProjectsResponseBody { return s.Body }
func (s *ListProjectsResponse) SetBody(v *ListProjectsResponseBody) *ListProjectsResponse {
	s.Body = v
	return s
}
func (s *ListProjectsResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *ListProjectsResponse) Validate() error {
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

var _ json.Marshaler = (*ListProjectsRequestFilter)(nil)
