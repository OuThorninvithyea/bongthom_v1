--
-- PostgreSQL database dump
--

\restrict Edp645sITNU2rBleqHEHeORxynVpVAvswblRq0vbG58cDmiJAGHHz2h7wCsj3pk

-- Dumped from database version 18.4 (Debian 18.4-1.pgdg13+1)
-- Dumped by pg_dump version 18.4 (Debian 18.4-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: update_updated_at(); Type: FUNCTION; Schema: public; Owner: outhorninvuth
--

CREATE FUNCTION public.update_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at() OWNER TO outhorninvuth;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: tbl_users; Type: TABLE; Schema: public; Owner: outhorninvuth
--

CREATE TABLE public.tbl_users (
    id bigint NOT NULL,
    user_name character varying(100) NOT NULL,
    password character varying(100) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    email character varying(100),
    role_name character varying(50) DEFAULT 'SuperAdmin'::character varying NOT NULL,
    role_id integer DEFAULT 0 NOT NULL,
    is_admin boolean DEFAULT false,
    login_session character varying(255),
    last_login timestamp without time zone,
    currency_id integer,
    language_id integer,
    status_id integer DEFAULT 1,
    "order" integer,
    created_by bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_by bigint,
    updated_at timestamp without time zone,
    deleted_by bigint,
    deleted_at timestamp without time zone
);


ALTER TABLE public.tbl_users OWNER TO outhorninvuth;

--
-- Name: tbl_users_id_seq; Type: SEQUENCE; Schema: public; Owner: outhorninvuth
--

CREATE SEQUENCE public.tbl_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tbl_users_id_seq OWNER TO outhorninvuth;

--
-- Name: tbl_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: outhorninvuth
--

ALTER SEQUENCE public.tbl_users_id_seq OWNED BY public.tbl_users.id;


--
-- Name: tbl_users id; Type: DEFAULT; Schema: public; Owner: outhorninvuth
--

ALTER TABLE ONLY public.tbl_users ALTER COLUMN id SET DEFAULT nextval('public.tbl_users_id_seq'::regclass);


--
-- Data for Name: tbl_users; Type: TABLE DATA; Schema: public; Owner: outhorninvuth
--

COPY public.tbl_users (id, user_name, password, first_name, last_name, email, role_name, role_id, is_admin, login_session, last_login, currency_id, language_id, status_id, "order", created_by, created_at, updated_by, updated_at, deleted_by, deleted_at) FROM stdin;
2	admin	$2a$10$jgwUpnQQs6xVwoxIIswIXubE7uey4QdXmPtcHuKsS99cYGlIuQJqu	\N	\N	\N	SuperAdmin	0	t	ce10c8ff-bbc7-43cb-937c-b516cdabd097	2026-06-15 03:17:02.735721	\N	\N	1	\N	\N	2026-06-11 05:24:05.008356	\N	2026-06-15 03:17:02.735721	\N	\N
\.


--
-- Name: tbl_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: outhorninvuth
--

SELECT pg_catalog.setval('public.tbl_users_id_seq', 2, true);


--
-- Name: tbl_users tbl_users_pkey; Type: CONSTRAINT; Schema: public; Owner: outhorninvuth
--

ALTER TABLE ONLY public.tbl_users
    ADD CONSTRAINT tbl_users_pkey PRIMARY KEY (id);


--
-- Name: tbl_users tbl_users_user_name_key; Type: CONSTRAINT; Schema: public; Owner: outhorninvuth
--

ALTER TABLE ONLY public.tbl_users
    ADD CONSTRAINT tbl_users_user_name_key UNIQUE (user_name);


--
-- Name: tbl_users trg_tbl_users_updated_at; Type: TRIGGER; Schema: public; Owner: outhorninvuth
--

CREATE TRIGGER trg_tbl_users_updated_at BEFORE UPDATE ON public.tbl_users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();


--
-- PostgreSQL database dump complete
--

\unrestrict Edp645sITNU2rBleqHEHeORxynVpVAvswblRq0vbG58cDmiJAGHHz2h7wCsj3pk

