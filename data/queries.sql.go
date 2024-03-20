// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package data

import (
	"context"

	"time"
)

const addCommand = `-- name: AddCommand :exec
insert into commands (chat_id, definition, substitution_text, substitution_photo)
values ($1, $2, $3, $4)
`

type AddCommandParams struct {
	ChatID            int64
	Definition        string
	SubstitutionText  string
	SubstitutionPhoto string
}

func (q *Queries) AddCommand(ctx context.Context, arg AddCommandParams) error {
	_, err := q.db.Exec(ctx, addCommand,
		arg.ChatID,
		arg.Definition,
		arg.SubstitutionText,
		arg.SubstitutionPhoto,
	)
	return err
}

const addMessage = `-- name: AddMessage :one
insert into messages (user_id, chat_id, content)
values ($1, $2, $3)
returning id
`

type AddMessageParams struct {
	UserID  int64
	ChatID  int64
	Content string
}

func (q *Queries) AddMessage(ctx context.Context, arg AddMessageParams) (int64, error) {
	row := q.db.QueryRow(ctx, addMessage, arg.UserID, arg.ChatID, arg.Content)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const addSticker = `-- name: AddSticker :exec
insert into stickers (message_id, file_id) values ($1, $2)
`

type AddStickerParams struct {
	MessageID int64
	FileID    string
}

func (q *Queries) AddSticker(ctx context.Context, arg AddStickerParams) error {
	_, err := q.db.Exec(ctx, addSticker, arg.MessageID, arg.FileID)
	return err
}

const deleteCommands = `-- name: DeleteCommands :exec
delete from commands where chat_id = $1 and definition = $2
`

type DeleteCommandsParams struct {
	ChatID     int64
	Definition string
}

func (q *Queries) DeleteCommands(ctx context.Context, arg DeleteCommandsParams) error {
	_, err := q.db.Exec(ctx, deleteCommands, arg.ChatID, arg.Definition)
	return err
}

const getChat = `-- name: GetChat :one
select id, title, data from chats where id = $1
`

func (q *Queries) GetChat(ctx context.Context, id int64) (Chat, error) {
	row := q.db.QueryRow(ctx, getChat, id)
	var i Chat
	err := row.Scan(&i.ID, &i.Title, &i.Data)
	return i, err
}

const getChatMember = `-- name: GetChatMember :one
select user_id, chat_id, custom_title from chat_members where user_id = $1 and chat_id = $2
`

type GetChatMemberParams struct {
	UserID int64
	ChatID int64
}

func (q *Queries) GetChatMember(ctx context.Context, arg GetChatMemberParams) (ChatMember, error) {
	row := q.db.QueryRow(ctx, getChatMember, arg.UserID, arg.ChatID)
	var i ChatMember
	err := row.Scan(&i.UserID, &i.ChatID, &i.CustomTitle)
	return i, err
}

const getUser = `-- name: GetUser :one
select id, first_name, last_name, username, is_premium from users where id = $1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.IsPremium,
	)
	return i, err
}

const instrumentMessage = `-- name: InstrumentMessage :exec
insert into handlers (message_id, handler, time, error)
values ($1, $2, $3, $4)
`

type InstrumentMessageParams struct {
	MessageID int64
	Handler   string
	Time      time.Duration
	Error     string
}

func (q *Queries) InstrumentMessage(ctx context.Context, arg InstrumentMessageParams) error {
	_, err := q.db.Exec(ctx, instrumentMessage,
		arg.MessageID,
		arg.Handler,
		arg.Time,
		arg.Error,
	)
	return err
}

const listCommands = `-- name: ListCommands :many
select id, chat_id, definition, substitution_text, substitution_photo from commands where chat_id = $1
`

func (q *Queries) ListCommands(ctx context.Context, chatID int64) ([]Command, error) {
	rows, err := q.db.Query(ctx, listCommands, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Command
	for rows.Next() {
		var i Command
		if err := rows.Scan(
			&i.ID,
			&i.ChatID,
			&i.Definition,
			&i.SubstitutionText,
			&i.SubstitutionPhoto,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMessages = `-- name: ListMessages :many
select m.id, m.user_id, m.chat_id, m.content, m.timestamp
  from messages m
       join handlers h
           on m.id = h.message_id
 where chat_id = $1
   and h.handler = '*handlers.Pass'
   and m.content != ''
`

func (q *Queries) ListMessages(ctx context.Context, chatID int64) ([]Message, error) {
	rows, err := q.db.Query(ctx, listMessages, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ChatID,
			&i.Content,
			&i.Timestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
select id, first_name, last_name, username, is_premium
  from users
 where (id, $1::bigint)
       in (select user_id, chat_id
             from active_users)
`

func (q *Queries) ListUsers(ctx context.Context, chatID int64) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Username,
			&i.IsPremium,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const randomUsers = `-- name: RandomUsers :many
select id, first_name, last_name, username, is_premium
  from users
 where (id, $2::bigint)
       in (select user_id, chat_id
             from active_users)
 order by random()
 limit $1
`

type RandomUsersParams struct {
	Limit  int32
	ChatID int64
}

func (q *Queries) RandomUsers(ctx context.Context, arg RandomUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, randomUsers, arg.Limit, arg.ChatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Username,
			&i.IsPremium,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const recentStickers = `-- name: RecentStickers :many
select s.message_id, s.file_id
  from stickers s
       join messages m
           on s.message_id = m.id
 where m.chat_id = $1
 order by m.timestamp desc
 limit 50
`

func (q *Queries) RecentStickers(ctx context.Context, chatID int64) ([]Sticker, error) {
	rows, err := q.db.Query(ctx, recentStickers, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Sticker
	for rows.Next() {
		var i Sticker
		if err := rows.Scan(&i.MessageID, &i.FileID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setChatStatus = `-- name: SetChatStatus :exec
update chats
   set data['active'] = to_jsonb($2::boolean)
 where id = $1
`

type SetChatStatusParams struct {
	ID     int64
	Active bool
}

func (q *Queries) SetChatStatus(ctx context.Context, arg SetChatStatusParams) error {
	_, err := q.db.Exec(ctx, setChatStatus, arg.ID, arg.Active)
	return err
}

const updateChat = `-- name: UpdateChat :exec
insert into chats (id, title)
values ($1, $2)
on conflict (id)
do update set
  title = excluded.title
where chats.* != excluded.*
`

type UpdateChatParams struct {
	ID    int64
	Title string
}

func (q *Queries) UpdateChat(ctx context.Context, arg UpdateChatParams) error {
	_, err := q.db.Exec(ctx, updateChat, arg.ID, arg.Title)
	return err
}

const updateChatMember = `-- name: UpdateChatMember :exec
insert into chat_members (user_id, chat_id, custom_title)
values ($1, $2, $3)
on conflict (user_id, chat_id)
do update set
  custom_title = excluded.custom_title
where chat_members.* != excluded.*
`

type UpdateChatMemberParams struct {
	UserID      int64
	ChatID      int64
	CustomTitle string
}

func (q *Queries) UpdateChatMember(ctx context.Context, arg UpdateChatMemberParams) error {
	_, err := q.db.Exec(ctx, updateChatMember, arg.UserID, arg.ChatID, arg.CustomTitle)
	return err
}

const updateDaily = `-- name: UpdateDaily :exec
update chats
   set data = $2
 where id = $1
   and coalesce((data ->> 'updated_at')::date < current_date, true)
`

type UpdateDailyParams struct {
	ID   int64
	Data ChatData
}

func (q *Queries) UpdateDaily(ctx context.Context, arg UpdateDailyParams) error {
	_, err := q.db.Exec(ctx, updateDaily, arg.ID, arg.Data)
	return err
}

const updateUser = `-- name: UpdateUser :exec
insert into users (id, first_name, last_name, username, is_premium)
values ($1, $2, $3, $4, $5)
on conflict (id)
do update set
  first_name = excluded.first_name,
  last_name = excluded.last_name,
  username = excluded.username,
  is_premium = excluded.is_premium
where users.* != excluded.*
`

type UpdateUserParams struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	IsPremium bool
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Username,
		arg.IsPremium,
	)
	return err
}
