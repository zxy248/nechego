create or replace function format_name(bigint, bigint)
  returns text
  language sql
as $$
  select concat_ws(' ',
                   nullif(u.first_name, ''),
                   nullif(u.last_name, ''),
                   '(' || nullif(cm.custom_title, '') || ')')
  from chat_members cm
  join users u
  on cm.user_id = u.id
  where cm.user_id = $1
  and cm.chat_id = $2;
$$;

create or replace function trim_text(text)
  returns text
  language sql
as $$
  select case when length($1) > 50 then left($1, 50) || '...' else $1 end;
$$;
