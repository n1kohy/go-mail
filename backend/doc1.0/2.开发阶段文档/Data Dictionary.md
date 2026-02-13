
---

# 数据字典 (Data Dictionary)

**项目名称**：基于 Gin & Go-Zero 的微服务电商平台
**版本**：1.0 (MVP 扩展版)
**日期**：2026-02-11
**数据库类型**：PostgreSQL 16

**设计原则 (遵循 Zero-Skills 最佳实践)**：
1.  **微服务隔离**：Database-Per-Service，严禁跨库 Join。
2.  **命名规范**：所有表名、字段名均使用 `snake_case`。
3.  **通用字段**：所有业务表**必须**包含 `create_time` 和 `update_time`。
4.  **字符集**：`UTF-8`（PostgreSQL 默认）。

> **PostgreSQL 类型映射说明**：
> - MySQL `BIGINT AUTO_INCREMENT` → PostgreSQL `BIGSERIAL`
> - MySQL `TINYINT` → PostgreSQL `SMALLINT`
> - MySQL `JSON` → PostgreSQL `JSONB`（支持索引）
> - MySQL `TIMESTAMP DEFAULT CURRENT_TIMESTAMP` → PostgreSQL `TIMESTAMPTZ DEFAULT NOW()`

---

## 1. 用户服务数据库 (`db_user`)

### 1.1 用户表 (`user`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **username** | VARCHAR(50) | **Y** | | 用户名，**唯一索引** |
| **password** | VARCHAR(100) | **Y** | | Bcrypt 加密密码 |
| **mobile** | VARCHAR(20) | **Y** | | 手机号，**唯一索引** |
| **avatar** | VARCHAR(255) | N | '' | 头像 URL |
| **gender** | SMALLINT | **Y** | 0 | 0-保密, 1-男, 2-女 |
| **role** | SMALLINT | **Y** | 0 | 0-普通用户, 1-管理员 |
| **member_level** | SMALLINT | **Y** | 0 | 0-普通, 1-银卡, 2-金卡, 3-钻石 |
| **create_time** | TIMESTAMPTZ | N | NOW() | 创建时间 |
| **update_time** | TIMESTAMPTZ | N | NOW() | 更新时间 |

### 1.2 用户地址表 (`user_address`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **user_id** | BIGINT | **Y** | | 关联用户 ID |
| **receiver** | VARCHAR(50) | **Y** | | 收货人姓名 |
| **phone** | VARCHAR(20) | **Y** | | 联系电话 |
| **province** | VARCHAR(20) | **Y** | | 省 |
| **city** | VARCHAR(20) | **Y** | | 市 |
| **district** | VARCHAR(20) | **Y** | | 区/县 |
| **detail** | VARCHAR(200) | **Y** | | 详细地址 |
| **is_default** | BOOLEAN | **Y** | FALSE | 是否默认地址 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

---

## 2. 认证服务数据库 (`db_auth`)

### 2.1 登录日志表 (`login_log`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **user_id** | BIGINT | **Y** | | 用户 ID |
| **login_type** | SMALLINT | **Y** | 1 | 1-密码, 2-验证码, 3-微信 |
| **ip** | INET | **Y** | | 登录 IP (PostgreSQL 原生 IP 类型) |
| **device** | VARCHAR(200) | N | | 设备信息 |
| **status** | SMALLINT | **Y** | 1 | 1-成功, 0-失败 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

---

## 3. 商品服务数据库 (`db_product`)

### 3.1 商品分类表 (`category`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **parent_id** | BIGINT | **Y** | 0 | 父分类 ID，0=根 |
| **name** | VARCHAR(50) | **Y** | | 分类名称 |
| **level** | SMALLINT | **Y** | 1 | 1/2/3 级 |
| **sort** | INTEGER | N | 0 | 排序值 |

### 3.2 商品表 (`product`) — SPU

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **category_id** | BIGINT | **Y** | | 关联分类 |
| **name** | VARCHAR(100) | **Y** | | 商品名称 (支持 `tsvector` 全文搜索) |
| **sub_title** | VARCHAR(255) | N | | 副标题 |
| **main_image** | VARCHAR(255) | N | | 主图 URL |
| **detail_html** | TEXT | N | | 详情富文本 |
| **status** | SMALLINT | **Y** | 1 | 1-上架, 0-下架 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

### 3.3 商品 SKU 表 (`product_sku`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **product_id** | BIGINT | **Y** | | 关联商品 |
| **specs** | JSONB | **Y** | | 规格 `{"Color":"Red"}` |
| **price** | NUMERIC(10,2) | **Y** | 0.00 | 售价 |
| **image** | VARCHAR(255) | N | | SKU 图片 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

---

## 4. 库存服务数据库 (`db_inventory`)

### 4.1 库存表 (`stock`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **sku_id** | BIGINT | **Y** | | SKU ID，**唯一索引** |
| **total** | INTEGER | **Y** | 0 | 总库存 |
| **available** | INTEGER | **Y** | 0 | 可用库存 |
| **locked** | INTEGER | **Y** | 0 | 锁定库存 |
| **version** | INTEGER | **Y** | 0 | 乐观锁版本号 |
| **update_time** | TIMESTAMPTZ | N | NOW() | |

---

## 5. 购物车服务数据库 (`db_cart`)

### 5.1 购物车明细表 (`cart_item`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **user_id** | BIGINT | **Y** | | 用户 ID |
| **product_id** | BIGINT | **Y** | | SPU ID |
| **sku_id** | BIGINT | **Y** | | SKU ID |
| **product_name** | VARCHAR(100) | **Y** | | 商品名称快照 |
| **sku_specs** | VARCHAR(200) | N | | 规格描述 |
| **price** | NUMERIC(10,2) | **Y** | | 加入时价格 |
| **quantity** | INTEGER | **Y** | 1 | 数量 |
| **selected** | BOOLEAN | **Y** | TRUE | 是否选中 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

> **唯一索引**：`UNIQUE(user_id, sku_id)`

---

## 6. 营销服务数据库 (`db_promotion`)

### 6.1 优惠券模板表 (`coupon`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **name** | VARCHAR(100) | **Y** | | 券名称 |
| **type** | SMALLINT | **Y** | 1 | 1-满减, 2-折扣, 3-兑换券 |
| **threshold** | NUMERIC(10,2) | **Y** | 0 | 使用门槛 |
| **discount** | NUMERIC(10,2) | **Y** | 0 | 优惠金额/折扣率 |
| **total_count** | INTEGER | **Y** | 0 | 发行总量 |
| **remain_count** | INTEGER | **Y** | 0 | 剩余可领 |
| **start_time** | TIMESTAMPTZ | **Y** | | 生效时间 |
| **end_time** | TIMESTAMPTZ | **Y** | | 失效时间 |
| **status** | SMALLINT | **Y** | 1 | 1-进行中, 0-已结束 |

### 6.2 用户领券记录表 (`coupon_record`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **user_id** | BIGINT | **Y** | | 用户 ID |
| **coupon_id** | BIGINT | **Y** | | 优惠券 ID |
| **status** | SMALLINT | **Y** | 0 | 0-未使用, 1-已使用, 2-已过期 |
| **used_order_id** | VARCHAR(64) | N | | 使用的订单号 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

### 6.3 活动配置表 (`promotion_activity`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **name** | VARCHAR(100) | **Y** | | 活动名称 |
| **type** | SMALLINT | **Y** | | 1-秒杀, 2-满减, 3-拼团 |
| **rules** | JSONB | **Y** | | 活动规则配置 |
| **start_time** | TIMESTAMPTZ | **Y** | | 开始时间 |
| **end_time** | TIMESTAMPTZ | **Y** | | 结束时间 |
| **status** | SMALLINT | **Y** | 0 | 0-未开始, 1-进行中, 2-已结束 |

---

## 7. 订单服务数据库 (`db_order`)

### 7.1 订单主表 (`order_master`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | VARCHAR(64) | **Y** | | 订单号 (雪花算法) |
| **user_id** | BIGINT | **Y** | | 用户 ID |
| **total_amount** | NUMERIC(10,2) | **Y** | 0.00 | 商品总金额 |
| **discount_amount** | NUMERIC(10,2) | **Y** | 0.00 | 优惠金额 |
| **freight_amount** | NUMERIC(10,2) | **Y** | 0.00 | 运费 |
| **pay_amount** | NUMERIC(10,2) | **Y** | 0.00 | 实付金额 |
| **coupon_id** | BIGINT | N | | 使用的优惠券 ID |
| **pay_type** | SMALLINT | N | 0 | 1-支付宝, 2-微信 |
| **status** | SMALLINT | **Y** | 0 | 0-待支付, 1-已支付, 2-已发货, 3-已完成, 4-已取消 |
| **address_snapshot** | JSONB | **Y** | | 收货地址快照 |
| **expire_time** | TIMESTAMPTZ | N | | 支付超时 (30分钟) |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

### 7.2 订单明细表 (`order_item`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **order_id** | VARCHAR(64) | **Y** | | 订单号 |
| **product_id** | BIGINT | **Y** | | 商品 ID |
| **product_name** | VARCHAR(100) | **Y** | | 商品名快照 |
| **sku_id** | BIGINT | **Y** | | SKU ID |
| **sku_specs** | VARCHAR(200) | N | | 规格快照 |
| **price** | NUMERIC(10,2) | **Y** | | 成交单价 |
| **quantity** | INTEGER | **Y** | 1 | 数量 |

### 7.3 本地消息表 (`local_message`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **tx_id** | VARCHAR(64) | **Y** | | 业务 ID (订单号) |
| **topic** | VARCHAR(50) | **Y** | | Kafka Topic |
| **content** | JSONB | **Y** | | 消息内容 |
| **state** | SMALLINT | **Y** | 0 | 0-待发送, 1-成功, 2-失败 |
| **retry_count** | INTEGER | N | 0 | 重试次数 |

---

## 8. 支付服务数据库 (`db_payment`)

### 8.1 支付流水表 (`payment_flow`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **order_id** | VARCHAR(64) | **Y** | | 内部订单号 |
| **trade_no** | VARCHAR(64) | N | | 第三方交易号 |
| **amount** | NUMERIC(10,2) | **Y** | | 交易金额 |
| **channel** | SMALLINT | **Y** | 1 | 1-支付宝, 2-微信 |
| **status** | SMALLINT | **Y** | 0 | 0-未支付, 1-成功, 2-失败, 3-已退款 |
| **callback_time** | TIMESTAMPTZ | N | | 回调时间 |
| **refund_amount** | NUMERIC(10,2) | N | 0.00 | 退款金额 |

---

## 9. 物流服务数据库 (`db_logistics`)

### 9.1 运费模板表 (`freight_template`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **name** | VARCHAR(50) | **Y** | | 模板名称 |
| **free_threshold** | NUMERIC(10,2) | N | 0 | 满X元免运费 |
| **base_fee** | NUMERIC(10,2) | **Y** | 0 | 基础运费 |
| **rules** | JSONB | N | | 区域运费规则 |

### 9.2 物流单表 (`shipping_order`)

| 字段名 | 数据类型 | 必填 | 默认值 | 描述/约束 |
| --- | --- | --- | --- | --- |
| **id** | BIGSERIAL | **Y** | 自增 | 主键 |
| **order_id** | VARCHAR(64) | **Y** | | 关联订单号 |
| **tracking_no** | VARCHAR(50) | N | | 物流单号 |
| **carrier** | VARCHAR(50) | N | | 快递公司 |
| **status** | SMALLINT | **Y** | 0 | 0-待发货, 1-已发货, 2-已签收 |
| **create_time** | TIMESTAMPTZ | N | NOW() | |

---

## 10. Redis 数据结构

| Key 格式 | 类型 | 说明 |
| :--- | :--- | :--- |
| `seckill:stock:{sku_id}` | String | 秒杀库存 (Lua 原子扣减) |
| `seckill:user:{uid}:{sku_id}` | String | 防重锁 |
| `cache:product:{id}` | String/JSON | 商品详情缓存，TTL=1h |
| `cart:{user_id}` | Hash | 购物车热数据 |
| `auth:token:{user_id}` | String | Token 黑名单/踢人 |
| `coupon:remain:{coupon_id}` | String | 优惠券剩余数量 (原子扣减) |
