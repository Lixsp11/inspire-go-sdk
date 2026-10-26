# SDK 用法

## 依赖与配置

module：`github.com/Lixsp11/inspire-go-sdk`。Go >=1.23；服务包位于 `src/` 下，例如 `src/openapi/client` 和 `src/trainv2/client`。RuntimeOptions 直接使用 `github.com/alibabacloud-go/tea/dara`。

```go
import (
    "os"
    openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
)

cfg := (&openapi.Config{}).
    SetEndpoint("qz.sii.edu.cn").
    SetSocks5Proxy(os.Getenv("INSPIRE_PROXY")).
    SetConnectTimeout(5000).
    SetReadTimeout(10000).
    SetOperationTimeout(60000)
```

需要代理时先核实进程中的 INSPIRE_PROXY 已设置；不要打印其值。SDK 不读取系统代理。SOCKS5h 由代理解析目标域名，NoProxy 使用 x/net 规则。Ca/Cert/Key 是 PEM 内容，TLS 默认校验。

已有会话通过 INSPIRE_COOKIE 提供，或显式配置 Cookie/provider。环境中 Cookie 和访问token同时存在会报歧义错误。刷新 provider 只负责加载已有管理材料，不等于实现了CAS登录。

字段优先级为 Runtime 非nil值、Config、默认；false/空串明确覆盖，nil才继承。header/query 按全局、Runtime.ExtendsParameters、请求层依次覆盖；provider最终设置管理认证头。Action、帧边界、代理认证和幂等头有保护。

## 返回与模型

返回为 Response.Headers/StatusCode/Body，原生字段位于 Body.ResponseMetadata 和 Body.Result。例：创建 ID 为 response.Body.Result.JobId。指针值可用 dara.StringValue 等读取。非类型化字段通过 result.ToMap() 读取，处理其 error。

结构可空构造、赋值或链式Setter，但发送前执行Validate。FromMap/JSON采用Go合并语义：缺字段保留，null清指针；未知字段以RawMessage保留。非nil空集合会保留，nil省略。RawBody和整个业务模型可能含敏感数据，诊断优先使用SdkContext.RequestId/Attempts。

## 有限只读重试

现代策略位于Config.RetryOptions，当前Go RuntimeOptions没有现代RetryOptions字段：

```go
cfg.SetRetryOptions(&dara.RetryOptions{
    Retryable: true, MaxAttempts: 3,
    RetryCondition: []*dara.RetryCondition{{
        MaxAttempts: 3, MaxDelay: 120000,
        Exception: []string{"TransportError", "ThrottlingError", "ServerError"},
        Backoff: &dara.FixedBackoffPolicy{Period: 1000},
    }},
})
```

两处MaxAttempts都设3以表达每个只读API最多三次发送，包含首次。条件匹配时用条件上限，不匹配才用顶层上限；0按上游默认3。退避1000ms，Retry-After是等待下限。整段查询通过外层context限制总时间；手动重试和同步/并发切换也计入约定的请求预算。

旧Runtime.Autoretry/MaxAttempts/Backoff*是兼容入口，不与Config现代策略同时指定。写请求不因任何重试设置自动重放。

| 错误Name | 判断 |
|---|---|
| TransportError | 未收到可信平台结果，不能据此判定Cookie失效；只读请求按预算重试 |
| AuthenticationError | 会话不被接受，复用已有刷新能力或有效管理材料 |
| PermissionDeniedError | 核实账号/资源权限，不能报告为空集合 |
| TlsError | 核对证书、目标和代理，不以关闭校验代替修复 |
| ValidationError / ConfigurationError | 修复输入，通常尚未发送 |
| InvalidResponseError | HTTP成功也不能证明业务结构正确 |
| MutationUncertainError | 先核对资源是否已创建或改变 |

错误可errors.As到*openapi.SDKError，Name/Code/StatusCode/RequestId/Attempts提供安全诊断；errors.Is保留context取消/超时。未封装接口使用CallApi/CallApiWithContext，只读标志必须来自真实动作语义。
