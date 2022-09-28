DROP DATABASE IF EXISTS ecommerce_service_test;    

CREATE DATABASE ecommerce_service_test WITH OWNER postgres;

\c ecommerce_service_test;        

CREATE TABLE IF NOT EXISTS public.customers
(
    customer_number integer NOT NULL,
    customer_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_last_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_first_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line1 character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line2 character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    state character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    postal_code character varying(15) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    country character varying(50) COLLATE pg_catalog."default" NOT NULL,
    sales_rep_employee_number integer,
    credit_limit numeric(10,2) DEFAULT NULL::numeric
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.customers
    OWNER to postgres;

GRANT ALL ON TABLE public.customers TO postgres;

CREATE TABLE IF NOT EXISTS public.customers
(
    customer_number integer NOT NULL,
    customer_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_last_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    contact_first_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line1 character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line2 character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    state character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    postal_code character varying(15) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    country character varying(50) COLLATE pg_catalog."default" NOT NULL,
    sales_rep_employee_number integer,
    credit_limit numeric(10,2) DEFAULT NULL::numeric
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.customers
    OWNER to postgres;

GRANT ALL ON TABLE public.customers TO postgres;


CREATE TABLE IF NOT EXISTS public.offices
(
    office_code character varying(10) COLLATE pg_catalog."default" NOT NULL,
    city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line1 character varying(50) COLLATE pg_catalog."default" NOT NULL,
    address_line2 character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    state character varying(50) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    country character varying(50) COLLATE pg_catalog."default" NOT NULL,
    postal_code character varying(15) COLLATE pg_catalog."default" NOT NULL,
    territory character varying(10) COLLATE pg_catalog."default" NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.offices
    OWNER to postgres;

GRANT ALL ON TABLE public.offices TO postgres;

CREATE TABLE IF NOT EXISTS public.orderdetails
(
    order_number integer NOT NULL,
    product_code character varying(15) COLLATE pg_catalog."default" NOT NULL,
    quantity_ordered integer NOT NULL,
    price_each numeric(10,2) NOT NULL,
    order_line_number smallint NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orderdetails
    OWNER to postgres;

GRANT ALL ON TABLE public.orderdetails TO postgres;

CREATE TABLE IF NOT EXISTS public.orders
(
    order_number integer NOT NULL,
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
CREATE TABLE IF NOT EXISTS public.payments
(
    customer_number integer NOT NULL,
    check_number character varying(50) COLLATE pg_catalog."default" NOT NULL,
    payment_date date NOT NULL,
    amount numeric(10,2) NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.payments
    OWNER to postgres;

GRANT ALL ON TABLE public.payments TO postgres;
CREATE TABLE IF NOT EXISTS public.product_lines
(
    product_line character varying(50) COLLATE pg_catalog."default" NOT NULL,
    text_description character varying(4000) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    html_description text COLLATE pg_catalog."default",
    image bytea
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.product_lines
    OWNER to postgres;

GRANT ALL ON TABLE public.product_lines TO postgres;
CREATE TABLE IF NOT EXISTS public.products
(
    product_code character varying(15) COLLATE pg_catalog."default" NOT NULL,
    product_name character varying(70) COLLATE pg_catalog."default" NOT NULL,
    product_line character varying(50) COLLATE pg_catalog."default" NOT NULL,
    product_scale character varying(10) COLLATE pg_catalog."default" NOT NULL,
    product_vendor character varying(50) COLLATE pg_catalog."default" NOT NULL,
    product_description text COLLATE pg_catalog."default" NOT NULL,
    quantity_in_stock smallint NOT NULL,
    buy_price numeric(10,2) NOT NULL,
    msrp numeric(10,2) NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.products
    OWNER to postgres;

GRANT ALL ON TABLE public.products TO postgres;



