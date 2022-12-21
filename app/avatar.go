package app

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

const (
	avatarDir       = "avatar"
	avatarMaxHeight = 1500
	avatarMaxWidth  = 1500
	avatarMaxSize   = UserError("Максимальный размер аватара - %dx%d пикселей.")
	avatarSet       = Response("✅ Аватар установлен.")
	avatarAttach    = UserError("Прикрепите изображение.")
)

func (a *App) handleAvatar(c tele.Context) error {
	user := getUser(c)
	if c.Message().Photo != nil {
		pic := c.Message().Photo
		if pic.Width > avatarMaxWidth || pic.Height > avatarMaxHeight {
			return respondUserError(c, avatarMaxSize.Fill(avatarMaxHeight, avatarMaxWidth))
		}
		if err := a.setAvatar(user.UID, pic.File); err != nil {
			return respondInternalError(c, err)
		}
		return respond(c, avatarSet)
	}
	if ava, ok := loadAvatar(user.UID); ok {
		return c.Send(ava)
	}
	return respondUserError(c, avatarAttach)
}

func (a *App) setAvatar(uid int64, f tele.File) error {
	rc, err := a.bot.File(&f)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	rc.Close()
	if err := os.MkdirAll(avatarDir, 0777); err != nil {
		return err
	}
	return os.WriteFile(avatarPath(uid), data, 0666)
}

func avatarPath(uid int64) string {
	return filepath.Join(avatarDir, strconv.FormatInt(uid, 10))
}

func loadAvatar(uid int64) (avatar *tele.Photo, ok bool) {
	ava := tele.FromDisk(avatarPath(uid))
	if ava.OnDisk() {
		return &tele.Photo{File: ava}, true
	}
	return nil, false
}
