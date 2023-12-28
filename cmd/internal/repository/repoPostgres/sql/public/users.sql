create table users
(
    user_id    uuid default gen_v7_uuid() not null
        constraint users_pk
            primary key,
    last_name  varchar,
    first_name varchar,
    email      varchar                    not null,
    password   varchar
);

alter table users
    owner to postgres;

create unique index users_email_idx
    on users (email);

