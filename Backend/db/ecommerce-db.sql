-- ecommerce-db.sql
-- สคริปต์สร้าง schema ของ database จาก 0 (bootstrap)
-- ใช้ตอนสร้าง DB ใหม่:  psql "<connection-string>" -f ecommerce-db.sql
--
-- gen_random_uuid() มีใน Postgres 13+ อยู่แล้ว (Supabase ใช้ได้เลย)
-- ถ้า Postgres เก่ากว่านั้น: CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ==========================================================
-- users (id = UUID)
-- ==========================================================
CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email          VARCHAR(255) NOT NULL,
    username       VARCHAR(255) UNIQUE,
    password       VARCHAR(255),                       -- เก็บ hash (NULL ได้สำหรับ OAuth)
    full_name      VARCHAR(150) NOT NULL,
    role           VARCHAR(20)  NOT NULL DEFAULT 'customer',
    status         VARCHAR(20)  NOT NULL DEFAULT 'active',
    active_flag    BOOLEAN      NOT NULL DEFAULT true,
    provider       VARCHAR(20)  NOT NULL DEFAULT 'local',  -- local | google | facebook
    provider_id    VARCHAR(255),
    email_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_users_email    ON users (lower(email));
CREATE UNIQUE INDEX uq_users_provider ON users (provider, provider_id) WHERE provider_id IS NOT NULL;

-- ==========================================================
-- Thailand geography: province -> district -> sub_district
-- id ใช้รหัสจากชุดข้อมูลมาตรฐานไทย (import มา)
-- ==========================================================
CREATE TABLE provinces (
    id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name_th VARCHAR(100) NOT NULL,
    name_en VARCHAR(100)
);

CREATE TABLE districts (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    province_id INT NOT NULL,
    name_th     VARCHAR(100) NOT NULL,
    name_en     VARCHAR(100),
    FOREIGN KEY (province_id) REFERENCES provinces(id)
);

CREATE TABLE sub_districts (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    district_id INT NOT NULL,
    name_th     VARCHAR(100) NOT NULL,
    name_en     VARCHAR(100),
    FOREIGN KEY (district_id) REFERENCES districts(id)
);

CREATE INDEX idx_districts_province     ON districts (province_id);
CREATE INDEX idx_sub_districts_district ON sub_districts (district_id);

-- ==========================================================
-- user addresses (1 user มีได้หลายที่อยู่) — id = running int
-- ==========================================================
CREATE TABLE user_addresses (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         UUID NOT NULL,
    recipient       VARCHAR(150) NOT NULL,       -- ชื่อผู้รับ
    phone           VARCHAR(20),
    line            VARCHAR(255) NOT NULL,       -- บ้านเลขที่/หมู่/ถนน/ซอย
    sub_district_id INT NOT NULL,
    district_id     INT NOT NULL,
    province_id     INT NOT NULL,
    postal_code     VARCHAR(5) NOT NULL,
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id)         REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sub_district_id) REFERENCES sub_districts(id),
    FOREIGN KEY (district_id)     REFERENCES districts(id),
    FOREIGN KEY (province_id)     REFERENCES provinces(id)
);

CREATE INDEX idx_user_addresses_user ON user_addresses (user_id);
