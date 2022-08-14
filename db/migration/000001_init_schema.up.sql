CREATE TYPE "location_status" AS ENUM (
  'original',
  'entrainment',
  'elsewhere'
);

CREATE TYPE "time_bucket" AS ENUM (
  'morning',
  'noon',
  'afternoon',
  'evening',
  'night'
);

CREATE TYPE "permission" AS ENUM (
  'level1',
  'level2',
  'level3'
);

CREATE TYPE "campus" AS ENUM (
  'qing',
  'sha'
);

CREATE TABLE "found" (
  "id" serial PRIMARY KEY,
  "create_at" timestamptz NOT NULL DEFAULT 'now()',
  "picker_openid" varchar NOT NULL,
  "found_date" date NOT NULL,
  "time_bucket" time_bucket NOT NULL,
  "location_id" smallserial NOT NULL,
  "location_info" varchar NOT NULL,
  "location_status" location_status NOT NULL,
  "type_id" smallserial NOT NULL,
  "item_info" varchar NOT NULL,
  "image" bytea NOT NULL,
  "image_key" varchar NOT NULL,
  "owner_info" varchar NOT NULL DEFAULT '',
  "addtional_info" varchar NOT NULL DEFAULT ''
);

CREATE TABLE "lost" (
  "id" serial PRIMARY KEY,
  "create_at" timestamptz NOT NULL DEFAULT 'now()',
  "owner_openid" varchar NOT NULL,
  "lost_date" date NOT NULL,
  "time_bucket" time_bucket NOT NULL,
  "type_id" smallserial NOT NULL,
  "item_info" varchar NOT NULL,
  "image" bytea NOT NULL DEFAULT '',
  "image_key" varchar NOT NULL DEFAULT '',
  "location_id" smallserial NOT NULL,
  "location_id1" smallserial NOT NULL DEFAULT 0,
  "location_id2" smallserial NOT NULL DEFAULT 0  
);

CREATE TABLE "match" (
  "id" serial PRIMARY KEY,
  "create_at" timestamptz NOT NULL DEFAULT 'now()',
  "picker_openid" varchar NOT NULL DEFAULT '',
  "owner_openid" varchar NOT NULL DEFAULT '',
  "found_date" date NOT NULL,
  "lost_date" date NOT NULL,
  "type_id" smallserial NOT NULL,
  "item_info" varchar NOT NULL,
  "image" bytea NOT NULL,
  "image_key" varchar NOT NULL,
  "comment" varchar NOT NULL DEFAULT ''
);

CREATE TABLE "location_wide" (
  "id" smallserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "campus" campus NOT NULL
);

CREATE TABLE "location_narrow" (
  "id" smallserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "wide_id" smallserial NOT NULL
);

CREATE TABLE "type_wide" (
  "id" smallserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "type_narrow" (
  "id" smallserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "wide_id" smallserial NOT NULL
);

CREATE TABLE "manager" (
  "id" smallserial PRIMARY KEY,
  "usr_openid" varchar UNIQUE NOT NULL,
  "permission" permission NOT NULL
);

CREATE TABLE "usr" (
  "openid" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "student_id" varchar NOT NULL,
  "avatar_url" varchar NOT NULL,
  "avatar" bytea NOT NULL
);

ALTER TABLE "found" ADD FOREIGN KEY ("picker_openid") REFERENCES "usr" ("openid");

ALTER TABLE "found" ADD FOREIGN KEY ("location_id") REFERENCES "location_narrow" ("id");

ALTER TABLE "found" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("owner_openid") REFERENCES "usr" ("openid");

ALTER TABLE "lost" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id") REFERENCES "location_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id1") REFERENCES "location_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id2") REFERENCES "location_narrow" ("id");

ALTER TABLE "match" ADD FOREIGN KEY ("picker_openid") REFERENCES "usr" ("openid");

ALTER TABLE "match" ADD FOREIGN KEY ("owner_openid") REFERENCES "usr" ("openid");

ALTER TABLE "match" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "location_narrow" ADD FOREIGN KEY ("wide_id") REFERENCES "location_wide" ("id");

ALTER TABLE "type_narrow" ADD FOREIGN KEY ("wide_id") REFERENCES "type_wide" ("id");

ALTER TABLE "manager" ADD FOREIGN KEY ("usr_openid") REFERENCES "usr" ("openid");
