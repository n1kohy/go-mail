-- =============================================
-- 库存服务数据库 (db_inventory) - PostgreSQL
-- =============================================

CREATE TABLE IF NOT EXISTS stock (
    id          BIGSERIAL    PRIMARY KEY,
    sku_id      BIGINT       NOT NULL UNIQUE,
    total       INTEGER      NOT NULL DEFAULT 0,
    available   INTEGER      NOT NULL DEFAULT 0,
    locked      INTEGER      NOT NULL DEFAULT 0,
    version     INTEGER      NOT NULL DEFAULT 0,
    update_time TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE stock IS '库存表';
COMMENT ON COLUMN stock.available IS '可用库存 = total - locked';
COMMENT ON COLUMN stock.locked IS '锁定库存 (下单未支付)';
COMMENT ON COLUMN stock.version IS '乐观锁版本号';
