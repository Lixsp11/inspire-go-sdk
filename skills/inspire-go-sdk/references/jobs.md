# 训练任务

使用`src/trainv2/client`的CreateJobConsoleRequest，调用CreateJobConsole/WithOptions/WithContext。创建需要实际workspace_id、project_id、logic_compute_group_id、quota_id、镜像、框架、命令和明确优先级；不能用discovery里的CreateJob替换合同。

| 字段 | 对应关系 |
|---|---|
| FrameworkConfig | 至少一项CreateJobConsoleRequestFrameworkConfig |
| Cpu / GpuCount / MemGi / ShmGi | 每实例资源；共享内存不能超过内存 |
| InstanceCount | 本规格实例数；总需求乘每实例资源 |
| Image / ImageType | 已核实的地址和来源，不用显示名猜测 |
| ResourceSpecPrice | 使用嵌套CreateJobConsoleRequestFrameworkConfigResourceSpecPrice |
| ResourceSpecPrice.LogicComputeGroupId / QuotaId | 与实际授权资源匹配的ID |
| ResourceSpecPrice.GpuType | 有GPU时必填机器类型，不能用展示名称代替 |
| TaskPriority | 显式1..10；空间/项目可能进一步限制 |
| MaxRunningTimeMs | 字符串形式毫秒数，由当前需求决定 |

标准空间与公平空间优先级规则可能不同，不能统一默认10。当前未提供完整类型化的镜像、quota和数据集目录；先核实已有配置或原生合同，不为让模型通过而编造材料。

创建前构造模型并Validate，说明字段来自哪个实际资源。构造/校验与提交分开，提交只按已有授权执行。

得到Body.Result.JobId仅表示有创建确认。再GetJob核对当前状态、关联资源和关键配置；排队、拉镜像不代表已运行。TaskPriority和所选资源最终由服务端核验。

## 创建不确定

遇到MutationUncertainError或发送期间取消，保留提交时间、提交者、空间、名称、请求配置和安全request ID。先ListJobs，再GetJob核对候选；名称不唯一，单个同名/配置相似对象也不能自动证明是这次提交。无法唯一确认时报告不确定性，不靠重建来确认。

## 停止与等待

StopJobRequest.JobId应来自当前明确授权的目标。先GetJob核对对象和状态，已满足用户目标则跳过重复写入；受理后继续查询，空Result不是已经释放资源的证明。

`client.WaitForJob(ctx, jobID, &train.WaiterOptions{Timeout: 10*time.Minute, Interval: 5*time.Second})`等待默认成功终态SUCCEEDED/job_succeeded；FAILED/CANCELLED/STOPPED等属于失败。需要等待运行或驻留时按真实状态配置TargetStates/FailureStates，不能用等待退出的默认值替换用户意图。

等待超时不停止任务。当前没有封装更新、删除、Pod日志/事件或容器执行接口，额外需求先核实合同与支持范围。
