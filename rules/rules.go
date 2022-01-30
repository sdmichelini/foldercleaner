package rules

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RuleMatcher interface {
	Matches(string) bool
}

type AllFiles struct{}

func (AllFiles) Matches(string) bool {
	return true
}

type FileExtensions struct {
	Extensions []string
}

func (extensions FileExtensions) Matches(filepath string) bool {
	for _, extension := range extensions.Extensions {
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(extension)) {
			return true
		}
	}
	return false
}

type FileLocation string

func FilesForLocationOlderThan(location FileLocation, olderThanCutoff time.Time) ([]FileInfo, error) {
	files, err := os.ReadDir(string(location))

	if err != nil {
		return nil, err
	}

	ret := make([]FileInfo, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := filepath.Join(string(location), file.Name())

		fileInfo, err := os.Stat(filename)

		if err != nil {
			log.Warnf("could not stat file: %s: %v", filename, err)
			continue
		}

		// ignore files after the cutoff
		if fileInfo.ModTime().After(olderThanCutoff) {
			continue
		}

		ret = append(ret, FileInfo{
			Name:         file.Name(),
			Path:         filename,
			LastModified: fileInfo.ModTime(),
		})
	}

	return ret, nil
}

type FileRule struct {
	RuleAction
	RuleMatcher
	FileLocation
	// rule only applies to files older than this
	Cutoff time.Duration
}
