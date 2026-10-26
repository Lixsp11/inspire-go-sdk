# 查询

## 当前账号

```go
client, err := user.NewClient(cfg)
if err != nil { return err }
defer client.Close()
response, err := client.GetUserDetailWithContext(ctx, &user.GetUserDetailRequest{}, nil)
if err != nil { return err }
account, err := response.Body.Result.ToMap()
if err != nil { return err }
```

user别名对应`github.com/Lixsp11/inspire-go-sdk/src/userv2/client`。用户ID取Result.Id或UserId，不能用显示名称填created_by。姓名、英文名、邮箱、角色等从实际map返回；空值如实为空，不推测。created_at换算时先核实单位并标时区。

## 各类资源

| 查询 | Request及必要字段 | Result集合 |
|---|---|---|
| 项目 | ListProjectsRequest；可选Filter.WorkspaceId | Items |
| 计算组 | ListLogicComputeGroupsRequest；Filter.WorkspaceId，显式PageSize | LogicComputeGroups |
| 训练任务 | ListJobsRequest；WorkspaceId；本人筛选CreatedBy=user ID | Jobs |
| 任务详情 | GetJobRequest.JobId | 单个任务 |
| 开发机 | ListNotebooksRequest；WorkspaceId、FilterBy.UserId列表 | List |
| 开发机详情 | GetNotebookRequest.NotebookId | 单个开发机 |

项目可省略Filter做全局查询，不臆造尚未封装的ListWorkspaces。项目原生额外字段可能含workspace关联，从ToMap核实，缺失时再查已知合同或补问。

计算组PageSize必须明确，可用-1取全部；遗漏该字段可能导致误导性的空集合。计算组名称不能直接说明GPU余量、quota或镜像适用性。当前没有类型化的完整配额/镜像目录发现接口。

## 训练分页

```go
request := (&train.ListJobsRequest{}).
    SetWorkspaceId(workspaceID).SetCreatedBy(userID).SetPageSize(50)
pager, err := train.NewListJobsPaginator(client, request, runtime)
if err != nil { return err }
for pager.HasMorePages() {
    page, err := pager.NextPage(ctx)
    if err != nil { return err }
    // page.Body.Result.Jobs 是当前页，page.SdkContext 保留实际诊断信息。
}
total, present := pager.ReportedTotal()
lastPage, pagesRead := pager.LastPageNum(), pager.PagesRead()
```

train别名对应`github.com/Lixsp11/inspire-go-sdk/src/trainv2/client`。分页器顺序使用，保留实际请求页号和平台total，检测重复ID/页、total变化和截断。报错后只能报告部分结果，不能称完整清单。需要限制页数时在首次调用前SetMaxPages。

ListJobs使用PageNum/PageSize，上限999。项目和Notebook使用Page/PageSize；项目/计算组允许-1。其他接口不照搬999或5000上限。Notebook筛选使用FilterBy.UserId列表，分页需逐页对照实际Total；不要把返回条数冒充平台总数。

GetDiscovery/WithOptions/WithContext是公共GET，无请求模型；返回Services/Version仅帮助定位接口，不证明所有列出的动作都已封装或可安全执行。
