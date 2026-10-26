package client

import (
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
)

type Job struct {
	model.Extra `json:"-"`
	JobId       *string `json:"job_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (s Job) String() string   { return dara.Prettify(s) }
func (s Job) GoString() string { return s.String() }
func (s Job) MarshalJSON() ([]byte, error) {
	type wire Job
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *Job) UnmarshalJSON(data []byte) error {
	type wire Job
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *Job) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *Job) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *Job) GetJobId() *string                          { return s.JobId }
func (s *Job) SetJobId(v string) *Job                     { s.JobId = &v; return s }
func (s *Job) GetName() *string                           { return s.Name }
func (s *Job) SetName(v string) *Job                      { s.Name = &v; return s }
func (s *Job) GetStatus() *string                         { return s.Status }
func (s *Job) SetStatus(v string) *Job                    { s.Status = &v; return s }
func (s *Job) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListJobsRequest struct {
	model.Extra `json:"-"`
	WorkspaceId *string `json:"workspace_id,omitempty"`
	CreatedBy   *string `json:"created_by,omitempty"`
	Status      *string `json:"status,omitempty"`
	Keyword     *string `json:"keyword,omitempty"`
	PageNum     *int32  `json:"page_num,omitempty"`
	PageSize    *int32  `json:"page_size,omitempty"`
}

func (s ListJobsRequest) String() string   { return dara.Prettify(s) }
func (s ListJobsRequest) GoString() string { return s.String() }
func (s ListJobsRequest) MarshalJSON() ([]byte, error) {
	type wire ListJobsRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListJobsRequest) UnmarshalJSON(data []byte) error {
	type wire ListJobsRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListJobsRequest) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *ListJobsRequest) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *ListJobsRequest) GetWorkspaceId() *string                    { return s.WorkspaceId }
func (s *ListJobsRequest) SetWorkspaceId(v string) *ListJobsRequest   { s.WorkspaceId = &v; return s }
func (s *ListJobsRequest) GetCreatedBy() *string                      { return s.CreatedBy }
func (s *ListJobsRequest) SetCreatedBy(v string) *ListJobsRequest     { s.CreatedBy = &v; return s }
func (s *ListJobsRequest) GetStatus() *string                         { return s.Status }
func (s *ListJobsRequest) SetStatus(v string) *ListJobsRequest        { s.Status = &v; return s }
func (s *ListJobsRequest) GetKeyword() *string                        { return s.Keyword }
func (s *ListJobsRequest) SetKeyword(v string) *ListJobsRequest       { s.Keyword = &v; return s }
func (s *ListJobsRequest) GetPageNum() *int32                         { return s.PageNum }
func (s *ListJobsRequest) SetPageNum(v int32) *ListJobsRequest        { s.PageNum = &v; return s }
func (s *ListJobsRequest) GetPageSize() *int32                        { return s.PageSize }
func (s *ListJobsRequest) SetPageSize(v int32) *ListJobsRequest       { s.PageSize = &v; return s }
func (s *ListJobsRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if err := model.Required(s.WorkspaceId, "workspace_id"); err != nil {
		return err
	}
	return model.Page(s.PageNum, s.PageSize, false, 999)
}

type GetJobRequest struct {
	model.Extra `json:"-"`
	JobId       *string `json:"job_id,omitempty"`
}

func (s GetJobRequest) String() string   { return dara.Prettify(s) }
func (s GetJobRequest) GoString() string { return s.String() }
func (s GetJobRequest) MarshalJSON() ([]byte, error) {
	type wire GetJobRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetJobRequest) UnmarshalJSON(data []byte) error {
	type wire GetJobRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetJobRequest) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *GetJobRequest) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *GetJobRequest) GetJobId() *string                          { return s.JobId }
func (s *GetJobRequest) SetJobId(v string) *GetJobRequest           { s.JobId = &v; return s }
func (s *GetJobRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Required(s.JobId, "job_id")
}

type StopJobRequest struct {
	model.Extra `json:"-"`
	JobId       *string `json:"job_id,omitempty"`
}

func (s StopJobRequest) String() string   { return dara.Prettify(s) }
func (s StopJobRequest) GoString() string { return s.String() }
func (s StopJobRequest) MarshalJSON() ([]byte, error) {
	type wire StopJobRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *StopJobRequest) UnmarshalJSON(data []byte) error {
	type wire StopJobRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *StopJobRequest) ToMap() (map[string]interface{}, error)     { return model.ToMap(s) }
func (s *StopJobRequest) FromMap(value map[string]interface{}) error { return model.FromMap(value, s) }
func (s *StopJobRequest) GetJobId() *string                          { return s.JobId }
func (s *StopJobRequest) SetJobId(v string) *StopJobRequest          { s.JobId = &v; return s }
func (s *StopJobRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Required(s.JobId, "job_id")
}

type CreateJobConsoleRequestFrameworkConfigResourceSpecPrice struct {
	model.Extra         `json:"-"`
	CpuType             *string  `json:"cpu_type,omitempty"`
	CpuCount            *float64 `json:"cpu_count,omitempty"`
	GpuType             *string  `json:"gpu_type,omitempty"`
	GpuCount            *int32   `json:"gpu_count,omitempty"`
	MemorySizeGib       *float64 `json:"memory_size_gib,omitempty"`
	LogicComputeGroupId *string  `json:"logic_compute_group_id,omitempty"`
	QuotaId             *string  `json:"quota_id,omitempty"`
}

func (s CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) String() string {
	return dara.Prettify(s)
}
func (s CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GoString() string { return s.String() }
func (s CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleRequestFrameworkConfigResourceSpecPrice
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleRequestFrameworkConfigResourceSpecPrice
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetCpuType() *string {
	return s.CpuType
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetCpuType(v string) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.CpuType = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetCpuCount() *float64 {
	return s.CpuCount
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetCpuCount(v float64) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.CpuCount = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetGpuType() *string {
	return s.GpuType
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetGpuType(v string) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.GpuType = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetGpuCount() *int32 {
	return s.GpuCount
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetGpuCount(v int32) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.GpuCount = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetMemorySizeGib() *float64 {
	return s.MemorySizeGib
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetMemorySizeGib(v float64) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.MemorySizeGib = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetLogicComputeGroupId() *string {
	return s.LogicComputeGroupId
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetLogicComputeGroupId(v string) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.LogicComputeGroupId = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) GetQuotaId() *string {
	return s.QuotaId
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) SetQuotaId(v string) *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	s.QuotaId = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if err := model.Required(s.LogicComputeGroupId, "logic_compute_group_id"); err != nil {
		return err
	}
	if err := model.Required(s.QuotaId, "quota_id"); err != nil {
		return err
	}
	if s.GpuCount != nil && *s.GpuCount > 0 {
		return model.Required(s.GpuType, "gpu_type")
	}
	return nil
}

type CreateJobConsoleRequestFrameworkConfig struct {
	model.Extra       `json:"-"`
	Cpu               *float64                                                 `json:"cpu,omitempty"`
	GpuCount          *int32                                                   `json:"gpu_count,omitempty"`
	MemGi             *float64                                                 `json:"mem_gi,omitempty"`
	ShmGi             *float64                                                 `json:"shm_gi,omitempty"`
	Image             *string                                                  `json:"image,omitempty"`
	ImageType         *string                                                  `json:"image_type,omitempty"`
	InstanceCount     *int32                                                   `json:"instance_count,omitempty"`
	ResourceSpecPrice *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice `json:"resource_spec_price,omitempty"`
}

func (s CreateJobConsoleRequestFrameworkConfig) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleRequestFrameworkConfig) GoString() string { return s.String() }
func (s CreateJobConsoleRequestFrameworkConfig) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleRequestFrameworkConfig
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleRequestFrameworkConfig) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleRequestFrameworkConfig
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleRequestFrameworkConfig) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *CreateJobConsoleRequestFrameworkConfig) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetCpu() *float64 { return s.Cpu }
func (s *CreateJobConsoleRequestFrameworkConfig) SetCpu(v float64) *CreateJobConsoleRequestFrameworkConfig {
	s.Cpu = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetGpuCount() *int32 { return s.GpuCount }
func (s *CreateJobConsoleRequestFrameworkConfig) SetGpuCount(v int32) *CreateJobConsoleRequestFrameworkConfig {
	s.GpuCount = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetMemGi() *float64 { return s.MemGi }
func (s *CreateJobConsoleRequestFrameworkConfig) SetMemGi(v float64) *CreateJobConsoleRequestFrameworkConfig {
	s.MemGi = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetShmGi() *float64 { return s.ShmGi }
func (s *CreateJobConsoleRequestFrameworkConfig) SetShmGi(v float64) *CreateJobConsoleRequestFrameworkConfig {
	s.ShmGi = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetImage() *string { return s.Image }
func (s *CreateJobConsoleRequestFrameworkConfig) SetImage(v string) *CreateJobConsoleRequestFrameworkConfig {
	s.Image = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetImageType() *string { return s.ImageType }
func (s *CreateJobConsoleRequestFrameworkConfig) SetImageType(v string) *CreateJobConsoleRequestFrameworkConfig {
	s.ImageType = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetInstanceCount() *int32 { return s.InstanceCount }
func (s *CreateJobConsoleRequestFrameworkConfig) SetInstanceCount(v int32) *CreateJobConsoleRequestFrameworkConfig {
	s.InstanceCount = &v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) GetResourceSpecPrice() *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice {
	return s.ResourceSpecPrice
}
func (s *CreateJobConsoleRequestFrameworkConfig) SetResourceSpecPrice(v *CreateJobConsoleRequestFrameworkConfigResourceSpecPrice) *CreateJobConsoleRequestFrameworkConfig {
	s.ResourceSpecPrice = v
	return s
}
func (s *CreateJobConsoleRequestFrameworkConfig) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if err := model.Required(s.Image, "image"); err != nil {
		return err
	}
	if err := model.Required(s.ImageType, "image_type"); err != nil {
		return err
	}
	if err := model.Positive(s.Cpu, "cpu"); err != nil {
		return err
	}
	if err := model.Positive(s.MemGi, "mem_gi"); err != nil {
		return err
	}
	if s.GpuCount == nil || *s.GpuCount < 0 || s.InstanceCount == nil || *s.InstanceCount < 1 || s.ResourceSpecPrice == nil {
		return openapi.ValidationError("GPU count, instance count and price specification are required.")
	}
	if s.ShmGi != nil {
		if err := model.Positive(s.ShmGi, "shm_gi"); err != nil {
			return err
		}
		if *s.ShmGi > *s.MemGi {
			return openapi.ValidationError("Shared memory cannot exceed memory.")
		}
	}
	return nil
}

type CreateJobConsoleRequestEnvs struct {
	model.Extra `json:"-"`
	Name        *string `json:"name,omitempty"`
	Value       *string `json:"value,omitempty"`
}

func (s CreateJobConsoleRequestEnvs) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleRequestEnvs) GoString() string { return s.String() }
func (s CreateJobConsoleRequestEnvs) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleRequestEnvs
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleRequestEnvs) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleRequestEnvs
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleRequestEnvs) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *CreateJobConsoleRequestEnvs) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleRequestEnvs) GetName() *string { return s.Name }
func (s *CreateJobConsoleRequestEnvs) SetName(v string) *CreateJobConsoleRequestEnvs {
	s.Name = &v
	return s
}
func (s *CreateJobConsoleRequestEnvs) GetValue() *string { return s.Value }
func (s *CreateJobConsoleRequestEnvs) SetValue(v string) *CreateJobConsoleRequestEnvs {
	s.Value = &v
	return s
}
func (s *CreateJobConsoleRequestEnvs) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return model.Required(s.Name, "name")
}

type CreateJobConsoleRequest struct {
	model.Extra         `json:"-"`
	Name                *string                                   `json:"name,omitempty"`
	WorkspaceId         *string                                   `json:"workspace_id,omitempty"`
	ProjectId           *string                                   `json:"project_id,omitempty"`
	LogicComputeGroupId *string                                   `json:"logic_compute_group_id,omitempty"`
	Command             *string                                   `json:"command,omitempty"`
	Framework           *string                                   `json:"framework,omitempty"`
	TaskPriority        *int32                                    `json:"task_priority,omitempty"`
	FrameworkConfig     []*CreateJobConsoleRequestFrameworkConfig `json:"framework_config,omitempty"`
	EnableNotification  *bool                                     `json:"enable_notification,omitempty"`
	DatasetInfo         []map[string]interface{}                  `json:"dataset_info,omitempty"`
	Envs                []*CreateJobConsoleRequestEnvs            `json:"envs,omitempty"`
	SpecifiedNodes      []*string                                 `json:"specified_nodes,omitempty"`
	ExcludeNodes        []*string                                 `json:"exclude_nodes,omitempty"`
	MaxRunningTimeMs    *string                                   `json:"max_running_time_ms,omitempty"`
}

func (s CreateJobConsoleRequest) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleRequest) GoString() string { return s.String() }
func (s CreateJobConsoleRequest) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleRequest
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleRequest) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleRequest
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleRequest) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *CreateJobConsoleRequest) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleRequest) GetName() *string                          { return s.Name }
func (s *CreateJobConsoleRequest) SetName(v string) *CreateJobConsoleRequest { s.Name = &v; return s }
func (s *CreateJobConsoleRequest) GetWorkspaceId() *string                   { return s.WorkspaceId }
func (s *CreateJobConsoleRequest) SetWorkspaceId(v string) *CreateJobConsoleRequest {
	s.WorkspaceId = &v
	return s
}
func (s *CreateJobConsoleRequest) GetProjectId() *string { return s.ProjectId }
func (s *CreateJobConsoleRequest) SetProjectId(v string) *CreateJobConsoleRequest {
	s.ProjectId = &v
	return s
}
func (s *CreateJobConsoleRequest) GetLogicComputeGroupId() *string { return s.LogicComputeGroupId }
func (s *CreateJobConsoleRequest) SetLogicComputeGroupId(v string) *CreateJobConsoleRequest {
	s.LogicComputeGroupId = &v
	return s
}
func (s *CreateJobConsoleRequest) GetCommand() *string { return s.Command }
func (s *CreateJobConsoleRequest) SetCommand(v string) *CreateJobConsoleRequest {
	s.Command = &v
	return s
}
func (s *CreateJobConsoleRequest) GetFramework() *string { return s.Framework }
func (s *CreateJobConsoleRequest) SetFramework(v string) *CreateJobConsoleRequest {
	s.Framework = &v
	return s
}
func (s *CreateJobConsoleRequest) GetTaskPriority() *int32 { return s.TaskPriority }
func (s *CreateJobConsoleRequest) SetTaskPriority(v int32) *CreateJobConsoleRequest {
	s.TaskPriority = &v
	return s
}
func (s *CreateJobConsoleRequest) GetFrameworkConfig() []*CreateJobConsoleRequestFrameworkConfig {
	return s.FrameworkConfig
}
func (s *CreateJobConsoleRequest) SetFrameworkConfig(v []*CreateJobConsoleRequestFrameworkConfig) *CreateJobConsoleRequest {
	s.FrameworkConfig = v
	return s
}
func (s *CreateJobConsoleRequest) GetEnableNotification() *bool { return s.EnableNotification }
func (s *CreateJobConsoleRequest) SetEnableNotification(v bool) *CreateJobConsoleRequest {
	s.EnableNotification = &v
	return s
}
func (s *CreateJobConsoleRequest) GetDatasetInfo() []map[string]interface{} { return s.DatasetInfo }
func (s *CreateJobConsoleRequest) SetDatasetInfo(v []map[string]interface{}) *CreateJobConsoleRequest {
	s.DatasetInfo = v
	return s
}
func (s *CreateJobConsoleRequest) GetEnvs() []*CreateJobConsoleRequestEnvs { return s.Envs }
func (s *CreateJobConsoleRequest) SetEnvs(v []*CreateJobConsoleRequestEnvs) *CreateJobConsoleRequest {
	s.Envs = v
	return s
}
func (s *CreateJobConsoleRequest) GetSpecifiedNodes() []*string { return s.SpecifiedNodes }
func (s *CreateJobConsoleRequest) SetSpecifiedNodes(v []*string) *CreateJobConsoleRequest {
	s.SpecifiedNodes = v
	return s
}
func (s *CreateJobConsoleRequest) GetExcludeNodes() []*string { return s.ExcludeNodes }
func (s *CreateJobConsoleRequest) SetExcludeNodes(v []*string) *CreateJobConsoleRequest {
	s.ExcludeNodes = v
	return s
}
func (s *CreateJobConsoleRequest) GetMaxRunningTimeMs() *string { return s.MaxRunningTimeMs }
func (s *CreateJobConsoleRequest) SetMaxRunningTimeMs(v string) *CreateJobConsoleRequest {
	s.MaxRunningTimeMs = &v
	return s
}
func (s *CreateJobConsoleRequest) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	for name, value := range map[string]*string{"name": s.Name, "workspace_id": s.WorkspaceId, "project_id": s.ProjectId, "logic_compute_group_id": s.LogicComputeGroupId, "command": s.Command, "framework": s.Framework} {
		if err := model.Required(value, name); err != nil {
			return err
		}
	}
	if s.TaskPriority == nil || *s.TaskPriority < 1 || *s.TaskPriority > 10 || len(s.FrameworkConfig) == 0 {
		return openapi.ValidationError("Explicit task_priority 1..10 and framework_config are required.")
	}
	return nil
}

type ListJobsResponseBodyResult struct {
	model.Extra `json:"-"`
	Jobs        []*Job      `json:"jobs,omitempty"`
	Total       interface{} `json:"total,omitempty"`
}

func (s ListJobsResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s ListJobsResponseBodyResult) GoString() string { return s.String() }
func (s ListJobsResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire ListJobsResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListJobsResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire ListJobsResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListJobsResponseBodyResult) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListJobsResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListJobsResponseBodyResult) GetJobs() []*Job { return s.Jobs }
func (s *ListJobsResponseBodyResult) SetJobs(v []*Job) *ListJobsResponseBodyResult {
	s.Jobs = v
	return s
}
func (s *ListJobsResponseBodyResult) GetTotal() interface{} { return s.Total }
func (s *ListJobsResponseBodyResult) SetTotal(v interface{}) *ListJobsResponseBodyResult {
	s.Total = v
	return s
}
func (s *ListJobsResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.Jobs == nil {
		return openapi.NewError("InvalidResponseError", "InvalidCollection", "Response omitted its collection.")
	}
	return nil
}

type GetJobResponseBodyResult struct {
	model.Extra `json:"-"`
	JobId       *string `json:"job_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (s GetJobResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s GetJobResponseBodyResult) GoString() string { return s.String() }
func (s GetJobResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire GetJobResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetJobResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire GetJobResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetJobResponseBodyResult) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetJobResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetJobResponseBodyResult) GetJobId() *string { return s.JobId }
func (s *GetJobResponseBodyResult) SetJobId(v string) *GetJobResponseBodyResult {
	s.JobId = &v
	return s
}
func (s *GetJobResponseBodyResult) GetName() *string                           { return s.Name }
func (s *GetJobResponseBodyResult) SetName(v string) *GetJobResponseBodyResult { s.Name = &v; return s }
func (s *GetJobResponseBodyResult) GetStatus() *string                         { return s.Status }
func (s *GetJobResponseBodyResult) SetStatus(v string) *GetJobResponseBodyResult {
	s.Status = &v
	return s
}
func (s *GetJobResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type CreateJobConsoleResponseBodyResult struct {
	model.Extra `json:"-"`
	JobId       *string `json:"job_id,omitempty"`
}

func (s CreateJobConsoleResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleResponseBodyResult) GoString() string { return s.String() }
func (s CreateJobConsoleResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleResponseBodyResult) ToMap() (map[string]interface{}, error) {
	return model.ToMap(s)
}
func (s *CreateJobConsoleResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleResponseBodyResult) GetJobId() *string { return s.JobId }
func (s *CreateJobConsoleResponseBodyResult) SetJobId(v string) *CreateJobConsoleResponseBodyResult {
	s.JobId = &v
	return s
}
func (s *CreateJobConsoleResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	if s.JobId == nil || *s.JobId == "" {
		return openapi.NewError("InvalidResponseError", "MissingJobId", "Create response omitted job_id.")
	}
	return nil
}

type StopJobResponseBodyResult struct {
	model.Extra `json:"-"`
}

func (s StopJobResponseBodyResult) String() string   { return dara.Prettify(s) }
func (s StopJobResponseBodyResult) GoString() string { return s.String() }
func (s StopJobResponseBodyResult) MarshalJSON() ([]byte, error) {
	type wire StopJobResponseBodyResult
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *StopJobResponseBodyResult) UnmarshalJSON(data []byte) error {
	type wire StopJobResponseBodyResult
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *StopJobResponseBodyResult) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *StopJobResponseBodyResult) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *StopJobResponseBodyResult) Validate() error {
	if s == nil {
		return openapi.ValidationError("Model is nil.")
	}
	if err := model.ValidateChildren(s); err != nil {
		return err
	}
	return nil
}

type ListJobsResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata   `json:"ResponseMetadata,omitempty"`
	Result           *ListJobsResponseBodyResult `json:"Result,omitempty"`
}

func (s ListJobsResponseBody) String() string   { return dara.Prettify(s) }
func (s ListJobsResponseBody) GoString() string { return s.String() }
func (s ListJobsResponseBody) MarshalJSON() ([]byte, error) {
	type wire ListJobsResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListJobsResponseBody) UnmarshalJSON(data []byte) error {
	type wire ListJobsResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListJobsResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListJobsResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListJobsResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *ListJobsResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *ListJobsResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *ListJobsResponseBody) GetResult() *ListJobsResponseBodyResult { return s.Result }
func (s *ListJobsResponseBody) SetResult(v *ListJobsResponseBodyResult) *ListJobsResponseBody {
	s.Result = v
	return s
}
func (s *ListJobsResponseBody) Validate() error {
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

type ListJobsResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string       `json:"headers,omitempty"`
	StatusCode  *int32                   `json:"statusCode,omitempty"`
	Body        *ListJobsResponseBody    `json:"body,omitempty"`
	RawBody     []byte                   `json:"-"`
	SdkContext  *openapi.ResponseContext `json:"-"`
}

func (s ListJobsResponse) String() string   { return dara.Prettify(s) }
func (s ListJobsResponse) GoString() string { return s.String() }
func (s ListJobsResponse) MarshalJSON() ([]byte, error) {
	type wire ListJobsResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *ListJobsResponse) UnmarshalJSON(data []byte) error {
	type wire ListJobsResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *ListJobsResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *ListJobsResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *ListJobsResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *ListJobsResponse) SetHeaders(v map[string]*string) *ListJobsResponse {
	s.Headers = v
	return s
}
func (s *ListJobsResponse) GetStatusCode() *int32                             { return s.StatusCode }
func (s *ListJobsResponse) SetStatusCode(v int32) *ListJobsResponse           { s.StatusCode = &v; return s }
func (s *ListJobsResponse) GetBody() *ListJobsResponseBody                    { return s.Body }
func (s *ListJobsResponse) SetBody(v *ListJobsResponseBody) *ListJobsResponse { s.Body = v; return s }
func (s *ListJobsResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *ListJobsResponse) Validate() error {
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

type GetJobResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata `json:"ResponseMetadata,omitempty"`
	Result           *GetJobResponseBodyResult `json:"Result,omitempty"`
}

func (s GetJobResponseBody) String() string   { return dara.Prettify(s) }
func (s GetJobResponseBody) GoString() string { return s.String() }
func (s GetJobResponseBody) MarshalJSON() ([]byte, error) {
	type wire GetJobResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetJobResponseBody) UnmarshalJSON(data []byte) error {
	type wire GetJobResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetJobResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *GetJobResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *GetJobResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *GetJobResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *GetJobResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *GetJobResponseBody) GetResult() *GetJobResponseBodyResult { return s.Result }
func (s *GetJobResponseBody) SetResult(v *GetJobResponseBodyResult) *GetJobResponseBody {
	s.Result = v
	return s
}
func (s *GetJobResponseBody) Validate() error {
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

type GetJobResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string       `json:"headers,omitempty"`
	StatusCode  *int32                   `json:"statusCode,omitempty"`
	Body        *GetJobResponseBody      `json:"body,omitempty"`
	RawBody     []byte                   `json:"-"`
	SdkContext  *openapi.ResponseContext `json:"-"`
}

func (s GetJobResponse) String() string   { return dara.Prettify(s) }
func (s GetJobResponse) GoString() string { return s.String() }
func (s GetJobResponse) MarshalJSON() ([]byte, error) {
	type wire GetJobResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *GetJobResponse) UnmarshalJSON(data []byte) error {
	type wire GetJobResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *GetJobResponse) ToMap() (map[string]interface{}, error)          { return model.ToMap(s) }
func (s *GetJobResponse) FromMap(value map[string]interface{}) error      { return model.FromMap(value, s) }
func (s *GetJobResponse) GetHeaders() map[string]*string                  { return s.Headers }
func (s *GetJobResponse) SetHeaders(v map[string]*string) *GetJobResponse { s.Headers = v; return s }
func (s *GetJobResponse) GetStatusCode() *int32                           { return s.StatusCode }
func (s *GetJobResponse) SetStatusCode(v int32) *GetJobResponse           { s.StatusCode = &v; return s }
func (s *GetJobResponse) GetBody() *GetJobResponseBody                    { return s.Body }
func (s *GetJobResponse) SetBody(v *GetJobResponseBody) *GetJobResponse   { s.Body = v; return s }
func (s *GetJobResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *GetJobResponse) Validate() error {
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

type CreateJobConsoleResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata           `json:"ResponseMetadata,omitempty"`
	Result           *CreateJobConsoleResponseBodyResult `json:"Result,omitempty"`
}

func (s CreateJobConsoleResponseBody) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleResponseBody) GoString() string { return s.String() }
func (s CreateJobConsoleResponseBody) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleResponseBody) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *CreateJobConsoleResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *CreateJobConsoleResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *CreateJobConsoleResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *CreateJobConsoleResponseBody) GetResult() *CreateJobConsoleResponseBodyResult {
	return s.Result
}
func (s *CreateJobConsoleResponseBody) SetResult(v *CreateJobConsoleResponseBodyResult) *CreateJobConsoleResponseBody {
	s.Result = v
	return s
}
func (s *CreateJobConsoleResponseBody) Validate() error {
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

type CreateJobConsoleResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string            `json:"headers,omitempty"`
	StatusCode  *int32                        `json:"statusCode,omitempty"`
	Body        *CreateJobConsoleResponseBody `json:"body,omitempty"`
	RawBody     []byte                        `json:"-"`
	SdkContext  *openapi.ResponseContext      `json:"-"`
}

func (s CreateJobConsoleResponse) String() string   { return dara.Prettify(s) }
func (s CreateJobConsoleResponse) GoString() string { return s.String() }
func (s CreateJobConsoleResponse) MarshalJSON() ([]byte, error) {
	type wire CreateJobConsoleResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *CreateJobConsoleResponse) UnmarshalJSON(data []byte) error {
	type wire CreateJobConsoleResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *CreateJobConsoleResponse) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *CreateJobConsoleResponse) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *CreateJobConsoleResponse) GetHeaders() map[string]*string { return s.Headers }
func (s *CreateJobConsoleResponse) SetHeaders(v map[string]*string) *CreateJobConsoleResponse {
	s.Headers = v
	return s
}
func (s *CreateJobConsoleResponse) GetStatusCode() *int32 { return s.StatusCode }
func (s *CreateJobConsoleResponse) SetStatusCode(v int32) *CreateJobConsoleResponse {
	s.StatusCode = &v
	return s
}
func (s *CreateJobConsoleResponse) GetBody() *CreateJobConsoleResponseBody { return s.Body }
func (s *CreateJobConsoleResponse) SetBody(v *CreateJobConsoleResponseBody) *CreateJobConsoleResponse {
	s.Body = v
	return s
}
func (s *CreateJobConsoleResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *CreateJobConsoleResponse) Validate() error {
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

type StopJobResponseBody struct {
	model.Extra      `json:"-"`
	ResponseMetadata *openapi.ResponseMetadata  `json:"ResponseMetadata,omitempty"`
	Result           *StopJobResponseBodyResult `json:"Result,omitempty"`
}

func (s StopJobResponseBody) String() string   { return dara.Prettify(s) }
func (s StopJobResponseBody) GoString() string { return s.String() }
func (s StopJobResponseBody) MarshalJSON() ([]byte, error) {
	type wire StopJobResponseBody
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *StopJobResponseBody) UnmarshalJSON(data []byte) error {
	type wire StopJobResponseBody
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *StopJobResponseBody) ToMap() (map[string]interface{}, error) { return model.ToMap(s) }
func (s *StopJobResponseBody) FromMap(value map[string]interface{}) error {
	return model.FromMap(value, s)
}
func (s *StopJobResponseBody) GetResponseMetadata() *openapi.ResponseMetadata {
	return s.ResponseMetadata
}
func (s *StopJobResponseBody) SetResponseMetadata(v *openapi.ResponseMetadata) *StopJobResponseBody {
	s.ResponseMetadata = v
	return s
}
func (s *StopJobResponseBody) GetResult() *StopJobResponseBodyResult { return s.Result }
func (s *StopJobResponseBody) SetResult(v *StopJobResponseBodyResult) *StopJobResponseBody {
	s.Result = v
	return s
}
func (s *StopJobResponseBody) Validate() error {
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

type StopJobResponse struct {
	model.Extra `json:"-"`
	Headers     map[string]*string       `json:"headers,omitempty"`
	StatusCode  *int32                   `json:"statusCode,omitempty"`
	Body        *StopJobResponseBody     `json:"body,omitempty"`
	RawBody     []byte                   `json:"-"`
	SdkContext  *openapi.ResponseContext `json:"-"`
}

func (s StopJobResponse) String() string   { return dara.Prettify(s) }
func (s StopJobResponse) GoString() string { return s.String() }
func (s StopJobResponse) MarshalJSON() ([]byte, error) {
	type wire StopJobResponse
	return model.Marshal(wire(s), s.ExtraFields)
}
func (s *StopJobResponse) UnmarshalJSON(data []byte) error {
	type wire StopJobResponse
	return model.Unmarshal(data, (*wire)(s), &s.ExtraFields)
}
func (s *StopJobResponse) ToMap() (map[string]interface{}, error)           { return model.ToMap(s) }
func (s *StopJobResponse) FromMap(value map[string]interface{}) error       { return model.FromMap(value, s) }
func (s *StopJobResponse) GetHeaders() map[string]*string                   { return s.Headers }
func (s *StopJobResponse) SetHeaders(v map[string]*string) *StopJobResponse { s.Headers = v; return s }
func (s *StopJobResponse) GetStatusCode() *int32                            { return s.StatusCode }
func (s *StopJobResponse) SetStatusCode(v int32) *StopJobResponse           { s.StatusCode = &v; return s }
func (s *StopJobResponse) GetBody() *StopJobResponseBody                    { return s.Body }
func (s *StopJobResponse) SetBody(v *StopJobResponseBody) *StopJobResponse  { s.Body = v; return s }
func (s *StopJobResponse) SetDiagnostics(raw []byte, ctx *openapi.ResponseContext) {
	s.RawBody = raw
	s.SdkContext = ctx
}
func (s *StopJobResponse) Validate() error {
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

var _ json.Marshaler = (*Job)(nil)
