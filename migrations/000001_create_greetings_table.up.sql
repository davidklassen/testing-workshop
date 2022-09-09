create extension if not exists "uuid-ossp";

create table if not exists greetings
(
    id uuid primary key,
    template text
);

insert into greetings(id, template)
values (uuid_generate_v4(), 'Hello, {{.}}'),
       (uuid_generate_v4(), 'Hi, {{.}}');
