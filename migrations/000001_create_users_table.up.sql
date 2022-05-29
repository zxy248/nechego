create table if not exists users (
       id integer primary key autoincrement,
       group_id integer not null,
       user_id integer not null
);
