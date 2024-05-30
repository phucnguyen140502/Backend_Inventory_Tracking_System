-- +migrate Up
CREATE TABLE "warehouse" (
  "warehouse_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "location" varchar(50),
  "capacity" integer
);

CREATE TABLE "supplier" (
  "supplier_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "phone_number" varchar(50),
  "email" varchar(50)
);

CREATE TABLE "product" (
  "product_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "unit_price" Float,
  "category" varchar(50),
  "warehouse_id" integer,
  "inventory_quantity" integer
);

CREATE TABLE "goods_delivery_note" (
  "goods_delivery_note_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "product_id" integer,
  "amounts" integer,
  "price" Float,
  "exportDate" datetime
  "status" varchar(50)
);

CREATE TABLE "goods_received_note" (
  "goods_received_note_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "product_id" integer,
  "amounts" integer,
  "price" Float,
  "importDate" datetime,
  "supplier_id" integer
);

CREATE TABLE "staff" (
  "staff_id" Serial PRIMARY KEY,
  "name" varchar(50),
  "phone_number" integer,
  "email" varchar(50)
);

ALTER TABLE "product" ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouse" ("warehouse_id");

ALTER TABLE "goods_received_note" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("product_id");

ALTER TABLE "goods_delivery_note" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("product_id");

ALTER TABLE "goods_received_note" ADD FOREIGN KEY ("supplier_id") REFERENCES "supplier" ("supplier_id");

-- +migrate Down
DROP TABLE warehouse;
DROP TABLE supplier;
DROP TABLE staff;
DROP TABLE goods_received_note;
DROP TABLE goods_delivery_note;
DROP TABLE product;