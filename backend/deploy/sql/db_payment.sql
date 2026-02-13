-- =============================================
-- 支付服务数据库 (db_payment) - PostgreSQL
-- =============================================

CREATE TABLE IF NOT EXISTS payment_flow (
    id            BIGSERIAL      PRIMARY KEY,
    order_id      VARCHAR(64)    NOT NULL,
    trade_no      VARCHAR(64),
    amount        NUMERIC(10,2)  NOT NULL,
    channel       SMALLINT       NOT NULL DEFAULT 1,
    status        SMALLINT       NOT NULL DEFAULT 0,
    callback_time TIMESTAMPTZ,
    refund_amount NUMERIC(10,2)  DEFAULT 0.00,
    create_time   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    update_time   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_order ON payment_flow(order_id);
COMMENT ON TABLE payment_flow IS '支付流水表';
COMMENT ON COLUMN payment_flow.channel IS '1-支付宝, 2-微信';
COMMENT ON COLUMN payment_flow.status IS '0-未支付, 1-成功, 2-失败, 3-已退款';
