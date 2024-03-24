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

-- name: ListUsers :many
select *
  from users
 where (id, @chat_id::bigint)
       in (select user_id, chat_id
             from active_users);

-- name: RandomUsers :many
select *
  from users
 where (id, @chat_id::bigint)
       in (select user_id, chat_id
             from active_users)
 order by random()
 limit $1;

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

-- name: MessageCount :one
select jsonb_agg(q)::text from (
  select to_char(date_trunc('day', timestamp), 'DD.MM') as x,
         count(*) as y
    from messages
   where chat_id = $1
     and timestamp > '2024-03-16'
   group by x
   order by x
) q;

-- name: CommandCount :one
select jsonb_agg(q)::text from (
  select to_char(date_trunc('day', timestamp), 'DD.MM') as x,
         count(*) as y
    from messages m
         left join handlers h
             on m.id = h.message_id
   where chat_id = $1
     and handler <> '*handlers.Pass'
     and timestamp > '2024-03-16'
   group by x
   order by x
) q;

-- name: TopCommands :one
select jsonb_agg(q)::text from (
  select count(*) as x,
         usage as y
    from handlers h
         join messages m
             on h.message_id = m.id
         join handlers_info hi
             on h.handler = hi.handler
   where chat_id = $1
     and h.handler <> '*handlers.Pass'
   group by usage
   order by x desc
) q;

-- name: TopUsers :one
select jsonb_agg(q)::text from (
  select count(*) as x,
         format_name(user_id, chat_id) as y
    from messages
   where chat_id = $1
   group by y
   order by x desc
) q;
