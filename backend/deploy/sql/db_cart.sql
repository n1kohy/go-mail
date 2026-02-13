-- =============================================
-- 购物车服务数据库 (db_cart) - PostgreSQL
-- =============================================

CREATE TABLE IF NOT EXISTS cart_item (
    id           BIGSERIAL      PRIMARY KEY,
    user_id      BIGINT         NOT NULL,
    product_id   BIGINT         NOT NULL,
    sku_id       BIGINT         NOT NULL,
    product_name VARCHAR(100)   NOT NULL,
    sku_specs    VARCHAR(200),
    price        NUMERIC(10,2)  NOT NULL,
    quantity     INTEGER        NOT NULL DEFAULT 1,
    selected     BOOLEAN        NOT NULL DEFAULT TRUE,
    create_time  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    update_time  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, sku_id)
);

CREATE INDEX idx_cart_user ON cart_item(user_id);
COMMENT ON TABLE cart_item IS '购物车明细表';
