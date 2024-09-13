CREATE DATABASE goproject;

ALTER DATABASE goproject OWNER TO postgres;

\connect goproject

CREATE TABLE users (username text PRIMARY KEY, "password" text);
ALTER TABLE users OWNER TO postgres;


CREATE TABLE usersessions (username text PRIMARY KEY REFERENCES Users(username), token text);
ALTER TABLE usersessions OWNER TO postgres;


CREATE TABLE messages(uuid text PRIMARY KEY, created_at timestamp, sender text REFERENCES Users(username), upvotes integer, downvotes integer, content VARCHAR(1024));
ALTER TABLE messages OWNER TO postgres;