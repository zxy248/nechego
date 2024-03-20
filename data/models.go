// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package data

import (
	"time"
)

type ActiveUser struct {
	UserID int64
	ChatID int64
}

type Chat struct {
	ID    int64
	Title string
	Data  ChatData
}

type ChatMember struct {
	UserID      int64
	ChatID      int64
	CustomTitle string
}

type Command struct {
	ID                int64
	ChatID            int64
	Definition        string
	SubstitutionText  string
	SubstitutionPhoto string
}

type Handler struct {
	MessageID int64
	Handler   string
	Time      time.Duration
	Error     string
}

type Message struct {
	ID        int64
	UserID    int64
	ChatID    int64
	Content   string
	Timestamp time.Time
}

type Sticker struct {
	MessageID int64
	FileID    string
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	IsPremium bool
}
