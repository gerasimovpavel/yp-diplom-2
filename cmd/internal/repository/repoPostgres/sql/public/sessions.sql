create table sessions
(
    session_id    uuid default gen_v7_uuid() not null
        constraint sessions_pk
            primary key,
    user_id       uuid                       not null,
    expires_at    timestamp with time zone   not null,
    refresh_token varchar
);

alter table sessions
    owner to postgres;

create index sessions_user_id_idx
    on sessions (user_id);

