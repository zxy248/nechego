package model

import (
	"errors"
	"nechego/input"
	"time"
)

var ErrNoEblan = errors.New("no eblan")
var ErrNoAdmin = errors.New("no admin")
var ErrNoPair = errors.New("no pair")
var ErrNoUser = errors.New("no user")

type Model struct {
	Admins    AdminsModel
	Bans      BansModel
	Eblans    EblansModel
	Forbid    ForbidModel
	Pairs     PairsModel
	Status    StatusModel
	Users     UsersModel
	Whitelist WhitelistModel
	Messages  MessagesModel
	Energy    EnergyModel
}

type AdminsModel interface {
	Insert(int64) error
	Delete(int64) error
	List(int64) ([]int64, error)
	Authorize(int64, int64) (bool, error)
	InsertDaily(int64, int64) error
	GetDaily(int64) (int64, error)
	DeleteDaily(int64) error
}

type BansModel interface {
	Ban(int64) error
	Unban(int64) error
	List() ([]int64, error)
	Banned(int64) (bool, error)
}

type EblansModel interface {
	Insert(int64, int64) error
	Get(int64) (int64, error)
	Delete(int64) error
}

type ForbidModel interface {
	Forbid(int64, input.Command) error
	Permit(int64, input.Command) error
	Forbidden(int64, input.Command) (bool, error)
	List(int64) ([]input.Command, error)
}

type PairsModel interface {
	Insert(int64, int64, int64) error
	Get(int64) (int64, int64, error)
}

type StatusModel interface {
	Enable(int64) error
	Active(int64) (bool, error)
	Disable(int64) error
}

type User struct {
	GID    int64
	UID    int64
	Energy int
}

type UsersModel interface {
	Insert(int64, int64) error
	Delete(int64, int64) error
	List(int64) ([]int64, error)
	All() ([]User, error)
	Exists(int64, int64) (bool, error)
	Random(int64) (int64, error)
	NRandom(int64, int) ([]int64, error)
}

type WhitelistModel interface {
	Insert(int64) error
	Allow(int64) (bool, error)
}

type MessagesModel interface {
	UserCount(int64, int64, time.Time) (int, error)
	TotalCount(int64, time.Time) (int, error)
}

type EnergyModel interface {
	Energy(int64, int64) (int, error)
	Update(int64, int64, int) error
}
