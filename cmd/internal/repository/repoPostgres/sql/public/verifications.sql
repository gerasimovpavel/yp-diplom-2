create table verifications
(
    user_id    uuid                     not null,
    code       varchar                  not null,
    verifed_at timestamp with time zone,
    expired_at timestamp with time zone not null,
    constraint verifications_pk
        primary key (user_id, code)
);

alter table verifications
    owner to postgres;

