-- =============================================
-- 订单服务数据库 (db_order) - PostgreSQL
-- =============================================

-- 订单主表
CREATE TABLE IF NOT EXISTS order_master (
    id               VARCHAR(64)    PRIMARY KEY,
    user_id          BIGINT         NOT NULL,
    total_amount     NUMERIC(10,2)  NOT NULL DEFAULT 0.00,
    discount_amount  NUMERIC(10,2)  NOT NULL DEFAULT 0.00,
    freight_amount   NUMERIC(10,2)  NOT NULL DEFAULT 0.00,
    pay_amount       NUMERIC(10,2)  NOT NULL DEFAULT 0.00,
    coupon_id        BIGINT,
    pay_type         SMALLINT       DEFAULT 0,
    status           SMALLINT       NOT NULL DEFAULT 0,
    address_snapshot JSONB          NOT NULL DEFAULT '{}',
    expire_time      TIMESTAMPTZ,
    create_time      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    update_time      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_user ON order_master(user_id);
CREATE INDEX idx_order_status ON order_master(status);
COMMENT ON TABLE order_master IS '订单主表';
COMMENT ON COLUMN order_master.status IS '0-待支付, 1-已支付, 2-已发货, 3-已完成, 4-已取消';

-- 订单明细表
CREATE TABLE IF NOT EXISTS order_item (
    id           BIGSERIAL      PRIMARY KEY,
    order_id     VARCHAR(64)    NOT NULL,
    product_id   BIGINT         NOT NULL,
    product_name VARCHAR(100)   NOT NULL,
    sku_id       BIGINT         NOT NULL,
    sku_specs    VARCHAR(200),
    price        NUMERIC(10,2)  NOT NULL,
    quantity     INTEGER        NOT NULL DEFAULT 1,
    create_time  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_item_order ON order_item(order_id);
COMMENT ON TABLE order_item IS '订单明细表';

-- 本地消息表 (分布式事务)
CREATE TABLE IF NOT EXISTS local_message (
    id          BIGSERIAL    PRIMARY KEY,
    tx_id       VARCHAR(64)  NOT NULL,
    topic       VARCHAR(50)  NOT NULL,
    content     JSONB        NOT NULL DEFAULT '{}',
    state       SMALLINT     NOT NULL DEFAULT 0,
    retry_count INTEGER      NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_local_msg_state ON local_message(state);
COMMENT ON TABLE local_message IS '本地消息表 (分布式事务最终一致性)';
COMMENT ON COLUMN local_message.state IS '0-待发送, 1-成功, 2-失败';
