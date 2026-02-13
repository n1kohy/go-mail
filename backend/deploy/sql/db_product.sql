-- =============================================
-- 商品服务数据库 (db_product) - PostgreSQL
-- =============================================

-- 商品分类表
CREATE TABLE IF NOT EXISTS category (
    id        BIGSERIAL    PRIMARY KEY,
    parent_id BIGINT       NOT NULL DEFAULT 0,
    name      VARCHAR(50)  NOT NULL,
    level     SMALLINT     NOT NULL DEFAULT 1,
    sort      INTEGER      NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE category IS '商品分类表 (支持三级)';

-- 商品表 (SPU)
CREATE TABLE IF NOT EXISTS product (
    id           BIGSERIAL     PRIMARY KEY,
    category_id  BIGINT        NOT NULL,
    name         VARCHAR(100)  NOT NULL,
    sub_title    VARCHAR(255),
    main_image   VARCHAR(255),
    detail_html  TEXT,
    status       SMALLINT      NOT NULL DEFAULT 1,
    create_time  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    update_time  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_category ON product(category_id);
COMMENT ON TABLE product IS '商品 SPU 表';

-- 商品 SKU 表
CREATE TABLE IF NOT EXISTS product_sku (
    id         BIGSERIAL      PRIMARY KEY,
    product_id BIGINT         NOT NULL,
    specs      JSONB          NOT NULL DEFAULT '{}',
    price      NUMERIC(10,2)  NOT NULL DEFAULT 0.00,
    image      VARCHAR(255),
    create_time TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sku_product ON product_sku(product_id);
COMMENT ON TABLE product_sku IS '商品 SKU 表 (库存量单位)';
