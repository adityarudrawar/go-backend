-- https://cadu.dev/creating-a-docker-image-with-database-preloaded/

CREATE DATABASE nimble;

ALTER DATABASE nimble OWNER TO postgres;

\connect nimble

CREATE TABLE Users (id integer PRIMARY KEY, username text, hashed_password text);

CREATE TABLE Friends(id integer, friend_with integer);

CREATE TABLE Message(id integer, created_at DATE, sender integer, receiver integer, upvotes integer, downvotes integer, content text, PRIMARY KEY(id));

ALTER TABLE Users OWNER TO postgres;
ALTER TABLE Friends OWNER TO postgres;
ALTER TABLE Message OWNER TO postgres;

INSERT INTO Users VALUES
(101, 'aditya1', 'password'),
(102, 'aditya2', 'password'),
(103, 'aditya3', 'password'),
(104, 'aditya4', 'password'),
(105, 'aditya5', 'password'),
(106, 'aditya6', 'password'),
(107, 'aditya7', 'password'),
(108, 'aditya8', 'password');

INSERT INTO Friends VALUES
(101, 102),
(102, 101),
(101, 103),
(103, 101);
