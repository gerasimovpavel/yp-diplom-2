create table club_users
(
    club_id uuid not null,
    user_id uuid not null,
    constraint club_users_pk
        primary key (club_id, user_id)
);

alter table club_users
    owner to postgres;

create index club_users_club_id_idx
    on club_users (club_id);

create index club_users_user_id_idx
    on club_users (user_id);

