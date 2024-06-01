CREATE EXTENSION IF NOT EXISTS plpgsql;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users
(
    user_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    full_name text COLLATE pg_catalog."default",
    email text COLLATE pg_catalog."default",
    password text COLLATE pg_catalog."default",
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT validate_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z]{2,}$'::text)
)

CREATE TABLE IF NOT EXISTS public.warehouse
(
    warehouse_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) COLLATE pg_catalog."default",
    location character varying(50) COLLATE pg_catalog."default",
    capacity integer,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT warehouse_pkey PRIMARY KEY (warehouse_id),
    CONSTRAINT check_capacity_gt_zero CHECK (capacity > 0)
)

TABLESPACE pg_default;

ALTER TABLE public.warehouse
    OWNER to postgres;


CREATE TABLE IF NOT EXISTS public.product
(
    product_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) COLLATE pg_catalog."default",
    unit_price double precision,
    category character varying(50) COLLATE pg_catalog."default",
    warehouse_id uuid,
    inventory_quantity integer,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT product_pkey PRIMARY KEY (product_id),
    CONSTRAINT product_warehouse_id_fkey FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse (warehouse_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.product
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.supplier
(
    supplier_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) COLLATE pg_catalog."default",
    phone_number character varying(50) COLLATE pg_catalog."default",
    email character varying(50) COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT supplier_pkey PRIMARY KEY (supplier_id),
    CONSTRAINT valid_phone_number CHECK (phone_number::text ~ '^[0-9]{10}$'::text),
    CONSTRAINT valid_email CHECK (email::text ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'::text)
)

TABLESPACE pg_default;

ALTER TABLE public.supplier
    OWNER to postgres;


CREATE TABLE IF NOT EXISTS public.goods_received_note
(
    goods_received_note_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) COLLATE pg_catalog."default",
    product_id uuid,
    amounts integer,
    price double precision,
    supplier_id uuid,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT goods_received_note_pkey PRIMARY KEY (goods_received_note_id),
    CONSTRAINT goods_received_note_product_id_fkey FOREIGN KEY (product_id)
        REFERENCES public.product (product_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT goods_received_note_supplier_id_fkey FOREIGN KEY (supplier_id)
        REFERENCES public.supplier (supplier_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT check_amount_import CHECK (amounts > 0)
)

TABLESPACE pg_default;

ALTER TABLE public.goods_received_note
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.goods_delivery_note
(
    goods_delivery_note_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) COLLATE pg_catalog."default",
    product_id uuid,
    amounts integer,
    price double precision,
    status character varying(50) COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT goods_delivery_note_pkey PRIMARY KEY (goods_delivery_note_id),
    CONSTRAINT goods_delivery_note_product_id_fkey FOREIGN KEY (product_id)
        REFERENCES public.product (product_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT check_amount_export CHECK (amounts > 0)
)

TABLESPACE pg_default;

ALTER TABLE public.goods_delivery_note
    OWNER to postgres;


CREATE FUNCTION public.delete_products_when_warehouse_deleted()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    -- Xóa toàn bộ các sản phẩm (products) liên quan đến warehouse được xóa
    DELETE FROM product WHERE warehouse_id = OLD.warehouse_id;
    RETURN OLD;
END;
$BODY$;

ALTER FUNCTION public.delete_products_when_warehouse_deleted()
    OWNER TO postgres;



CREATE FUNCTION public.inventory_capacity_export()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    -- Trừ giảm amounts từ inventory của product
    UPDATE product
    SET inventory_quantity = inventory_quantity - NEW.amounts
    WHERE product_id = NEW.product_id;

    -- Cộng amounts vào capacity của warehouse
    UPDATE warehouse
    SET capacity = capacity + NEW.amounts
    WHERE warehouse_id = (SELECT warehouse_id FROM product WHERE product_id = NEW.product_id);

    RETURN NEW;
END;
$BODY$;

ALTER FUNCTION public.inventory_capacity_export()
    OWNER TO postgres;


CREATE FUNCTION public.inventory_capacity_import()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    -- Cộng amounts vào inventory của product
    UPDATE product
    SET inventory_quantity = inventory_quantity + NEW.amounts
    WHERE product_id = NEW.product_id;

    -- Trừ giảm amounts từ capacity của warehouse
    UPDATE warehouse
    SET capacity = capacity - NEW.amounts
    WHERE warehouse_id = (SELECT warehouse_id FROM product WHERE product_id = NEW.product_id);

    RETURN NEW;
END;
$BODY$;

ALTER FUNCTION public.inventory_capacity_import()
    OWNER TO postgres;


CREATE FUNCTION public.update_export_warehouse_product()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    -- Cập nhật giá trị mới trong bảng warehouse
    UPDATE warehouse 
    SET capacity = capacity - (OLD.amounts - NEW.amounts),
        updated_at = CURRENT_TIMESTAMP
    WHERE warehouse_id = (
        SELECT warehouse_id FROM product WHERE product_id = NEW.product_id
    );

    -- Cập nhật giá trị mới trong bảng product
    UPDATE product
    SET inventory_quantity = inventory_quantity + (OLD.amounts - NEW.amounts),
        updated_at = CURRENT_TIMESTAMP
    WHERE product_id = NEW.product_id;

    RETURN NEW;
END;
$BODY$;

ALTER FUNCTION public.update_export_warehouse_product()
    OWNER TO postgres;

-- Trigger: delete_products_trigger
CREATE OR REPLACE TRIGGER delete_products_trigger
    BEFORE DELETE
    ON public.warehouse
    FOR EACH ROW
    EXECUTE FUNCTION public.delete_products_when_warehouse_deleted();


-- Trigger: inventory_capacity_export_trigger
CREATE OR REPLACE TRIGGER inventory_capacity_export_trigger
    AFTER INSERT
    ON public.goods_delivery_note
    FOR EACH ROW
    EXECUTE FUNCTION public.inventory_capacity_export();
-- Trigger: update_export_warehouse_product_trigger
CREATE OR REPLACE TRIGGER update_export_warehouse_product_trigger
    AFTER UPDATE OF amounts
    ON public.goods_delivery_note
    FOR EACH ROW
    EXECUTE FUNCTION public.update_export_warehouse_product();

-- Trigger: inventory_capacity_import_trigger
CREATE OR REPLACE TRIGGER inventory_capacity_import_trigger
    AFTER INSERT
    ON public.goods_received_note
    FOR EACH ROW
    EXECUTE FUNCTION public.inventory_capacity_import();
-- Trigger: update_import_warehouse_product_trigger
CREATE OR REPLACE TRIGGER update_import_warehouse_product_trigger
    AFTER UPDATE OF amounts
    ON public.goods_received_note
    FOR EACH ROW
    EXECUTE FUNCTION public.update_import_warehouse_product();
