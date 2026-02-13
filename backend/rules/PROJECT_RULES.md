# 项目开发规范 (Project Rules)

本文档定义 go-mail 电商微服务项目的开发规范。**所有代码生成和编写必须严格遵循。**

---

## 1. 全局规范

### 1.1 语言要求 (强制)

- **AI 回复、代码注释、文档内容、Git Commit Message** 必须使用 **中文 (简体)**。
- 变量名、函数名等代码标识符保持英文，注释为中文。

### 1.2 文件命名

| 语言 | 格式 | 示例 |
|:---|:---|:---|
| Go 文件 | `snake_case.go` | `user_model.go` |
| React 组件 | `PascalCase.tsx` | `UserProfile.tsx` |
| TS/JS 文件 | `camelCase.ts` | `apiClient.ts` |
| API 定义 | `service.api` | `user.api` |
| Proto 定义 | `service.proto` | `user.proto` |

### 1.3 代码风格

- 遵循 DRY 原则，删除未使用的导入和变量。
- Go 文件使用 `gofmt` / `goimports` 格式化。

---

## 2. go-zero 开发规范 (核心)

> **最高优先级**：编写 go-zero 相关代码时，**必须严格遵守** 以下知识库：
>
> 1. `.cursorrules/` — AI 工作流指令（决策树、工具用法、代码模式）
> 2. `.ai-context/zero-skills/` — 完整模式参考（REST API / RPC / Database / Resilience）

### 2.1 Spec-First 原则 (强制)

**先定义 `.api` / `.proto` 文件，再用 `goctl` 生成代码。禁止手写 handler、routes、types 等自动生成文件。**

```
工作流：定义 .api/.proto → goctl 生成 → 编写 Logic 层业务逻辑
```

### 2.2 .api 文件规范

遵循 `.cursorrules/patterns.md` 和 `zero-skills/references/rest-api-patterns.md`：

```api
syntax = "v1"

info (
    title:   "服务名称"
    desc:    "服务描述"
    author:  "go-mail"
    version: "v1"
)

type (
    CreateReq {
        Name string `json:"name" validate:"required,min=2"`
    }
    CreateResp {
        Id int64 `json:"id"`
    }
)

// 公开路由
@server (
    prefix: /api/v1
    group:  moduleName
)
service service-api {
    @doc "接口描述"
    @handler HandlerName
    post /path (CreateReq) returns (CreateResp)
}

// JWT 保护路由
@server (
    prefix: /api/v1
    group:  moduleName
    jwt:    Auth
)
service service-api {
    @doc "需要认证的接口"
    @handler ProtectedHandler
    get /protected/path returns (Response)
}
```

**必须遵守的规则**：
- `syntax = "v1"` + `info()` 头部声明
- `@server(prefix/group/jwt)` 路由分组
- `validate` 标签进行参数校验
- `@doc` 文档注释 + `@handler` 命名
- `json:"field_name"` / `path:"id"` / `form:"param"` 标签

### 2.3 .proto 文件规范

遵循 `zero-skills/references/rpc-patterns.md`：

```protobuf
syntax = "proto3";

package servicename;
option go_package = "./servicename";

service ServiceRpc {
  rpc MethodName(MethodReq) returns (MethodResp);
}

message MethodReq {
  int64 id = 1;
}

message MethodResp {
  int64  id   = 1;
  string name = 2;
}
```

### 2.4 goctl 代码生成命令

```bash
# API 服务
goctl api go -api desc/service.api -dir . -style gozero

# RPC 服务
goctl rpc protoc service.proto --go_out=. --go-grpc_out=. --zrpc_out=. -style gozero

# PostgreSQL Model
goctl model pg datasource -url="postgres://..." -table="table_name" -dir ./model -c -style gozero
```

### 2.5 三层架构 (强制)

遵循 `zero-skills/SKILL.md` 的 Handler → Logic → Model 三层分离：

| 层 | 职责 | 是否可手动修改 |
|:---|:---|:---|
| **Handler** | HTTP 路由、请求解析、响应输出 | ❌ goctl 生成，不手动修改 |
| **Logic** | 业务逻辑实现 | ✅ 这里写代码 |
| **Model** | 数据库操作 | ⚠️ 自定义方法在 `customXxxModel` 中扩展 |
| **Config** | 配置结构定义 | ✅ 可添加自定义字段 |
| **ServiceContext** | 依赖注入容器 | ✅ 注入 Model/RPC Client |

### 2.6 错误处理

参考 `.cursorrules/patterns.md`：

```go
// ✅ 正确：使用 xerr 包
return nil, xerr.NewCodeError(xerr.BadRequest)
return nil, xerr.NewCodeErrorMsg(xerr.Unauthorized, "用户名或密码错误")

// ✅ RPC 层使用 gRPC status
return nil, status.Error(codes.NotFound, "用户不存在")

// ❌ 错误：不要使用 fmt.Errorf 返回 API 错误
return nil, fmt.Errorf("用户不存在")
```

### 2.7 配置模式

```go
// Config 结构
type Config struct {
    rest.RestConf
    DataSource string           // PostgreSQL 连接串
    Cache      cache.CacheConf  // Redis 缓存
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
}
```

```yaml
# YAML 配置
Name: service-api
Host: 0.0.0.0
Port: 8801
DataSource: "postgres://user:pass@host:5432/db?sslmode=disable"
Cache:
  - Host: localhost:6379
Auth:
  AccessSecret: "secret"
  AccessExpire: 86400
```

### 2.8 JWT 认证

```go
// 生成 Token
claims := jwt.MapClaims{"userId": user.Id, "exp": now + expire}
token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
    SignedString([]byte(secret))

// 从 context 提取 userId
userId := l.ctx.Value("userId").(json.Number)
uid, _ := userId.Int64()
```

### 2.9 Handler 统一响应

所有 Handler 使用 `common/response` 包返回统一格式 `{code, msg, data}`：

```go
// ✅ 正确
response.Response(w, resp, err)

// ❌ 错误
httpx.OkJsonCtx(r.Context(), w, resp)
```

---

## 3. 开发工作流

### 3.1 新增 HTTP API

```
1. 编写/修改 desc/*.api 文件
2. goctl api go -api desc/main.api -dir . -style gozero
3. 在 internal/logic/ 实现业务逻辑
4. 修改 handler 使用统一响应
```

### 3.2 新增 RPC 服务

```
1. 编写 *.proto 文件
2. goctl rpc protoc *.proto --go_out=. --go-grpc_out=. --zrpc_out=. -style gozero
3. 在 internal/logic/ 实现 RPC 方法
```

### 3.3 数据库 Model

```
1. 编写 DDL SQL (deploy/sql/)
2. goctl model pg datasource -url="..." -table="xxx" -dir model -c -style gozero
3. 在 ServiceContext 中注入 Model
4. 在 Logic 中使用
```

---

## 4. 数据库规范

- 使用 **PostgreSQL 16** 作为主数据库。
- 数据类型：`BIGSERIAL` (主键)、`VARCHAR` (字符串)、`NUMERIC(10,2)` (金额)、`TIMESTAMPTZ` (时间)、`JSONB` (结构化数据)、`BOOLEAN` (布尔)。
- SQL 脚本放在 `deploy/sql/` 目录。

---

## 5. 禁止事项

参考 `.cursorrules/00-instructions.md` 和 `zero-skills/SKILL.md`：

| ❌ 禁止 | ✅ 应该 |
|:---|:---|
| 手写 handler / routes / types | 用 goctl 生成 |
| 不写 .api/.proto 直接写代码 | 先 Spec 再生成 |
| `fmt.Errorf` 返回 API 错误 | 用 `xerr.NewCodeError` |
| 硬编码配置值 | 用 Config + YAML |
| 跳过 validate 校验 | 添加 `validate` 标签 |
| 在 Handler 层写业务逻辑 | 逻辑放 Logic 层 |
| 修改 goctl 自动生成的代码 | 通过 Logic / customModel 扩展 |
| 空 stub (`// todo`) | 生成完整实现 |

---

## 6. Git 提交规范

使用 Angular 规范，中文描述：

```
feat: 添加用户登录接口
fix: 修复登录页 token 验证失败
docs: 更新 API 设计文档
refactor: 重构用户服务配置
style: 格式化代码
test: 添加注册逻辑单元测试
chore: 更新 go.mod 依赖
```

---

## 7. 知识库参考索引

编写代码时必须查阅以下知识库（按优先级排序）：

| 优先级 | 文件 | 内容 |
|:---|:---|:---|
| 1 | `.cursorrules/workflows.md` | 工作流步骤 |
| 2 | `.cursorrules/tools.md` | MCP 工具用法 |
| 3 | `.cursorrules/patterns.md` | 代码模式速查 |
| 4 | `.ai-context/zero-skills/references/rest-api-patterns.md` | REST API 详细模式 |
| 5 | `.ai-context/zero-skills/references/rpc-patterns.md` | RPC 详细模式 |
| 6 | `.ai-context/zero-skills/references/database-patterns.md` | 数据库详细模式 |
| 7 | `.ai-context/zero-skills/references/resilience-patterns.md` | 弹性模式 |
| 8 | `.ai-context/zero-skills/troubleshooting/common-issues.md` | 常见问题排查 |
