package app

import (
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

const statusDir = "status"

func getStatus(uid int64) string {
	f, err := os.ReadFile(statusPath(uid))
	if err != nil {
		return ""
	}
	return string(f)
}

func setStatus(uid int64, s string) error {
	if err := os.MkdirAll(statusDir, 0777); err != nil {
		return err
	}
	return os.WriteFile(statusPath(uid), []byte(s), 0666)
}

func statusPath(uid int64) string {
	return filepath.Join(statusDir, strconv.FormatInt(uid, 10))
}

const maxStatusLength = 120

var (
	statusSet         = Response("✅ Статус установлен.")
	errStatusLength   = UserError("Максимальная длина статуса %v символов.")
	errStatusNotFound = UserError("Статус не найден.")
)

func handleStatus(c tele.Context) error {
	u := getUser(c)
	a := getMessage(c).Argument()
	if a != "" {
		if utf8.RuneCountInString(a) > maxStatusLength {
			return respondUserError(c, errStatusLength.Fill(maxStatusLength))
		}
		if err := setStatus(u.UID, a); err != nil {
			return respondInternalError(c, err)
		}
		return respond(c, statusSet)
	}
	s := getStatus(u.UID)
	if s == "" {
		return respondUserError(c, errStatusNotFound)
	}
	return respond(c, Response(getStatus(u.UID)))
}
