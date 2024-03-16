-- name: GetUser :one
select * from users where id = $1;

-- name: UpdateUser :exec
insert into users (id, first_name, last_name, username, is_premium)
values ($1, $2, $3, $4, $5)
on conflict (id)
do update set
  first_name = excluded.first_name,
  last_name = excluded.last_name,
  username = excluded.username,
  is_premium = excluded.is_premium
where users.* != excluded.*;

-- name: RecentUsers :many
select *
  from users
 where id in (
   select user_id
     from messages
    where chat_id = $1
      and timestamp > now() - '1 week'::interval
 );

-- name: GetChat :one
select * from chats where id = $1;

-- name: UpdateChat :exec
insert into chats (id, title)
values ($1, $2)
on conflict (id)
do update set
  title = excluded.title
where chats.* != excluded.*;

-- name: GetChatMember :one
select * from chat_members where user_id = $1 and chat_id = $2;

-- name: UpdateChatMember :exec
insert into chat_members (user_id, chat_id, custom_title)
values ($1, $2, $3)
on conflict (user_id, chat_id)
do update set
  custom_title = excluded.custom_title
where chat_members.* != excluded.*;

-- name: AddMessage :one
insert into messages (user_id, chat_id, content)
values ($1, $2, $3)
returning id;

-- name: ListMessages :many
select m.*
  from messages m
       join handlers h
           on m.id = h.message_id
 where chat_id = $1
   and h.handler = '*handlers.Pass'
   and m.content != '';

-- name: InstrumentMessage :exec
insert into handlers (message_id, handler, time, error)
values ($1, $2, $3, $4);

-- name: AddSticker :exec
insert into stickers (message_id, file_id) values ($1, $2);

-- name: RecentStickers :many
select s.*
  from stickers s
       join messages m
           on s.message_id = m.id
 where m.chat_id = $1
 order by m.timestamp desc
 limit 50;

-- name: RandomHelloSticker :one
select *
  from hello_stickers
 order by random()
 limit 1;

-- name: UpdateDaily :exec
update chats
   set data = $2
 where id = $1
   and coalesce((data ->> 'updated_at')::date < current_date, true);

-- name: SetChatStatus :exec
update chats
   set data['active'] = to_jsonb(@active::boolean)
 where id = $1;

-- name: AddCommand :exec
insert into commands (chat_id, definition, substitution_text, substitution_photo)
values ($1, $2, $3, $4);

-- name: ListCommands :many
select * from commands where chat_id = $1;

-- name: DeleteCommands :exec
delete from commands where chat_id = $1 and definition = $2;
