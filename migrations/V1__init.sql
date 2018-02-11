--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.10
-- Dumped by pg_dump version 9.5.10

-- Started on 2018-02-12 00:37:00 MSK

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

DROP INDEX IF EXISTS public.tournaments_title_index;
DROP INDEX IF EXISTS public.tournaments_id_index;
DROP INDEX IF EXISTS public.teams_team_name_index;
DROP INDEX IF EXISTS public.teams_id_index;
DROP INDEX IF EXISTS public.persons_id_index;
DROP INDEX IF EXISTS public.matches_tourn_id_id_index;
DROP INDEX IF EXISTS public.games_title_index;
DROP INDEX IF EXISTS public.games_id_index;
DROP TABLE IF EXISTS public.tournaments;
DROP TABLE IF EXISTS public.timeline;
DROP TABLE IF EXISTS public.teams;
DROP TABLE IF EXISTS public.players;
DROP TABLE IF EXISTS public.persons;
DROP TABLE IF EXISTS public.matches;
DROP TABLE IF EXISTS public.games;
DROP TABLE IF EXISTS public.game_tourney;
DROP TABLE IF EXISTS public.auth;
DROP EXTENSION IF EXISTS plpgsql;
DROP SCHEMA IF EXISTS public;
--
-- TOC entry 6 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- TOC entry 2190 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 1 (class 3079 OID 12397)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2192 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 188 (class 1259 OID 36519)
-- Name: auth; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE auth (
    email character varying(50) NOT NULL,
    pass bytea NOT NULL,
    person_id uuid NOT NULL,
    staff boolean DEFAULT false NOT NULL
);


--
-- TOC entry 189 (class 1259 OID 36632)
-- Name: game_tourney; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE game_tourney (
    game_id uuid NOT NULL,
    tourney_id uuid,
    started timestamp without time zone NOT NULL,
    title character varying(50) NOT NULL
);


--
-- TOC entry 182 (class 1259 OID 36203)
-- Name: games; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE games (
    title character varying(50) NOT NULL,
    about text NOT NULL,
    id uuid NOT NULL,
    CONSTRAINT not_empty CHECK (((title)::bpchar <> ''::bpchar))
);


--
-- TOC entry 184 (class 1259 OID 36240)
-- Name: matches; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE matches (
    id uuid NOT NULL,
    tourn_id uuid NOT NULL,
    team_id_1 uuid,
    team_id_2 uuid,
    team_score_1 integer DEFAULT 0 NOT NULL,
    team_score_2 integer DEFAULT 0 NOT NULL,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone,
    link character varying(300) DEFAULT 'unavailable'::character varying NOT NULL,
    prev_match_id_1 uuid,
    prev_match_id_2 uuid,
    next_match_id uuid,
    organize_id uuid NOT NULL
);


--
-- TOC entry 185 (class 1259 OID 36247)
-- Name: persons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE persons (
    id uuid NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    about text
);


--
-- TOC entry 186 (class 1259 OID 36254)
-- Name: players; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE players (
    id uuid NOT NULL,
    person_id uuid NOT NULL,
    nickname character varying(50) NOT NULL,
    team_id uuid NOT NULL,
    team_name character varying(50),
    retire boolean DEFAULT false NOT NULL
);


--
-- TOC entry 181 (class 1259 OID 36124)
-- Name: teams; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE teams (
    team_name character varying(50) NOT NULL,
    about text NOT NULL,
    id uuid NOT NULL,
    coach_id uuid NOT NULL,
    coach_name character varying(101) NOT NULL
);


--
-- TOC entry 187 (class 1259 OID 36258)
-- Name: timeline; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE timeline (
    match_id uuid NOT NULL,
    event_time timestamp without time zone NOT NULL,
    about text NOT NULL
);


--
-- TOC entry 183 (class 1259 OID 36234)
-- Name: tournaments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE tournaments (
    id uuid NOT NULL,
    started timestamp without time zone NOT NULL,
    ended timestamp without time zone NOT NULL,
    about text NOT NULL,
    title character varying(50) NOT NULL,
    organize_id uuid NOT NULL,
    organize_name character varying(101) NOT NULL,
    game_id uuid NOT NULL
);


--
-- TOC entry 2065 (class 1259 OID 36558)
-- Name: games_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX games_id_index ON games USING btree (id);


--
-- TOC entry 2066 (class 1259 OID 36296)
-- Name: games_title_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX games_title_index ON games USING btree (title);


--
-- TOC entry 2069 (class 1259 OID 36362)
-- Name: matches_tourn_id_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX matches_tourn_id_id_index ON matches USING btree (tourn_id, id);


--
-- TOC entry 2070 (class 1259 OID 36340)
-- Name: persons_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX persons_id_index ON persons USING btree (id);


--
-- TOC entry 2063 (class 1259 OID 36316)
-- Name: teams_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX teams_id_index ON teams USING btree (id);


--
-- TOC entry 2064 (class 1259 OID 36315)
-- Name: teams_team_name_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX teams_team_name_index ON teams USING btree (team_name);


--
-- TOC entry 2067 (class 1259 OID 36361)
-- Name: tournaments_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tournaments_id_index ON tournaments USING btree (id);


--
-- TOC entry 2068 (class 1259 OID 36360)
-- Name: tournaments_title_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tournaments_title_index ON tournaments USING btree (title);


--
-- TOC entry 2191 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: -
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2018-02-12 00:37:01 MSK

--
-- PostgreSQL database dump complete
--

