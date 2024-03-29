create table if not exists users (
  id bigint primary key,
  first_name text not null,
  last_name text not null,
  username text not null,
  is_premium boolean not null
);

create table if not exists chats (
  id bigint primary key,
  title text not null,
  data jsonb not null default '{"active": true}'
);

create table if not exists chat_members (
  user_id bigint not null references users,
  chat_id bigint not null references chats,
  custom_title text not null,
  primary key (user_id, chat_id)
);

create table if not exists messages (
  id bigserial primary key,
  user_id bigint not null references users,
  chat_id bigint not null references chats,
  content text not null,
  timestamp timestamptz not null default now()
);

create table if not exists stickers (
  message_id bigint primary key references messages,
  file_id text not null
);

create table if not exists commands (
  id bigserial primary key,
  chat_id bigint not null references chats,
  definition text not null,
  substitution_text text not null,
  substitution_photo text not null,
  created_at timestamptz not null default now()
);

create table if not exists handlers (
  message_id bigint primary key references messages,
  handler text not null,
  time interval not null,
  error text not null
);

create table if not exists handlers_info (
  handler text primary key,
  usage text not null
);

create or replace view active_users as (
  select distinct user_id, chat_id
    from messages
   where timestamp > now() - interval '1 week'
);
