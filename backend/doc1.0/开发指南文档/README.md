# Go-Zero 微服务电商平台 — 开发指南

**环境**：Go 1.25+ / goctl v1.9+ / PostgreSQL 16 / Redis
**版本**：1.0 (MVP 10 核心服务)

---

## 文档索引

| 阶段 | 文档 | 状态 | 说明 |
|:---|:---|:---|:---|
| 1 | [阶段1-分析与映射](01-分析与映射.md) | ✅ | 服务清单、目录结构、开发顺序 |
| 2 | [阶段2-项目初始化](02-项目初始化.md) | ✅ | Go Module、建表脚本、公共包、Spec-First 生成 |
| 3 | [阶段3-数据模型与业务逻辑](03-数据模型与业务逻辑.md) | ✅ | PostgreSQL Model、Config、Logic、Handler 统一响应 |
| 4 | [阶段4-商品与库存](04-商品与库存.md) | ✅ | Product API+RPC、Inventory RPC、乐观锁扣减 |
| 5 | [阶段5-搜索与购物车](05-搜索与购物车.md) | ✅ | Search API、Cart API+RPC、Upsert 购物车 |
| 6 | [阶段6-营销与订单](06-营销与订单.md) | ✅ | Promotion API+RPC、Order API+RPC |
| 7 | [阶段7-支付与秒杀](07-支付与秒杀.md) | ✅ | Payment API+RPC、Seckill API |
| 8 | [阶段8-物流与集成](08-物流与集成.md) | ✅ | Logistics API+RPC、Payment/Logistics→OrderRPC 集成 |
| 9 | [阶段9-集成与优化](09-集成与优化.md) | ✅ | Order→4RPC 完整下单链路、CancelOrder 库存回滚 |
| 10 | [阶段10-部署与运维](10-部署与运维.md) | 📋 | 下一阶段：Docker 容器化 + 监控 |
| — | [附录-速查参考](附录-速查参考.md) | ✅ | 工作流图、goctl 命令、知识库索引 |

---

## 核心工作流

```
定义 .api/.proto  →  goctl 生成  →  改 Config/SvcCtx  →  写 Logic  →  改 Handler  →  go build
```

## 服务清单 (全部 ✅)

| 服务 | 端口 | 类型 | 跨服务调用 |
|:---|:---|:---|:---|
| User API | 8801 | API+RPC | — |
| Auth API | 8802 | API | → User RPC |
| Product API | 8803 | API+RPC | — |
| Search API | 8804 | API | — |
| Cart API | 8805 | API+RPC | — |
| Promotion API | 8806 | API+RPC | — |
| **Order API** | **8807** | **API+RPC** | **→ Cart/Product/Inventory/Promotion RPC** |
| **Payment API** | **8808** | **API+RPC** | **→ Order RPC** |
| Seckill API | 8809 | API | — |
| **Logistics API** | **8810** | **API+RPC** | **→ Order RPC** |

**全部 10 个核心服务（18 个进程）开发完成，跨服务 RPC 调用链集成完毕 ✅**
