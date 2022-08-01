package app

import (
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

var (
	basiliPath = "basili"
	casperPath = "casper"
	zeusPath   = "zeus"
	picPath    = "pic"
)

// handleBasili sends a photo of the Basili's cat.
func (a *App) handleBasili(c tele.Context) error {
	return sendRandomPicture(c, a.Locate(basiliPath))
}

// handleBasili sends a photo of the Leonid's cat.
func (a *App) handleCasper(c tele.Context) error {
	return sendRandomPicture(c, a.Locate(casperPath))
}

// handleZeus sends a photo of the Solar's cat.
func (a *App) handleZeus(c tele.Context) error {
	return sendRandomPicture(c, a.Locate(zeusPath))
}

// handlePic sends a photo from a hierarchy of directories located at picPath.
func (a *App) handlePic(c tele.Context) error {
	return sendRandomPictureWith(c, a.Locate(picPath), randomFileFromHierarchy)
}

// sendRandomPicture sends a random picture from the directory.
func sendRandomPicture(c tele.Context, dir string) error {
	return sendRandomPictureWith(c, dir, randomFile)
}

// sendRandomPictureWith sends a random picture chosen by f from the directory.
func sendRandomPictureWith(c tele.Context, dir string, f randomFileFunc) error {
	path, err := f(dir)
	if err != nil {
		return respondInternalError(c, err)
	}
	return sendPicture(c, path)
}

// sendPicture sends the file located at the path.
func sendPicture(c tele.Context, path string) error {
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

type randomFileFunc func(string) (string, error)

// randomFile returns a random filename from a directory.
func randomFile(dir string) (string, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	f := fs[rand.Intn(len(fs))]
	return filepath.Join(dir, f.Name()), nil
}

// randomFileFromHierarchy returns a random filename from a hierarchy of directories.
func randomFileFromHierarchy(root string) (string, error) {
	var filenames []string
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			filenames = append(filenames, path)
		}
		return nil
	}); err != nil {
		return "", err
	}
	return filenames[rand.Intn(len(filenames))], nil
}
