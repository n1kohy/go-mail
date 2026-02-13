# 系统架构与设计文档 (System Design & Architecture)

**项目名称**：基于 Gin & Go-Zero 的微服务电商平台
**版本**：1.0 (MVP 扩展版)
**日期**：2026-02-11

---

## 1. 概述 (Overview)

本系统采用 **云原生微服务架构**，后端基于 **Go (Go-Zero)** 构建 **10 个核心微服务**，前端采用 **React** 构建 SPA 单页应用，通过 **Gin** 网关统一接入。

### 1.1 设计目标
1.  **高并发**：Redis + Lua 原子操作抗住秒杀流量
2.  **高可用**：熔断、限流、多副本部署
3.  **可扩展**：服务按域拆分，独立扩缩容
4.  **数据一致性**：本地消息表 + Kafka 最终一致性

---

## 2. 技术选型 (Technology Stack)

| 层次 | 技术组件 | 选型理由 |
| :--- | :--- | :--- |
| **前端** | React | 组件化 SPA |
| **API 网关** | Gin | HTTP 协议转换、鉴权、限流 |
| **微服务框架** | Go-Zero | 内置服务发现、熔断、zRPC |
| **服务注册** | Etcd | 强一致性 KV，服务发现 |
| **数据库** | PostgreSQL | 核心业务数据，Database-Per-Service，支持 JSONB/全文搜索 |
| **缓存** | Redis | 热点缓存、分布式锁、库存预热 |
| **搜索引擎** | Elasticsearch | 商品全文搜索 |
| **消息队列** | Kafka | 流量削峰、异步解耦 |
| **容器化** | Docker & K8s | 部署编排 |
| **监控** | Prometheus + Grafana | 指标监控与告警 |
| **日志** | ELK / Loki | 分布式日志收集 |

---

## 3. 微服务架构 (10 个核心服务)

```
                        [React 前端]
                             │
                             ▼
                    [API Gateway (Gin)]
                       JWT 认证 / 限流
                             │
        ┌────────┬──────┬────┴────┬──────┬────────┐
        ▼        ▼      ▼        ▼      ▼        ▼
    ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
    │ User │ │ Auth │ │Produc│ │Search│ │ Cart │ │Promo │
    │  API │ │  API │ │t API │ │  API │ │  API │ │  API │
    └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘
       ▼        ▼        ▼        ▼        ▼        ▼
    ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
    │ User │ │ Auth │ │Produc│ │  ES  │ │ Cart │ │Promo │
    │  RPC │ │  RPC │ │t RPC │ │Client│ │  RPC │ │  RPC │
    └──┬───┘ └──┬───┘ └──┬───┘ └──────┘ └──┬───┘ └──┬───┘
       ▼        ▼        ▼                  ▼        ▼
    db_user  db_auth  db_product         Redis+   db_promo
                      + Inventory RPC    db_cart
                           │
        ┌──────────────────┘
        ▼
    ┌──────┐    ┌──────┐    ┌──────┐    ┌──────┐
    │Invent│    │Order │    │Paymen│    │Logist│
    │  RPC │    │  API │    │t API │    │  API │
    └──┬───┘    └──┬───┘    └──┬───┘    └──┬───┘
       ▼           ▼           ▼           ▼
    db_invent   db_order    db_payment  db_logistics
```

### 3.1 用户域

#### 用户服务 (User Service)
- **职责**：用户注册、个人信息、地址管理、会员等级
- **HTTP**：`/user/register`, `/user/info`, `/user/address/*`
- **RPC**：`UserRpc.GetUser()`, `UserRpc.GetUserByMobile()`
- **数据库**：`db_user` → `user`, `user_address`

#### 认证服务 (Auth Service)
- **职责**：JWT 签发/刷新/注销，登录态管理
- **HTTP**：`/auth/login`, `/auth/refresh`, `/auth/logout`
- **RPC**：`AuthRpc.VerifyToken()` — 供所有服务验证 Token
- **数据库**：`db_auth` → `login_log`

### 3.2 商品域

#### 商品服务 (Product Service)
- **职责**：SPU/SKU 管理、分类、品牌
- **HTTP**：`/product/:id`, `/product/list`
- **RPC**：`ProductRpc.GetProduct()`, `ProductRpc.GetSku()`
- **数据库**：`db_product` → `category`, `product`, `product_sku`

#### 库存服务 (Inventory Service)
- **职责**：库存查询、扣减、预扣、释放，秒杀库存
- **RPC**：`InventoryRpc.DeductStock()`, `InventoryRpc.ReleaseStock()`, `InventoryRpc.PreDeductStock()`
- **数据库**：`db_inventory` → `stock`
- **缓存**：Redis Lua 原子扣减

#### 搜索服务 (Search Service)
- **职责**：商品全文搜索、搜索联想、热搜
- **HTTP**：`/search/product`, `/search/suggest`, `/search/hot`
- **技术**：Elasticsearch

### 3.3 交易域

#### 购物车服务 (Cart Service)
- **职责**：购物车增删改查、结算
- **HTTP**：`/cart/add`, `/cart/list`, `/cart/update`, `/cart/delete`, `/cart/checkout`
- **RPC**：`CartRpc.GetCartItems()` — 供订单服务获取
- **数据库**：`db_cart` → `cart_item`
- **缓存**：Redis (热数据) + PostgreSQL (持久化)

#### 营销服务 (Promotion Service)
- **职责**：优惠券管理、满减规则、活动配置
- **HTTP**：`/coupon/list`, `/coupon/claim`, `/coupon/mine`, `/promotion/calculate`
- **RPC**：`PromotionRpc.CalculateDiscount()` — 供订单服务计算优惠
- **数据库**：`db_promotion` → `coupon`, `coupon_record`, `promotion_activity`

#### 订单服务 (Order Service)
- **职责**：下单、订单状态机、超时取消
- **HTTP**：`/order/create`, `/order/:id`, `/order/list`
- **RPC 调用**：Cart → Inventory → Promotion → Logistics
- **状态机**：`PENDING → PAID → SHIPPED → COMPLETED / CANCELED`
- **数据库**：`db_order` → `order_master`, `order_item`, `local_message`

#### 支付服务 (Payment Service)
- **职责**：支付网关、回调处理、退款
- **HTTP**：`/payment/pay`, `/payment/callback`, `/payment/refund`
- **数据库**：`db_payment` → `payment_flow`

### 3.4 履约域

#### 物流服务 (Logistics Service)
- **职责**：运费计算、物流轨迹查询
- **HTTP**：`/logistics/freight`, `/logistics/track/:order_id`
- **RPC**：`LogisticsRpc.CalculateFreight()` — 供订单服务计算运费
- **数据库**：`db_logistics` → `freight_template`, `shipping_order`

---

## 4. 关键业务流程

### 4.1 下单流程（完整链路）
```
用户点击"提交订单"
  → [Cart RPC] 获取购物车选中商品
  → [Product RPC] 校验商品状态和价格
  → [Promotion RPC] 计算优惠（优惠券、满减）
  → [Logistics RPC] 计算运费
  → [Inventory RPC] 预扣库存（锁定）
  → [Order] 创建订单 + 写入本地消息表（事务）
  → [Kafka] 发送"超时取消"延迟消息（30分钟）
  → 返回 order_id + 待支付金额
```

### 4.2 秒杀流程
```
用户点击"立即秒杀"
  → [Gateway] JWT 校验 + 限流
  → [Redis] Lua 脚本原子扣减 seckill:stock:{sku_id}
  → [Redis] 设置防重锁 seckill:user:{uid}:{sku_id}
  → [Kafka] 发送异步下单消息
  → 返回 "排队中"
  → [Order 消费者] 消费 Kafka → 创建订单 → 扣 PostgreSQL 库存
  → 用户轮询 /seckill/result 查询结果
```

### 4.3 支付回调
```
用户完成支付
  → [支付宝/微信] 回调 /payment/callback
  → [Payment] 验签 → 更新 payment_flow
  → [Order RPC] 更新订单状态为 PAID
  → [Kafka] 发送"发货通知"消息
```

---

## 5. 数据库设计概要

遵循 **Database-Per-Service** 模式，共 8 个独立数据库：

| 数据库 | 服务 | 核心表 |
| :--- | :--- | :--- |
| `db_user` | 用户服务 | `user`, `user_address` |
| `db_auth` | 认证服务 | `login_log` |
| `db_product` | 商品服务 | `category`, `product`, `product_sku` |
| `db_inventory` | 库存服务 | `stock` |
| `db_cart` | 购物车服务 | `cart_item` |
| `db_promotion` | 营销服务 | `coupon`, `coupon_record`, `promotion_activity` |
| `db_order` | 订单服务 | `order_master`, `order_item`, `local_message` |
| `db_payment` | 支付服务 | `payment_flow` |

*(详细字段定义请参考 `Data Dictionary.md`)*

---

## 6. 接口规范

- **RESTful**：URL 资源化，标准 HTTP 方法
- **统一响应**：
  ```json
  { "code": 200, "msg": "success", "data": { ... } }
  ```
- **内部 RPC**：Protobuf 强类型接口
- **认证**：`Authorization: Bearer <jwt_token>`

---

## 7. 部署架构

- **Dev**：Docker Compose 一键拉起全套环境
- **Prod**：Kubernetes 集群，HPA 自动扩缩容
- **网络**：仅网关暴露公网，微服务置于内网 VPC
