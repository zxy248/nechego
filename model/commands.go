package model

import "nechego/input"

const forbidCommand = `
insert into forbidden_commands (gid, command)
values (?, ?)`

func (m *Model) ForbidCommand(g Group, c input.Command) bool {
	n, err := m.db.MustExec(forbidCommand, g.GID, c).RowsAffected()
	failOn(err)
	return n == 1
}

const permitCommand = `
delete from forbidden_commands
where gid = ? and command = ?`

func (m *Model) PermitCommand(g Group, c input.Command) bool {
	n, err := m.db.MustExec(permitCommand, g.GID, c).RowsAffected()
	failOn(err)
	return n == 1
}

const forbiddenCommands = `
select command
from forbidden_commands
where gid = ?`

func (m *Model) ForbiddenCommands(g Group) ([]input.Command, error) {
	commands := []input.Command{}
	err := m.db.Select(&commands, forbiddenCommands, g.GID)
	return commands, err
}