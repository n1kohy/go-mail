-- =============================================
-- 用户服务数据库 (db_user) - PostgreSQL
-- =============================================

-- 用户表
CREATE TABLE IF NOT EXISTS "user" (
    id            BIGSERIAL     PRIMARY KEY,
    username      VARCHAR(50)   NOT NULL UNIQUE,
    password      VARCHAR(100)  NOT NULL,
    mobile        VARCHAR(20)   NOT NULL UNIQUE,
    avatar        VARCHAR(255)  NOT NULL DEFAULT '',
    gender        SMALLINT      NOT NULL DEFAULT 0,
    role          SMALLINT      NOT NULL DEFAULT 0,
    member_level  SMALLINT      NOT NULL DEFAULT 0,
    create_time   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    update_time   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "user" IS '用户基础信息表';
COMMENT ON COLUMN "user".gender IS '0-保密, 1-男, 2-女';
COMMENT ON COLUMN "user".role IS '0-普通用户, 1-管理员';
COMMENT ON COLUMN "user".member_level IS '0-普通, 1-银卡, 2-金卡, 3-钻石';

-- 用户地址表
CREATE TABLE IF NOT EXISTS user_address (
    id          BIGSERIAL     PRIMARY KEY,
    user_id     BIGINT        NOT NULL,
    receiver    VARCHAR(50)   NOT NULL,
    phone       VARCHAR(20)   NOT NULL,
    province    VARCHAR(20)   NOT NULL,
    city        VARCHAR(20)   NOT NULL,
    district    VARCHAR(20)   NOT NULL,
    detail      VARCHAR(200)  NOT NULL,
    is_default  BOOLEAN       NOT NULL DEFAULT FALSE,
    create_time TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_address_user_id ON user_address(user_id);
COMMENT ON TABLE user_address IS '用户收货地址表';
