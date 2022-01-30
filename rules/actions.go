package rules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type FileInfo struct {
	Name         string
	Path         string
	LastModified time.Time
}

type RuleAction interface {
	Run(FileInfo) error
}

// MoveToDateFolder - basic rule to move a file to an archive folder and categorize by date created
type MoveToDateFolder struct {
	Path string
}

type Delete struct{}

func folderNameForTime(time time.Time) string {
	return filepath.Join(fmt.Sprintf("%d", time.Year()), fmt.Sprintf("%02d", time.Month()))
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (folder MoveToDateFolder) Run(file FileInfo) error {
	if exists, err := FileExists(file.Path); (err != nil) || !exists {
		log.WithFields(log.Fields{"name": file.Name}).Warn("file did not exist")
		return nil
	}

	destFile := filepath.Join(folder.Path, folderNameForTime(file.LastModified), file.Name)

	_ = os.MkdirAll(filepath.Dir(destFile), os.ModePerm)

	log.WithFields(log.Fields{"name": file.Name, "originalPath": file.Path, "destFile": destFile}).Info("MoveToDateFolder - moving file")

	return os.Rename(file.Path, destFile)
}

func (Delete) Run(file FileInfo) error {
	log.WithFields(log.Fields{"name": file.Name, "originalPath": file.Path}).Info("Delete - deleting file")

	return os.Remove(file.Path)
}
