create table "user"
(
    id           uuid not null
        constraint user_pk_2
            primary key,
    phone_number varchar(25)
        constraint user_pk
            unique,
    first_name   varchar(20),
    last_name    varchar(20),
    address      text,
    created_at   timestamp,
    updated_at   timestamp,
    version      integer,
    pin          varchar(255),
    balance      integer
);

alter table "user"
    owner to postgres;

create table transaction
(
    id               uuid not null,
    remarks          varchar(59),
    amount           integer,
    balance_before   integer,
    balance_after    integer,
    transaction_type varchar(6),
    user_id          uuid
        constraint transaction_user_id_fk
            references "user",
    created_at       timestamp,
    version          integer
);

alter table transaction
    owner to postgres;

