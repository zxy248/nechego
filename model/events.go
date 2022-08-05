package model

type event int

const (
	parliamentMemberEvent event = iota
	impeachmentEvent
	depositEvent
)

const selectEvents = `
select u.* from real_users as u
inner join events as e on u.id = e.user_id
where e.event = ? and e.gid = ?
and e.happen >= date('now', 'localtime')`

const insertEvent = `
insert into events (gid, user_id, event, happen)
values (?, ?, ?, datetime('now', 'localtime'))`

const countEventsToday = `
select count(1) from events
where event = ? and gid = ?
and happen >= date('now', 'localtime')`

const countUserEventsToday = `
select count(1) from events
where event = ? and gid = ? and user_id = ?
and happen >= date('now', 'localtime')`

const existsEventToday = `
select exists(
    select 1 from events
    where event = ? and gid = ?
    and happen >= date('now', 'localtime')
    limit 1)`

const existsUserEventToday = `
select exists(
    select 1 from events
    where event = ? and gid = ? and user_id = ?
    and happen >= date('now', 'localtime')
    limit 1)`
