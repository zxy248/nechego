package sqlite

import (
	"database/sql"
	"errors"
	"nechego/model"
)

type Energy struct {
	DB *DB
}

const energyEnergyQuery = "select energy from users where gid = ? and uid = ?"

// Energy returns an energy value of the user.
func (e *Energy) Energy(gid, uid int64) (int, error) {
	var energy int
	if err := e.DB.QueryRow(energyEnergyQuery, gid, uid).Scan(&energy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, model.ErrNoUser
		}
		return 0, err
	}
	return energy, nil
}

const energyUpdateQuery = "update users set energy = energy + ? where gid = ? and uid = ?"

// Restore updates an energy value of the user by delta.
func (e *Energy) Update(gid, uid int64, delta int) error {
	_, err := e.DB.Exec(energyUpdateQuery, delta, gid, uid)
	return err
}
