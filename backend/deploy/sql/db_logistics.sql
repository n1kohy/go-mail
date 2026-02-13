-- =============================================
-- 物流服务数据库 (db_logistics) - PostgreSQL
-- =============================================

-- 运费模板表
CREATE TABLE IF NOT EXISTS freight_template (
    id             BIGSERIAL      PRIMARY KEY,
    name           VARCHAR(50)    NOT NULL,
    free_threshold NUMERIC(10,2)  DEFAULT 0,
    base_fee       NUMERIC(10,2)  NOT NULL DEFAULT 0,
    rules          JSONB,
    create_time    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE freight_template IS '运费模板表';

-- 物流单表
CREATE TABLE IF NOT EXISTS shipping_order (
    id          BIGSERIAL    PRIMARY KEY,
    order_id    VARCHAR(64)  NOT NULL,
    tracking_no VARCHAR(50),
    carrier     VARCHAR(50),
    status      SMALLINT     NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipping_order ON shipping_order(order_id);
COMMENT ON TABLE shipping_order IS '物流单表';
COMMENT ON COLUMN shipping_order.status IS '0-待发货, 1-已发货, 2-已签收';
