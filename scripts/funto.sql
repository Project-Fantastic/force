-- -------------------------------------------------------------
-- TablePlus 2.11.2(278)
--
-- https://tableplus.com/
--
-- Database: funto
-- Generation Time: 2019-11-21 21:21:28.7040
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."billing_requests";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS billing_requests_id_seq;

-- Table Definition
CREATE TABLE "public"."billing_requests" (
    "id" int4 NOT NULL DEFAULT nextval('billing_requests_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_product_id" int4 NOT NULL,
    "user_product_member_id" int4 NOT NULL,
    "status" int4 NOT NULL DEFAULT 0,
    "amount" numeric NOT NULL,
    "start_date" timestamptz,
    "end_date" timestamptz,
    "expiration_date" timestamptz,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."notifications";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS notifications_id_seq;

-- Table Definition
CREATE TABLE "public"."notifications" (
    "id" int4 NOT NULL DEFAULT nextval('notifications_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" int4 NOT NULL,
    "notification_type" int4 NOT NULL,
    "object_id" int4 NOT NULL,
    "status" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."products";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS products_id_seq;

-- Table Definition
CREATE TABLE "public"."products" (
    "id" int4 NOT NULL DEFAULT nextval('products_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name" varchar(512) NOT NULL,
    "billing_type" int4 NOT NULL DEFAULT 0,
    "is_fixed_price" bool NOT NULL DEFAULT true,
    "max_member_count" int8 NOT NULL DEFAULT '1'::bigint,
    "account_support_types" int4 DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."user_product_accounts";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_product_accounts_id_seq;

-- Table Definition
CREATE TABLE "public"."user_product_accounts" (
    "id" int4 NOT NULL DEFAULT nextval('user_product_accounts_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" int4 NOT NULL,
    "product_id" int4 NOT NULL,
    "account" varchar(512) NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."user_product_members";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_product_members_id_seq;

-- Table Definition
CREATE TABLE "public"."user_product_members" (
    "id" int4 NOT NULL DEFAULT nextval('user_product_members_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" int4 NOT NULL,
    "user_product_id" int4 NOT NULL,
    "is_host" bool NOT NULL DEFAULT false,
    "status" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."user_products";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_products_id_seq;

-- Table Definition
CREATE TABLE "public"."user_products" (
    "id" int4 NOT NULL DEFAULT nextval('user_products_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "host_id" int4 NOT NULL,
    "product_id" int4 NOT NULL,
    "title" varchar(512) NOT NULL,
    "description" text,
    "active" bool NOT NULL DEFAULT true,
    "min_member_count" int8 NOT NULL DEFAULT '1'::bigint,
    "max_member_count" int8 NOT NULL DEFAULT '1'::bigint,
    "min_price" numeric NOT NULL,
    "max_price" numeric NOT NULL,
    "total_price" numeric NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."user_third_party_payments";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_third_party_payments_id_seq;

-- Table Definition
CREATE TABLE "public"."user_third_party_payments" (
    "id" int4 NOT NULL DEFAULT nextval('user_third_party_payments_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" int4 NOT NULL,
    "payment_type" int4 NOT NULL DEFAULT 0,
    "payment_id" varchar(512) NOT NULL,
    "payment_id_type" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."users";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "email" varchar(255) NOT NULL,
    "phone_number" varchar(20),
    "password" varchar(512) NOT NULL,
    "first_name" varchar(100) NOT NULL,
    "last_name" varchar(100) NOT NULL,
    "status" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."billing_requests" ADD FOREIGN KEY ("user_product_member_id") REFERENCES "public"."user_product_members"("id") ON DELETE CASCADE;
ALTER TABLE "public"."billing_requests" ADD FOREIGN KEY ("user_product_id") REFERENCES "public"."user_products"("id") ON DELETE CASCADE;
ALTER TABLE "public"."notifications" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_product_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_product_accounts" ADD FOREIGN KEY ("product_id") REFERENCES "public"."products"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_product_members" ADD FOREIGN KEY ("user_product_id") REFERENCES "public"."user_products"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_product_members" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_products" ADD FOREIGN KEY ("host_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_products" ADD FOREIGN KEY ("product_id") REFERENCES "public"."products"("id") ON DELETE CASCADE;
ALTER TABLE "public"."user_third_party_payments" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
