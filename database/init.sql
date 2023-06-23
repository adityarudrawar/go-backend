-- https://cadu.dev/creating-a-docker-image-with-database-preloaded/

CREATE DATABASE nimble;

ALTER DATABASE nimble OWNER TO postgres;

\connect nimble

CREATE TABLE Users (id text PRIMARY KEY, username text UNIQUE, hashed_password text);

CREATE TABLE Messages(id text, created_at DATE, sender_id text, upvotes integer, downvotes integer, content text, sender_name text, PRIMARY KEY(id));

ALTER TABLE Users OWNER TO postgres;
ALTER TABLE Messages OWNER TO postgres;

INSERT INTO Users VALUES
(101, 'aditya1', 'password'),
(102, 'aditya2', 'password'),
(103, 'aditya3', 'password'),
(104, 'aditya4', 'password'),
(105, 'aditya5', 'password'),
(106, 'aditya6', 'password'),
(107, 'aditya7', 'password'),
(108, 'aditya8', 'password');