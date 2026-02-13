# 开发工作流 (Workflows)

> 以下工作流严格遵循 `.cursorrules/workflows.md` 和 `.ai-context/zero-skills/SKILL.md` 规范。

---

## 1. 新增 API 服务

```
1. 在 services/<name>/api/desc/ 创建 .api 文件
   - syntax = "v1" + info() 头部
   - 定义 type（含 validate 标签）
   - 定义 @server(prefix/group/jwt) 路由
   - 每个接口写 @doc 和 @handler
2. 运行 goctl 生成代码
   goctl api go -api desc/<name>.api -dir . -style gozero
3. 修改 Config，添加 DataSource / Cache / Auth
4. 修改 ServiceContext，注入 Model / RPC Client
5. 在 internal/logic/ 实现业务逻辑（完整实现，不要空 stub）
6. 修改 Handler 使用统一响应 response.Response()
7. go build ./... 验证编译
```

## 2. 新增 RPC 服务

```
1. 在 services/<name>/rpc/ 创建 .proto 文件
   - syntax = "proto3"
   - 定义 package 和 option go_package
   - 定义 message 和 service
2. 运行 goctl 生成代码
   goctl rpc protoc <name>.proto --go_out=. --go-grpc_out=. --zrpc_out=. -style gozero
3. 修改 Config，添加 DataSource / Cache
4. 修改 ServiceContext，注入 Model
5. 在 internal/logic/ 实现 RPC 方法
6. go build ./... 验证编译
```

## 3. 数据库 Model 变更

```
1. 在 deploy/sql/ 编写/更新 PostgreSQL DDL
2. 运行 goctl 生成 Model
   goctl model pg datasource -url="postgres://..." -table="<table>" -dir model -c -style gozero
3. 在 ServiceContext 中注入新 Model
4. 在 Logic 中使用 Model 方法
```

## 4. 修改已有 API

```
1. 编辑 desc/*.api 文件（修改 type 或路由）
2. 重新运行 goctl api go（不会覆盖 Logic 中的自定义代码）
3. 更新 Logic 层逻辑
4. go build ./... 验证编译
```

## 5. 代码提交

```
1. go build ./... 确保编译通过
2. golangci-lint run 检查风格
3. 使用 Angular 规范 + 中文编写 commit message
   feat: 添加用户注册接口
   fix: 修复登录密码校验逻辑
```
