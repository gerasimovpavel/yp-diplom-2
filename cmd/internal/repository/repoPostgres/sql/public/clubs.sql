create table clubs
(
    club_id      uuid    default gen_v7_uuid()         not null
        constraint clubs_pk
            primary key,
    name         varchar default ''::character varying not null,
    passwordhash varchar default ''::character varying not null
);

alter table clubs
    owner to postgres;

