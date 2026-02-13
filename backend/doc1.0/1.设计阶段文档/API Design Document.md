
---

# API 设计文档 (API Design Document)

**项目名称**：基于 Gin & Go-Zero 的微服务电商平台
**版本**：1.0 (MVP 扩展版)
**基础 URL**：`http://api.mall-domain.com/api/v1`

---

## 1. 概述

所有接口基于 **HTTP/1.1**，遵循 **RESTful** 风格，数据格式为 **JSON**。

### 1.1 认证机制
- **JWT (Bearer Token)**：`Authorization: Bearer <token>`
- **公开接口**：注册、登录、商品浏览、搜索
- **受保护接口**：下单、购物车、个人信息等

### 1.2 通用响应格式
```json
{ "code": 200, "msg": "success", "data": { ... } }
```

### 1.3 错误码定义

| 码 | 名称 | 说明 |
|---|---|---|
| 200 | OK | 成功 |
| 400 | Bad Request | 参数错误 |
| 401 | Unauthorized | Token 无效 |
| 403 | Forbidden | 权限不足 |
| 404 | Not Found | 资源不存在 |
| 500 | Internal Error | 服务器错误 |
| 1001 | Business Error | 通用业务错误 |
| 2001 | Out of Stock | 库存不足 |
| 2002 | System Busy | 排队中 |

---

## 2. 认证服务 (Auth Service) — 新增

### 2.1 登录
- **POST** `/auth/login` | Auth: No
```json
// 请求
{ "username": "tom123", "password": "password123" }
// 响应
{ "code": 200, "data": { "token": "eyJ...", "refresh_token": "xxx", "expire_in": 86400 } }
```

### 2.2 刷新 Token
- **POST** `/auth/refresh` | Auth: No
```json
// 请求
{ "refresh_token": "xxx" }
// 响应
{ "code": 200, "data": { "token": "eyJ...(新)", "expire_in": 86400 } }
```

### 2.3 注销
- **POST** `/auth/logout` | Auth: **Yes**
```json
// 响应
{ "code": 200, "msg": "已注销" }
```

---

## 3. 用户服务 (User Service)

### 3.1 用户注册
- **POST** `/user/register` | Auth: No
```json
// 请求
{ "username": "tom123", "password": "password123", "mobile": "13800138000" }
// 响应
{ "code": 200, "data": { "user_id": 1001 } }
```

### 3.2 获取用户信息
- **GET** `/user/info` | Auth: **Yes**
```json
// 响应
{ "code": 200, "data": { "id": 1001, "username": "tom123", "mobile": "138****8000", "member_level": 0 } }
```

### 3.3 用户地址列表 — 新增
- **GET** `/user/address/list` | Auth: **Yes**
```json
// 响应
{ "code": 200, "data": { "list": [{ "id": 1, "receiver": "张三", "phone": "138...", "province": "广东", "city": "深圳", "detail": "..." }] } }
```

### 3.4 新增地址 — 新增
- **POST** `/user/address/add` | Auth: **Yes**
```json
// 请求
{ "receiver": "张三", "phone": "13800138000", "province": "广东", "city": "深圳", "district": "南山", "detail": "科技园", "is_default": true }
```

---

## 4. 商品服务 (Product Service)

### 4.1 商品搜索 → 迁移至搜索服务
> **注意**：商品搜索功能已迁移至 **搜索服务 (Search Service)**，见第 5 节。

### 4.2 商品详情
- **GET** `/product/:id` | Auth: No
```json
// 响应
{ "code": 200, "data": { "id": 1, "name": "iPhone 15", "skus": [{ "sku_id": 101, "specs": "黑色+128G", "price": 5999.00 }] } }
```

### 4.3 商品列表 (按分类)
- **GET** `/product/list?category_id=1&page=1&size=20` | Auth: No

---

## 5. 搜索服务 (Search Service) — 新增

### 5.1 商品搜索
- **GET** `/search/product?keyword=iPhone&page=1&size=20` | Auth: No
```json
// 响应
{ "code": 200, "data": { "list": [{ "id": 1, "name": "iPhone 15", "price": 5999.00, "image": "..." }], "total": 100 } }
```

### 5.2 搜索联想
- **GET** `/search/suggest?keyword=iPh` | Auth: No
```json
{ "code": 200, "data": { "suggestions": ["iPhone 15", "iPhone 15 Pro", "iPhone 手机壳"] } }
```

### 5.3 热搜词
- **GET** `/search/hot` | Auth: No
```json
{ "code": 200, "data": { "keywords": ["iPhone", "连衣裙", "switch"] } }
```

---

## 6. 购物车服务 (Cart Service) — 新增

### 6.1 加入购物车
- **POST** `/cart/add` | Auth: **Yes**
```json
// 请求
{ "sku_id": 101, "quantity": 1 }
// 响应
{ "code": 200, "msg": "已加入购物车" }
```

### 6.2 购物车列表
- **GET** `/cart/list` | Auth: **Yes**
```json
{ "code": 200, "data": { "items": [{ "sku_id": 101, "product_name": "iPhone 15", "specs": "黑色+128G", "price": 5999.00, "quantity": 1, "selected": true }], "total_amount": 5999.00 } }
```

### 6.3 修改数量
- **PUT** `/cart/update` | Auth: **Yes**
```json
{ "sku_id": 101, "quantity": 2 }
```

### 6.4 删除商品
- **DELETE** `/cart/delete` | Auth: **Yes**
```json
{ "sku_ids": [101, 102] }
```

---

## 7. 营销服务 (Promotion Service) — 新增

### 7.1 可领优惠券列表
- **GET** `/coupon/list` | Auth: No

### 7.2 领取优惠券
- **POST** `/coupon/claim` | Auth: **Yes**
```json
{ "coupon_id": 1 }
```

### 7.3 我的优惠券
- **GET** `/coupon/mine?status=0` | Auth: **Yes**
```json
{ "code": 200, "data": { "list": [{ "id": 1, "name": "满100减20", "threshold": 100, "discount": 20, "end_time": "2026-03-01" }] } }
```

### 7.4 计算优惠价格 (内部)
- **RPC**: `PromotionRpc.CalculateDiscount(sku_ids, coupon_id)`

---

## 8. 订单服务 (Order Service)

### 8.1 创建订单 — 扩展
- **POST** `/order/create` | Auth: **Yes**
```json
// 请求
{ "address_id": 1, "coupon_id": 1, "items": [{ "sku_id": 101, "quantity": 1 }], "remark": "请尽快发货" }
// 响应
{ "code": 200, "data": { "order_id": "SN202602110001", "pay_amount": 5979.00, "expire_time": "2026-02-11T22:16:00" } }
```

### 8.2 订单详情
- **GET** `/order/:id` | Auth: **Yes**

### 8.3 订单列表
- **GET** `/order/list?status=0&page=1` | Auth: **Yes**

### 8.4 取消订单
- **POST** `/order/cancel` | Auth: **Yes**
```json
{ "order_id": "SN202602110001", "reason": "不想要了" }
```

---

## 9. 秒杀服务 (Seckill)

### 9.1 执行秒杀
- **POST** `/seckill/action` | Auth: **Yes**
```json
{ "sku_id": 101 }
// 响应 (排队中)
{ "code": 2002, "msg": "排队中" }
// 响应 (库存不足)
{ "code": 2001, "msg": "秒杀已结束" }
```

### 9.2 查询结果
- **GET** `/seckill/result?sku_id=101` | Auth: **Yes**
```json
{ "code": 200, "data": { "status": 1, "order_id": "SN..." } }
```

---

## 10. 支付服务 (Payment Service)

### 10.1 发起支付
- **POST** `/payment/pay` | Auth: **Yes**
```json
{ "order_id": "SN202602110001", "channel": "alipay" }
// 响应
{ "code": 200, "data": { "pay_url": "https://openapi.alipaydev.com/..." } }
```

### 10.2 支付回调 (第三方调用)
- **POST** `/payment/callback`

### 10.3 退款 — 新增
- **POST** `/payment/refund` | Auth: **Yes**
```json
{ "order_id": "SN...", "amount": 5979.00, "reason": "七天无理由" }
```

---

## 11. 物流服务 (Logistics Service) — 新增

### 11.1 计算运费
- **POST** `/logistics/freight` | Auth: No
```json
// 请求
{ "address_id": 1, "items": [{ "sku_id": 101, "quantity": 1 }] }
// 响应
{ "code": 200, "data": { "freight": 10.00, "free_threshold": 99.00 } }
```

### 11.2 物流轨迹
- **GET** `/logistics/track/:order_id` | Auth: **Yes**
```json
{ "code": 200, "data": { "tracking_no": "SF1234567890", "carrier": "顺丰", "traces": [{ "time": "2026-02-11 14:00", "desc": "已揽收" }] } }
```
