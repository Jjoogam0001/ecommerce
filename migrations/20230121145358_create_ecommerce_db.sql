-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.customers
(
    customer_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_last_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_first_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line1 character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line2 character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    city character varying(50) COLLATE pg_catalog."default",
    state character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    postal_code character varying COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    country character varying(50) COLLATE pg_catalog."default",
    sales_rep_employee_number integer,
    credit_limit numeric(10,2) DEFAULT NULL::numeric,
    user_id character varying(255) COLLATE pg_catalog."default",
    server_id character varying COLLATE pg_catalog."default",
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    customer_number bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    CONSTRAINT customers_pkey PRIMARY KEY (email)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.customers
    OWNER to postgres;

GRANT ALL ON TABLE public.customers TO postgres;

-- Table: public.employees

-- DROP TABLE IF EXISTS public.employees;

CREATE TABLE IF NOT EXISTS public.employees
(
    employee_number bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    last_name character varying(255) COLLATE pg_catalog."default",
    first_name character varying(255) COLLATE pg_catalog."default",
    extension character varying(255) COLLATE pg_catalog."default",
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    office_code character varying(255) COLLATE pg_catalog."default",
    reports_to character varying COLLATE pg_catalog."default",
    job_title character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT employees_pk PRIMARY KEY (email)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.employees
    OWNER to postgres;

-- Table: public.offices

-- DROP TABLE IF EXISTS public.offices;

CREATE TABLE IF NOT EXISTS public.offices
(
    office_code bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    city character varying(255) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(255) COLLATE pg_catalog."default" NOT NULL,
    address_line1 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    address_line2 character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT NULL::character varying,
    state character varying(255) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    country character varying(255) COLLATE pg_catalog."default" NOT NULL,
    postal_code character varying(255) COLLATE pg_catalog."default",
    territory character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT offices_pk PRIMARY KEY (office_code)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.offices
    OWNER to postgres;

GRANT ALL ON TABLE public.offices TO postgres;

    -- Table: public.orderdetails

-- DROP TABLE IF EXISTS public.orderdetails;

CREATE TABLE IF NOT EXISTS public.orderdetails
(
    order_number bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    product_code character varying(15) COLLATE pg_catalog."default" NOT NULL,
    quantity_ordered integer NOT NULL,
    price_each integer NOT NULL,
    order_line_number character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT orderdetails_pk PRIMARY KEY (order_number)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orderdetails
    OWNER to postgres;

GRANT ALL ON TABLE public.orderdetails TO postgres;

-- Table: public.orders

-- DROP TABLE IF EXISTS public.orders;

CREATE TABLE IF NOT EXISTS public.orders
(
    order_number bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    order_date date NOT NULL,
    required_date date NOT NULL,
    shipped_date date,
    status character varying(15) COLLATE pg_catalog."default" NOT NULL,
    comments text COLLATE pg_catalog."default",
    customer_number integer NOT NULL
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;

GRANT ALL ON TABLE public.orders TO postgres;

-- Table: public.payments

-- DROP TABLE IF EXISTS public.payments;

CREATE TABLE IF NOT EXISTS public.payments
(
    customer_number bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    check_number character varying(50) COLLATE pg_catalog."default" NOT NULL,
    payment_date date NOT NULL,
    amount numeric(10,2) NOT NULL,
    CONSTRAINT payments_pk PRIMARY KEY (customer_number)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.payments
    OWNER to postgres;

GRANT ALL ON TABLE public.payments TO postgres;

-- Table: public.product_lines

-- DROP TABLE IF EXISTS public.product_lines;

CREATE TABLE IF NOT EXISTS public.product_lines
(
    product_line character varying(255) COLLATE pg_catalog."default" NOT NULL,
    text_description character varying(4000) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    html_description text COLLATE pg_catalog."default",
    image character varying COLLATE pg_catalog."default"
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.product_lines
    OWNER to postgres;

GRANT ALL ON TABLE public.product_lines TO postgres;

-- Table: public.products

-- DROP TABLE IF EXISTS public.products;

CREATE TABLE IF NOT EXISTS public.products
(
    product_code character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_line character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_scale character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_vendor character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_description text COLLATE pg_catalog."default" NOT NULL,
    quantity_in_stock smallint NOT NULL,
    buy_price integer NOT NULL,
    msrp character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT products_pk PRIMARY KEY (product_code)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.products
    OWNER to postgres;

GRANT ALL ON TABLE public.products TO postgres;
-- +goose StatementEnd


