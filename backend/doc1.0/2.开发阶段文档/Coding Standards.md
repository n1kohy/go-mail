
---

# 代码规范文档 (Coding Standards)

**项目名称**：基于 Gin 的微服务电商平台
**版本**：1.0
**日期**：2026-02-11
**适用范围**：后端（Go/Gin/Go-Zero）、前端（React）、数据库（MySQL）

---

## 1. 引言 (Introduction)

本规范旨在统一开发团队的代码风格，降低维护成本，提高代码的可读性与健壮性。所有提交到代码仓库的代码必须遵循本规范。项目采用 Go 语言作为后端核心 ，React 作为前端框架 ，并严格遵循微服务架构原则。

---

## 2. Go 语言后端规范 (Backend Standards)

### 2.1 格式化与排版

* **工具强制**：所有 Go 代码必须在提交前使用 `gofmt` 或 `goimports` 进行格式化。
* **缩进**：使用 **Tab** 缩进，而非空格（Go 标准）。
* **行长**：单行代码建议不超过 120 字符，过长应折行。
* **括号**：控制结构（`if`, `for`）的左大括号 `{` 必须在同一行，不可换行。

### 2.2 命名规范 (Naming Conventions)

遵循 `CamelCase`（驼峰命名法）。

* **包名 (Packages)**：
* 使用全小写单词，无下划线（如 `user`, `order`）。
* 包名应与目录名保持一致。


### 2.2 命名规范 (Naming Conventions - Zero-Skills Standard)

遵循 `CamelCase`（驼峰命名法）。

* **包名 (Packages)**：
    * 使用全小写单词，无下划线（如 `user`, `order`）。
    * 包名应与目录名保持一致。

* **文件命名 (遵循 go-zero 最佳实践)**：
    * **Handler**: `<resource><action>handler.go` (e.g., `createuserhandler.go`)
    * **Logic**: `<resource><action>logic.go` (e.g., `createuserlogic.go`)
    * **Model**: `<table>model.go` (e.g., `usermodel.go`)
    * **Middleware**: `<purpose>middleware.go` (e.g., `authmiddleware.go`)
    * 测试文件必须以 `_test.go` 结尾。

* **变量与常量**：
    * **局部变量**：使用小驼峰（camelCase），如 `userCount`。
    * **导出变量/常量**（Public）：使用大驼峰（PascalCase），如 `MaxRetryCount`。
    * **非导出变量/常量**（Private）：使用小驼峰，如 `defaultTimeout`。

### 2.3 错误处理 (Error Handling - Zero-Skills Standard)

* **定义错误**：在包级别定义通用错误。
    ```go
    var (
        ErrUserNotFound = errors.New("user not found")
        ErrInvalidInput = errors.New("invalid input")
    )
    ```
* **错误包装**：使用 `fmt.Errorf("...: %w", err)` 保留上下文。
* **统一响应**：使用 `httpx.SetErrorHandler` 统一处理错误响应，将其映射为标准 JSON 格式。

### 2.4 日志规范 (Logging - Zero-Skills Standard)

* **严禁使用 `fmt.Println` 或 `log.Println`**。
* **必须使用 `logx`**：
    ```go
    l.Logger.Infof("creating user: %s", req.Email)
    l.Logger.Errorf("failed to create user: %v", err)
    ```
* **敏感信息脱敏**：严禁打印密码、Token、银行卡号等敏感信息。

---

## 3. 框架特定规范 (Framework Specifics)

### 3.1 Project Structure (Go-Zero Standard)

```
service-name/
├── etc/                # 配置定义 (config.yaml)
├── internal/
│   ├── config/         # 配置结构体 (config.go)
│   ├── handler/        # HTTP 处理器 (仅处理请求解析与响应)
│   ├── logic/          # 业务逻辑 (核心业务代码)
│   ├── middleware/     # 中间件
│   ├── svc/            # 服务上下文 (依赖注入)
│   ├── types/          # 请求/响应类型定义
│   └── model/          # 数据库模型
├── service.go          # 程序入口
└── service.api         # API 定义文件
```

### 3.2 配置管理 (Configuration)
使用 `config.Config` 结构体映射 YAML 配置：
```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    Database struct {
        DataSource string
    }
}
```

---

## 4. 数据库规范 (Database Standards)

基于 MySQL 和 Redis 的使用规范。

### 4.1 MySQL 规范
* **表名/列名**：使用 `snake_case`。
* **主键**：必须拥有 `id` (BIGINT)。
* **必须字段**：`create_time`, `update_time`。
* **SQL 编写**：优先使用 `goctl model` 生成的方法，禁止在循环中查询数据库 (N+1 问题)。

### 4.2 Redis 规范
* **Key 命名**：`App:Service:Business:ID` (e.g., `mall:user:cache:1001`)。
* **TTL**：所有 Key 必须设置过期时间。

---

## 5. 前端规范 (Frontend - React)

* **Functional Components**: 必须使用函数式组件 + Hooks。
* **命名**:
    * 组件文件: `UserProfile.tsx` (PascalCase)
    * 工具/Hook: `useAuth.ts` (camelCase)
* **API 请求**: 封装在 `services/` 目录，禁止在组件中直接调用 `axios`。

---

## 6. 版本控制规范 (Git)

* **Commit Message**: 遵循 Angular 规范 (`feat`, `fix`, `docs`, `chore`)。
* **分支管理**: `main` (稳定), `dev` (开发), `feature/xxx` (特性)。