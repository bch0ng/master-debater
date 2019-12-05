--CREATE DATABASE dockerdb;

create database if not EXISTS psql_db;
  use psql_db;
create table if not exists contacts (
    id SERIAL primary key,
    email varchar(320) not null UNIQUE,
    passhash varchar(255) not null,
    user_name varchar(255) not null UNIQUE,
    first_name varchar(128) not null,
    last_name varchar(128) not null,
    photo_url varchar(128) not null
);
create table if not exists sessions (
    id SERIAL primary key,
    sign_in_time timestamp not null,
    ip varchar(15) not null UNIQUE
);
