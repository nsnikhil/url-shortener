create extension if not exists "pgcrypto";

create table if not exists shortener (
    id uuid primary key default gen_random_uuid(),
    url  varchar(500) not null unique,
    urlhash  varchar(20) not null unique,
    createdAt timestamp without time zone default (now() at time zone 'utc'),
    updatedAt timestamp without time zone default (now() at time zone 'utc'),
    CHECK (url <> ''),
    CHECK (urlhash <> '')
);