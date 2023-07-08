create table users
(
    id           serial primary key,
    full_name    varchar(50) not null,
    username     varchar(15) not null,
    email        varchar(30) not null,
    phone_number varchar(20)
)