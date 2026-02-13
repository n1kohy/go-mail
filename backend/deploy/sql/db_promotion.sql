-- =============================================
-- 营销服务数据库 (db_promotion) - PostgreSQL
-- =============================================

-- 优惠券模板表
CREATE TABLE IF NOT EXISTS coupon (
    id           BIGSERIAL      PRIMARY KEY,
    name         VARCHAR(100)   NOT NULL,
    type         SMALLINT       NOT NULL DEFAULT 1,
    threshold    NUMERIC(10,2)  NOT NULL DEFAULT 0,
    discount     NUMERIC(10,2)  NOT NULL DEFAULT 0,
    total_count  INTEGER        NOT NULL DEFAULT 0,
    remain_count INTEGER        NOT NULL DEFAULT 0,
    start_time   TIMESTAMPTZ    NOT NULL,
    end_time     TIMESTAMPTZ    NOT NULL,
    status       SMALLINT       NOT NULL DEFAULT 1,
    create_time  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE coupon IS '优惠券模板表';
COMMENT ON COLUMN coupon.type IS '1-满减, 2-折扣, 3-兑换券';

-- 用户领券记录表
CREATE TABLE IF NOT EXISTS coupon_record (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL,
    coupon_id     BIGINT       NOT NULL,
    status        SMALLINT     NOT NULL DEFAULT 0,
    used_order_id VARCHAR(64),
    create_time   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coupon_record_user ON coupon_record(user_id);
COMMENT ON TABLE coupon_record IS '用户领券记录';
COMMENT ON COLUMN coupon_record.status IS '0-未使用, 1-已使用, 2-已过期';

-- 活动配置表
CREATE TABLE IF NOT EXISTS promotion_activity (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    type       SMALLINT     NOT NULL,
    rules      JSONB        NOT NULL DEFAULT '{}',
    start_time TIMESTAMPTZ  NOT NULL,
    end_time   TIMESTAMPTZ  NOT NULL,
    status     SMALLINT     NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE promotion_activity IS '营销活动配置表';
COMMENT ON COLUMN promotion_activity.type IS '1-秒杀, 2-满减, 3-拼团';
