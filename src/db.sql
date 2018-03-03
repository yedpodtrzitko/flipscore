--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.7
-- Dumped by pg_dump version 9.5.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: score; Type: TABLE; Schema: public; Owner: yed
--

CREATE TABLE score (
    score integer,
    player character varying(32),
    created_at timestamp without time zone,
    id integer NOT NULL
);


ALTER TABLE score OWNER TO yed;

--
-- Name: score_id_seq; Type: SEQUENCE; Schema: public; Owner: yed
--

CREATE SEQUENCE score_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE score_id_seq OWNER TO yed;

--
-- Name: score_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: yed
--

ALTER SEQUENCE score_id_seq OWNED BY score.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: yed
--

ALTER TABLE ONLY score ALTER COLUMN id SET DEFAULT nextval('score_id_seq'::regclass);


--
-- Name: score_id_seq; Type: SEQUENCE SET; Schema: public; Owner: yed
--

SELECT pg_catalog.setval('score_id_seq', 1, false);


--
-- Name: score_pkey; Type: CONSTRAINT; Schema: public; Owner: yed
--

ALTER TABLE ONLY score
    ADD CONSTRAINT score_pkey PRIMARY KEY (id);
