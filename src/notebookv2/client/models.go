package client

import (
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type ListNotebooksRequestFilterBy struct {
	model.Extra         `json:"-"`
	Keyword             *string   `json:"keyword,omitempty"`
	UserId              []*string `json:"user_id,omitempty"`
	LogicComputeGroupId []*string `json:"logic_compute_group_id,omitempty"`
	Status              []*string `json:"status,omitempty"`
	MirrorUrl           []*string `json:"mirror_url,omitempty"`
}

func (s ListNotebooksRequestFilterBy) String() string   { return dara.Prettify(s) }
func (s ListNotebooksRequestFilterBy) GoString() string { return s.String() }
func (s ListNotebooksRequestFilterBy) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksRequestFilterBy
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksRequestFilterBy) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksRequestFilterBy
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksRequestFilterBy) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListNotebooksRequestFilterBy) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksRequestFilterBy) GetKeyword() *string { return s.Keyword }
func (s *ListNotebooksRequestFilterBy) SetKeyword(v string) *ListNotebooksRequestFilterBy {
	s.Keyword = &v
	return s
}
func (s *ListNotebooksRequestFilterBy) GetUserId() []*string { return s.UserId }
func (s *ListNotebooksRequestFilterBy) SetUserId(v []*string) *ListNotebooksRequestFilterBy {
	s.UserId = v
	return s
}
func (s *ListNotebooksRequestFilterBy) GetLogicComputeGroupId() []*string {
	return s.LogicComputeGroupId
}
func (s *ListNotebooksRequestFilterBy) SetLogicComputeGroupId(v []*string) *ListNotebooksRequestFilterBy {
	s.LogicComputeGroupId = v
	return s
}
func (s *ListNotebooksRequestFilterBy) GetStatus() []*string { return s.Status }
func (s *ListNotebooksRequestFilterBy) SetStatus(v []*string) *ListNotebooksRequestFilterBy {
	s.Status = v
	return s
}
func (s *ListNotebooksRequestFilterBy) GetMirrorUrl() []*string { return s.MirrorUrl }
func (s *ListNotebooksRequestFilterBy) SetMirrorUrl(v []*string) *ListNotebooksRequestFilterBy {
	s.MirrorUrl = v
	return s
}
func (s *ListNotebooksRequestFilterBy) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if len(s.UserId) == 0 {
		return openapi.ValidationError("filter_by.user_id must explicitly select users.")
	}
	for _, id := range s.UserId {
		if err := model.Required(id, "user_id"); err != nil {
			return err
		}
	}
	return nil
}

type ListNotebooksRequestOrderBy struct {
	model.Extra `json:"-"`
	Field       *string `json:"field,omitempty"`
	Order       *string `json:"order,omitempty"`
}

func (s ListNotebooksRequestOrderBy) String() string   { return dara.Prettify(s) }
func (s ListNotebooksRequestOrderBy) GoString() string { return s.String() }
func (s ListNotebooksRequestOrderBy) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksRequestOrderBy
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksRequestOrderBy) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksRequestOrderBy
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksRequestOrderBy) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListNotebooksRequestOrderBy) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksRequestOrderBy) GetField() *string { return s.Field }
func (s *ListNotebooksRequestOrderBy) SetField(v string) *ListNotebooksRequestOrderBy {
	s.Field = &v
	return s
}
func (s *ListNotebooksRequestOrderBy) GetOrder() *string { return s.Order }
func (s *ListNotebooksRequestOrderBy) SetOrder(v string) *ListNotebooksRequestOrderBy {
	s.Order = &v
	return s
}
func (s *ListNotebooksRequestOrderBy) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListNotebooksRequest struct {
	model.Extra `json:"-"`
	WorkspaceId *string                        `json:"workspace_id,omitempty"`
	Page        *int32                         `json:"page,omitempty"`
	PageSize    *int32                         `json:"page_size,omitempty"`
	FilterBy    *ListNotebooksRequestFilterBy  `json:"filter_by,omitempty"`
	OrderBy     []*ListNotebooksRequestOrderBy `json:"order_by,omitempty"`
}

func (s ListNotebooksRequest) String() string   { return dara.Prettify(s) }
func (s ListNotebooksRequest) GoString() string { return s.String() }
func (s ListNotebooksRequest) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksRequest) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksRequest) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListNotebooksRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksRequest) GetWorkspaceId() *string { return s.WorkspaceId }
func (s *ListNotebooksRequest) SetWorkspaceId(v string) *ListNotebooksRequest {
	s.WorkspaceId = &v
	return s
}
func (s *ListNotebooksRequest) GetPage() *int32                            { return s.Page }
func (s *ListNotebooksRequest) SetPage(v int32) *ListNotebooksRequest      { s.Page = &v; return s }
func (s *ListNotebooksRequest) GetPageSize() *int32                        { return s.PageSize }
func (s *ListNotebooksRequest) SetPageSize(v int32) *ListNotebooksRequest  { s.PageSize = &v; return s }
func (s *ListNotebooksRequest) GetFilterBy() *ListNotebooksRequestFilterBy { return s.FilterBy }
func (s *ListNotebooksRequest) SetFilterBy(v *ListNotebooksRequestFilterBy) *ListNotebooksRequest {
	s.FilterBy = v
	return s
}
func (s *ListNotebooksRequest) GetOrderBy() []*ListNotebooksRequestOrderBy { return s.OrderBy }
func (s *ListNotebooksRequest) SetOrderBy(v []*ListNotebooksRequestOrderBy) *ListNotebooksRequest {
	s.OrderBy = v
	return s
}
func (s *ListNotebooksRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if err := model.Required(s.WorkspaceId, "workspace_id"); err != nil {
		return err
	}
	if s.FilterBy == nil {
		return openapi.ValidationError("filter_by is required.")
	}
	return model.Page(s.Page, s.PageSize, false, 0)
}

type GetNotebookRequest struct {
	model.Extra `json:"-"`
	NotebookId  *string `json:"notebook_id,omitempty"`
}

func (s GetNotebookRequest) String() string   { return dara.Prettify(s) }
func (s GetNotebookRequest) GoString() string { return s.String() }
func (s GetNotebookRequest) MarshalJSON() ([]byte, error) {
	type wire GetNotebookRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetNotebookRequest) UnmarshalJSON(data []byte) error {
	type wire GetNotebookRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetNotebookRequest) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetNotebookRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetNotebookRequest) GetNotebookId() *string                     { return s.NotebookId }
func (s *GetNotebookRequest) SetNotebookId(v string) *GetNotebookRequest { s.NotebookId = &v; return s }
func (s *GetNotebookRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Required(s.NotebookId, "notebook_id")
}

type Notebook struct {
	model.Extra `json:"-"`
	Id          *string `json:"id,omitempty"`
	NotebookId  *string `json:"notebook_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (s Notebook) String() string   { return dara.Prettify(s) }
func (s Notebook) GoString() string { return s.String() }
func (s Notebook) MarshalJSON() ([]byte, error) {
	type wire Notebook
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *Notebook) UnmarshalJSON(data []byte) error {
	type wire Notebook
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *Notebook) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *Notebook) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *Notebook) GetId() *string                             { return s.Id }
func (s *Notebook) SetId(v string) *Notebook                   { s.Id = &v; return s }
func (s *Notebook) GetNotebookId() *string                     { return s.NotebookId }
func (s *Notebook) SetNotebookId(v string) *Notebook           { s.NotebookId = &v; return s }
func (s *Notebook) GetName() *string                           { return s.Name }
func (s *Notebook) SetName(v string) *Notebook                 { s.Name = &v; return s }
func (s *Notebook) GetStatus() *string                         { return s.Status }
func (s *Notebook) SetStatus(v string) *Notebook               { s.Status = &v; return s }
func (s *Notebook) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListNotebooksResponseBodyResult struct {
	model.Extra `json:"-"`
	List        []*Notebook `json:"list,omitempty"`
	Total       interface{} `json:"total,omitempty"`
}

func (s ListNotebooksResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s ListNotebooksResponseBodyResult) GoString() string { return s.String() }
func (s ListNotebooksResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *ListNotebooksResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksResponseBodyResult) GetList() []*Notebook { return s.List }
func (s *ListNotebooksResponseBodyResult) SetList(v []*Notebook) *ListNotebooksResponseBodyResult {
	s.List = v
	return s
}
func (s *ListNotebooksResponseBodyResult) GetTotal() interface{} { return s.Total }
func (s *ListNotebooksResponseBodyResult) SetTotal(v interface{}) *ListNotebooksResponseBodyResult {
	s.Total = v
	return s
}
func (s *ListNotebooksResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.List == nil {
		return openapi.NewError("InvalidResponseError", "InvalidCollection", "Response omitted its collection.")
	}
	return nil
}

type GetNotebookResponseBodyResult struct {
	model.Extra `json:"-"`
	Id          *string `json:"id,omitempty"`
	NotebookId  *string `json:"notebook_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (s GetNotebookResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s GetNotebookResponseBodyResult) GoString() string { return s.String() }
func (s GetNotebookResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire GetNotebookResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetNotebookResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire GetNotebookResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetNotebookResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *GetNotebookResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetNotebookResponseBodyResult) GetId() *string { return s.Id }
func (s *GetNotebookResponseBodyResult) SetId(v string) *GetNotebookResponseBodyResult {
	s.Id = &v
	return s
}
func (s *GetNotebookResponseBodyResult) GetNotebookId() *string { return s.NotebookId }
func (s *GetNotebookResponseBodyResult) SetNotebookId(v string) *GetNotebookResponseBodyResult {
	s.NotebookId = &v
	return s
}
func (s *GetNotebookResponseBodyResult) GetName() *string { return s.Name }
func (s *GetNotebookResponseBodyResult) SetName(v string) *GetNotebookResponseBodyResult {
	s.Name = &v
	return s
}
func (s *GetNotebookResponseBodyResult) GetStatus() *string { return s.Status }
func (s *GetNotebookResponseBodyResult) SetStatus(v string) *GetNotebookResponseBodyResult {
	s.Status = &v
	return s
}
func (s *GetNotebookResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListNotebooksResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata        `json:"ResponseMetadata,omitempty"`
	Result           *ListNotebooksResponseBodyResult `json:"Result,omitempty"`
}

func (s ListNotebooksResponseBody) String() string   { return dara.Prettify(s) }
func (s ListNotebooksResponseBody) GoString() string { return s.String() }
func (s ListNotebooksResponseBody) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksResponseBody) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListNotebooksResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *ListNotebooksResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *ListNotebooksResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *ListNotebooksResponseBody) GetResult() *ListNotebooksResponseBodyResult { return s.Result }
func (s *ListNotebooksResponseBody) SetResult(v *ListNotebooksResponseBodyResult) *ListNotebooksResponseBody {
	s.Result = v
	return s
}
func (s *ListNotebooksResponseBody) Validate() error {
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

type ListNotebooksResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string         `json:"headers,omitempty"`
	StatusCode  *int32                     `json:"statusCode,omitempty"`
	Body        *ListNotebooksResponseBody `json:"body,omitempty"`
	RawBody     []byte                     `json:"-"`
	SdkContext  *openapi.ResponseContext   `json:"-"`
}

func (s ListNotebooksResponse) String() string   { return dara.Prettify(s) }
func (s ListNotebooksResponse) GoString() string { return s.String() }
func (s ListNotebooksResponse) MarshalJSON() ([]byte, error) {
	type wire ListNotebooksResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListNotebooksResponse) UnmarshalJSON(data []byte) error {
	type wire ListNotebooksResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListNotebooksResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListNotebooksResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListNotebooksResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *ListNotebooksResponse) SetHeaders(v map[string]*string) *ListNotebooksResponse {
	s.Headers = v
	return s
}
func (s *ListNotebooksResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *ListNotebooksResponse) SetStatusCode(v int32) *ListNotebooksResponse {
	s.StatusCode = &v
	return s
}
func (s *ListNotebooksResponse) GetBody() *ListNotebooksResponseBody { return s.Body }
func (s *ListNotebooksResponse) SetBody(v *ListNotebooksResponseBody) *ListNotebooksResponse {
	s.Body = v
	return s
}
func (s *ListNotebooksResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *ListNotebooksResponse) Validate() error {
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

type GetNotebookResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata      `json:"ResponseMetadata,omitempty"`
	Result           *GetNotebookResponseBodyResult `json:"Result,omitempty"`
}

func (s GetNotebookResponseBody) String() string   { return dara.Prettify(s) }
func (s GetNotebookResponseBody) GoString() string { return s.String() }
func (s GetNotebookResponseBody) MarshalJSON() ([]byte, error) {
	type wire GetNotebookResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetNotebookResponseBody) UnmarshalJSON(data []byte) error {
	type wire GetNotebookResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetNotebookResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetNotebookResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetNotebookResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *GetNotebookResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *GetNotebookResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *GetNotebookResponseBody) GetResult() *GetNotebookResponseBodyResult { return s.Result }
func (s *GetNotebookResponseBody) SetResult(v *GetNotebookResponseBodyResult) *GetNotebookResponseBody {
	s.Result = v
	return s
}
func (s *GetNotebookResponseBody) Validate() error {
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

type GetNotebookResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string       `json:"headers,omitempty"`
	StatusCode  *int32                   `json:"statusCode,omitempty"`
	Body        *GetNotebookResponseBody `json:"body,omitempty"`
	RawBody     []byte                   `json:"-"`
	SdkContext  *openapi.ResponseContext `json:"-"`
}

func (s GetNotebookResponse) String() string   { return dara.Prettify(s) }
func (s GetNotebookResponse) GoString() string { return s.String() }
func (s GetNotebookResponse) MarshalJSON() ([]byte, error) {
	type wire GetNotebookResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetNotebookResponse) UnmarshalJSON(data []byte) error {
	type wire GetNotebookResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetNotebookResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetNotebookResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetNotebookResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *GetNotebookResponse) SetHeaders(v map[string]*string) *GetNotebookResponse {
	s.Headers = v
	return s
}
func (s *GetNotebookResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *GetNotebookResponse) SetStatusCode(v int32) *GetNotebookResponse {
	s.StatusCode = &v
	return s
}
func (s *GetNotebookResponse) GetBody() *GetNotebookResponseBody { return s.Body }
func (s *GetNotebookResponse) SetBody(v *GetNotebookResponseBody) *GetNotebookResponse {
	s.Body = v
	return s
}
func (s *GetNotebookResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *GetNotebookResponse) Validate() error {
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

var _ json.Marshaler = (*ListNotebooksRequestFilterBy)(nil)
