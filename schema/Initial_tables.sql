CREATE TABLE users
(
    id         bigserial primary key not null,
    email      varchar(50)           not null unique,
    password   varchar(255)          not null,
    name       varchar(50)           not null,
    last_name  varchar(50)           not null,
    last_login timestamptz           null,
    created_at timestamptz           not null,
    updated_at timestamptz           not null
);

CREATE TABLE pets
(
    id            bigserial primary key not null,
    user_id       int                   not null,
    name          varchar(50)           not null,
    age           int                   not null,
    breed         varchar(30)           not null,
    birth_date    date                  not null,
    register_date timestamptz           not null,
    updated_at    timestamptz           not null,
    constraint fk_user_pets foreign key (user_id) references users (id)
        on delete cascade on update cascade
);