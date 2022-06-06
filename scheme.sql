create table if not exists "users" (
    "id" integer,
    "group_id" integer not null,
    "user_id" integer not null,
    primary key ("id" autoincrement)
);

create table if not exists "pairs" (
    "id" integer not null,
    "group_id" integer not null,
    "user_id_x" integer not null,
    "user_id_y" integer not null,
    "last" text not null,
    primary key ("id" autoincrement)
);

create table if not exists "eblans" (
    "id" integer not null,
    "group_id" integer not null,
    "user_id" integer not null,
    "last" text not null,
    primary key ("id" autoincrement)
);
