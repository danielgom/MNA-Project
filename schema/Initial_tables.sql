CREATE TABLE users
(
    id         bigserial primary key not null,
    email      varchar(50)           not null unique,
    password   varchar(255)          not null,
    name       varchar(50)           not null,
    address    varchar(200)          null,
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
    type          varchar(30)           not null,
    breed         varchar(30)           not null,
    sex           varchar(10)           not null,
    color         varchar(20)           not null,
    birth_date    date                  not null,
    register_date timestamptz           not null,
    updated_at    timestamptz           not null,
    constraint fk_user_pets foreign key (user_id) references users (id)
        on delete cascade on update cascade
);

CREATE TABLE vaccines
(
    id         bigserial primary key not null,
    pet_id     int                   not null,
    type       varchar(50)           not null,
    vet_name   varchar(100)          not null,
    address    varchar(100)          not null,
    date       date                  not null,
    next_date  date                  null,
    updated_at timestamptz           not null,
    constraint fk_pet_vaccines foreign key (pet_id) references pets (id)
        on delete cascade on update cascade
);

CREATE TABLE vet_visits
(
    id         bigserial primary key not null,
    pet_id     int                   not null,
    address    varchar(50)           not null,
    vet_name   varchar(100)          not null,
    reason     varchar(50)           not null,
    comments   text                  null,
    date       date                  not null,
    updated_at timestamptz           not null,
    constraint fk_pet_vet_visits foreign key (pet_id) references pets (id)
        on delete cascade on update cascade
);

CREATE TABLE dewormings
(
    id         bigserial primary key not null,
    pet_id     int                   not null,
    address    varchar(50)           not null,
    vet_name   varchar(100)          not null,
    date       date                  not null,
    next_date  date                  not null,
    updated_at timestamptz           not null,
    constraint fk_pet_deworming foreign key (pet_id) references pets (id)
        on delete cascade on update cascade
);





