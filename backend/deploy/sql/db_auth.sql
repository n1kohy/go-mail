-- =============================================
-- 认证服务数据库 (db_auth) - PostgreSQL
-- =============================================

CREATE TABLE IF NOT EXISTS login_log (
    id          BIGSERIAL     PRIMARY KEY,
    user_id     BIGINT        NOT NULL,
    login_type  SMALLINT      NOT NULL DEFAULT 1,
    ip          INET          NOT NULL,
    device      VARCHAR(200),
    status      SMALLINT      NOT NULL DEFAULT 1,
    create_time TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_login_log_user_id ON login_log(user_id);
COMMENT ON TABLE login_log IS '登录日志表';
COMMENT ON COLUMN login_log.login_type IS '1-密码, 2-验证码, 3-微信';
COMMENT ON COLUMN login_log.status IS '1-成功, 0-失败';
