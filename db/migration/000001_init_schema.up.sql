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
  'level0',
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
  "create_at" timestampz NOT NULL DEFAULT 'now()',
  "picker_openid" varchar NOT NULL,
  "found_date" date NOT NULL,
  "bucket" time_bucket NOT NULL,
  "location_id" smallserial NOT NULL,
  "location_info" varcher NOT NULL,
  "location_status" location_status NOT NULL,
  "type_id" smallserial NOT NULL,
  "item_info" varchar NOT NULL,
  "image" varchar NOT NULL,
  "image_key" varchar NOT NULL,
  "owner_info" varchar,
  "addtional_info" varcher
);

CREATE TABLE "lost" (
  "id" serial PRIMARY KEY,
  "create_at" timestampz NOT NULL DEFAULT 'now()',
  "owner_openid" varchar NOT NULL,
  "found_date" date NOT NULL,
  "bucket" time_bucket NOT NULL,
  "location_id" smallserial NOT NULL,
  "location_id1" smallserial,
  "location_id2" smallserial,
  "type_id" smallserial NOT NULL
);

CREATE TABLE "match" (
  "id" serial PRIMARY KEY,
  "create_at" timestampz NOT NULL DEFAULT 'now()',
  "found_date" date NOT NULL,
  "lost_date" date NOT NULL,
  "picker_openid" varchar NOT NULL,
  "owner_openid" varchar NOT NULL,
  "type_id" smallserial NOT NULL,
  "item_info" varchar NOT NULL,
  "image" varchar NOT NULL,
  "image_key" varchar NOT NULL,
  "comment" varcher
);

CREATE TABLE "location_wide" (
  "id" smallserial PRIMARY KEY,
  "name" varcher NOT NULL,
  "campus" campus NOT NULL
);

CREATE TABLE "location_narrow" (
  "id" smallserial PRIMARY KEY,
  "name" varcher NOT NULL,
  "wide_id" smallserial NOT NULL
);

CREATE TABLE "type_wide" (
  "id" smallserial PRIMARY KEY,
  "name" varcher NOT NULL
);

CREATE TABLE "type_narrow" (
  "id" smallserial PRIMARY KEY,
  "name" varcher NOT NULL,
  "wide_id" smallserial NOT NULL
);

CREATE TABLE "manager" (
  "id" smallserial PRIMARY KEY,
  "user_openid" varcher NOT NULL,
  "permission" permission NOT NULL
);

CREATE TABLE "user" (
  "openid" varcher PRIMARY KEY,
  "name" varchar NOT NULL,
  "student_id" varcher NOT NULL,
  "avatar" varchar NOT NULL
);

ALTER TABLE "found" ADD FOREIGN KEY ("picker_openid") REFERENCES "user" ("openid");

ALTER TABLE "found" ADD FOREIGN KEY ("location_id") REFERENCES "location_narrow" ("id");

ALTER TABLE "found" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("owner_openid") REFERENCES "user" ("openid");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id") REFERENCES "location_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id1") REFERENCES "location_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("location_id2") REFERENCES "location_narrow" ("id");

ALTER TABLE "lost" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "match" ADD FOREIGN KEY ("picker_openid") REFERENCES "user" ("openid");

ALTER TABLE "match" ADD FOREIGN KEY ("owner_openid") REFERENCES "user" ("openid");

ALTER TABLE "match" ADD FOREIGN KEY ("type_id") REFERENCES "type_narrow" ("id");

ALTER TABLE "location_narrow" ADD FOREIGN KEY ("wide_id") REFERENCES "location_wide" ("id");

ALTER TABLE "type_narrow" ADD FOREIGN KEY ("wide_id") REFERENCES "type_wide" ("id");

ALTER TABLE "manager" ADD FOREIGN KEY ("user_openid") REFERENCES "user" ("openid");
