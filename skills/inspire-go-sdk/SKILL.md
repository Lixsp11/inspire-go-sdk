---
name: inspire-go-sdk
description: 使用 Inspire Go SDK 调用启智 qz.sii.edu.cn 原生 API，编写账号和资源查询、训练任务创建与停止、分页、状态等待代码，排查认证、代理和返回值。适用于 github.com/Lixsp11/inspire-go-sdk；不涵盖容器命令、SSH/SCP、隧道或模型推理接口。
---

# Inspire Go SDK

使用现有 Go SDK 完成平台调用；本资料适用于 SDK 0.1.x。按需读取：

| 任务 | 参考 |
|---|---|
| module、配置、代理、超时和重试 | [SDK 用法](references/sdk-usage.md) |
| 账号、项目、计算组、任务和开发机查询 | [查询](references/queries.md) |
| 创建、停止与等待训练任务 | [任务操作](references/jobs.md) |

## 项目与接口

检查目标项目 go.mod、Go 版本和已有调用方式，复用现有依赖管理。module 为 `github.com/Lixsp11/inspire-go-sdk`，服务包为 `src/trainv2/client` 等；优先采用类型化 Client 和指针模型。核对当前版本的字段、Setter、Validate 和方法签名，不能用名称相似的云 API 代替原生合同。

普通方法和 WithOptions 之外，WithContext 接收 context.Context；Go 并发使用 goroutine。超时、退避的配置单位是毫秒，WaiterOptions 的 Timeout/Interval 则为 time.Duration。所有客户端最终 Close；调用方提供的 Transport 仍由调用方管理。

## 连接与凭据

沿用当前已授权的管理 Cookie、provider 和代理。Cookie 为完整 `inspire-session=...`，SDK 环境入口是 INSPIRE_COOKIE；推理 key 不是管理登录。Config.Cookie、BearerToken、Credential 三选一。

用户要求代理时，明确读取已有代理配置并设置 Config.Socks5Proxy、HttpProxy 或 HttpsProxy；INSPIRE_PROXY 是示例约定，不是 SDK 自动读取的设置。不得因网络失败自行改为直连。凭据和代理密码不写源码、请求日志或公开结果。

## 执行与报告

读取现有信息确定空间、项目、用户、配额、镜像和目标对象，只补问无法从上下文或查询确定的内容。查询按原生筛选和分页执行，保留真实 total、访问页范围及实际数据；空结果、权限错误、会话失效和网络失败分别报告。账号查询返回所需身份字段，不以单独的 HTTP 200 代替结果。

默认查询单发；暂时网络故障可在明确次数和总预算内重试只读接口。401/403、参数和永久TLS错误不循环用相同材料重发。每个 API 的尝试数不是整段多页流程的总请求数，完整流程应有外层 context。

写操作按当前任务及既有授权执行，校验代码不自动授权提交资源。MutationUncertainError 或提交期间取消时先查询核对，不能盲目再创建；名称相同不证明是同一提交。请求受理、开始运行和最终成功分别通过状态确认。

先做本地模型/合成传输验证；真实冒烟查询仅使用当前授权的只读范围。生成构建包不安装 skill，也不读取或修改真实平台资源。未封装 API 先核实路径、Action、JSON 结构及副作用，再使用 CallApi；discovery 不能单独证明创建合同。
