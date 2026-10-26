package client

import (
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type ListLogicComputeGroupsRequestFilter struct {
	model.Extra `json:"-"`
	WorkspaceId *string `json:"workspace_id,omitempty"`
}

func (s ListLogicComputeGroupsRequestFilter) String() string   { return dara.Prettify(s) }
func (s ListLogicComputeGroupsRequestFilter) GoString() string { return s.String() }
func (s ListLogicComputeGroupsRequestFilter) MarshalJSON() ([]byte, error) {
	type wire ListLogicComputeGroupsRequestFilter
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListLogicComputeGroupsRequestFilter) UnmarshalJSON(data []byte) error {
	type wire ListLogicComputeGroupsRequestFilter
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListLogicComputeGroupsRequestFilter) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListLogicComputeGroupsRequestFilter) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListLogicComputeGroupsRequestFilter) GetWorkspaceId() *string { return s.WorkspaceId }
func (s *ListLogicComputeGroupsRequestFilter) SetWorkspaceId(v string) *ListLogicComputeGroupsRequestFilter {
	s.WorkspaceId = &v
	return s
}
func (s *ListLogicComputeGroupsRequestFilter) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Required(s.WorkspaceId, "workspace_id")
}

type ListLogicComputeGroupsRequest struct {
	model.Extra `json:"-"`
	PageSize    *int32                               `json:"page_size,omitempty"`
	PageNum     *int32                               `json:"page_num,omitempty"`
	Filter      *ListLogicComputeGroupsRequestFilter `json:"filter,omitempty"`
}

func (s ListLogicComputeGroupsRequest) String() string   { return dara.Prettify(s) }
func (s ListLogicComputeGroupsRequest) GoString() string { return s.String() }
func (s ListLogicComputeGroupsRequest) MarshalJSON() ([]byte, error) {
	type wire ListLogicComputeGroupsRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListLogicComputeGroupsRequest) UnmarshalJSON(data []byte) error {
	type wire ListLogicComputeGroupsRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListLogicComputeGroupsRequest) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListLogicComputeGroupsRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListLogicComputeGroupsRequest) GetPageSize() *int32 { return s.PageSize }
func (s *ListLogicComputeGroupsRequest) SetPageSize(v int32) *ListLogicComputeGroupsRequest {
	s.PageSize = &v
	return s
}
func (s *ListLogicComputeGroupsRequest) GetPageNum() *int32 { return s.PageNum }
func (s *ListLogicComputeGroupsRequest) SetPageNum(v int32) *ListLogicComputeGroupsRequest {
	s.PageNum = &v
	return s
}
func (s *ListLogicComputeGroupsRequest) GetFilter() *ListLogicComputeGroupsRequestFilter {
	return s.Filter
}
func (s *ListLogicComputeGroupsRequest) SetFilter(v *ListLogicComputeGroupsRequestFilter) *ListLogicComputeGroupsRequest {
	s.Filter = v
	return s
}
func (s *ListLogicComputeGroupsRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Filter == nil || s.PageSize == nil {
		return openapi.ValidationError("filter and explicit page_size are required.")
	}
	return model.Page(s.PageNum, s.PageSize, true, 0)
}

type LogicComputeGroup struct {
	model.Extra `json:"-"`
	Id          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (s LogicComputeGroup) String() string   { return dara.Prettify(s) }
func (s LogicComputeGroup) GoString() string { return s.String() }
func (s LogicComputeGroup) MarshalJSON() ([]byte, error) {
	type wire LogicComputeGroup
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *LogicComputeGroup) UnmarshalJSON(data []byte) error {
	type wire LogicComputeGroup
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *LogicComputeGroup) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *LogicComputeGroup) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *LogicComputeGroup) GetId() *string                      { return s.Id }
func (s *LogicComputeGroup) SetId(v string) *LogicComputeGroup   { s.Id = &v; return s }
func (s *LogicComputeGroup) GetName() *string                    { return s.Name }
func (s *LogicComputeGroup) SetName(v string) *LogicComputeGroup { s.Name = &v; return s }
func (s *LogicComputeGroup) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListLogicComputeGroupsResponseBodyResult struct {
	model.Extra        `json:"-"`
	LogicComputeGroups []*LogicComputeGroup `json:"logic_compute_groups,omitempty"`
	Total              interface{}          `json:"total,omitempty"`
}

func (s ListLogicComputeGroupsResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s ListLogicComputeGroupsResponseBodyResult) GoString() string { return s.String() }
func (s ListLogicComputeGroupsResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire ListLogicComputeGroupsResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire ListLogicComputeGroupsResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListLogicComputeGroupsResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListLogicComputeGroupsResponseBodyResult) GetLogicComputeGroups() []*LogicComputeGroup {
	return s.LogicComputeGroups
}
func (s *ListLogicComputeGroupsResponseBodyResult) SetLogicComputeGroups(v []*LogicComputeGroup) *ListLogicComputeGroupsResponseBodyResult {
	s.LogicComputeGroups = v
	return s
}
func (s *ListLogicComputeGroupsResponseBodyResult) GetTotal() interface{} { return s.Total }
func (s *ListLogicComputeGroupsResponseBodyResult) SetTotal(v interface{}) *ListLogicComputeGroupsResponseBodyResult {
	s.Total = v
	return s
}
func (s *ListLogicComputeGroupsResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.LogicComputeGroups == nil {
		return openapi.NewError("InvalidResponseError", "InvalidCollection", "Response omitted its collection.")
	}
	return nil
}

type ListLogicComputeGroupsResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata                 `json:"ResponseMetadata,omitempty"`
	Result           *ListLogicComputeGroupsResponseBodyResult `json:"Result,omitempty"`
}

func (s ListLogicComputeGroupsResponseBody) String() string   { return dara.Prettify(s) }
func (s ListLogicComputeGroupsResponseBody) GoString() string { return s.String() }
func (s ListLogicComputeGroupsResponseBody) MarshalJSON() ([]byte, error) {
	type wire ListLogicComputeGroupsResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponseBody) UnmarshalJSON(data []byte) error {
	type wire ListLogicComputeGroupsResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponseBody) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListLogicComputeGroupsResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListLogicComputeGroupsResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *ListLogicComputeGroupsResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *ListLogicComputeGroupsResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *ListLogicComputeGroupsResponseBody) GetResult() *ListLogicComputeGroupsResponseBodyResult {
	return s.Result
}
func (s *ListLogicComputeGroupsResponseBody) SetResult(v *ListLogicComputeGroupsResponseBodyResult) *ListLogicComputeGroupsResponseBody {
	s.Result = v
	return s
}
func (s *ListLogicComputeGroupsResponseBody) Validate() error {
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

type ListLogicComputeGroupsResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string                  `json:"headers,omitempty"`
	StatusCode  *int32                              `json:"statusCode,omitempty"`
	Body        *ListLogicComputeGroupsResponseBody `json:"body,omitempty"`
	RawBody     []byte                              `json:"-"`
	SdkContext  *openapi.ResponseContext            `json:"-"`
}

func (s ListLogicComputeGroupsResponse) String() string   { return dara.Prettify(s) }
func (s ListLogicComputeGroupsResponse) GoString() string { return s.String() }
func (s ListLogicComputeGroupsResponse) MarshalJSON() ([]byte, error) {
	type wire ListLogicComputeGroupsResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponse) UnmarshalJSON(data []byte) error {
	type wire ListLogicComputeGroupsResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListLogicComputeGroupsResponse) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListLogicComputeGroupsResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListLogicComputeGroupsResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *ListLogicComputeGroupsResponse) SetHeaders(v map[string]*string) *ListLogicComputeGroupsResponse {
	s.Headers = v
	return s
}
func (s *ListLogicComputeGroupsResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *ListLogicComputeGroupsResponse) SetStatusCode(v int32) *ListLogicComputeGroupsResponse {
	s.StatusCode = &v
	return s
}
func (s *ListLogicComputeGroupsResponse) GetBody() *ListLogicComputeGroupsResponseBody { return s.Body }
func (s *ListLogicComputeGroupsResponse) SetBody(v *ListLogicComputeGroupsResponseBody) *ListLogicComputeGroupsResponse {
	s.Body = v
	return s
}
func (s *ListLogicComputeGroupsResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *ListLogicComputeGroupsResponse) Validate() error {
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

var _ json.Marshaler = (*ListLogicComputeGroupsRequestFilter)(nil)
