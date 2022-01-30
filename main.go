package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"time"

	"foldercleaner/rules"
)

func main() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("couldn't get homedir: %v", err)
	}

	fileRules := []rules.FileRule{
		{
			RuleAction:   rules.Delete{},
			RuleMatcher:  rules.FileExtensions{Extensions: []string{".exe", ".msi"}},
			FileLocation: rules.FileLocation(path.Join(homeDir, "Downloads")),
			Cutoff:       time.Hour * 24 * 7,
		},
		{
			RuleAction: rules.MoveToDateFolder{
				Path: filepath.Join(homeDir, "archive"),
			},
			RuleMatcher:  rules.AllFiles{},
			FileLocation: rules.FileLocation(path.Join(homeDir, "Downloads")),
			Cutoff:       time.Hour * 24 * 7,
		},
	}

	for _, fileRule := range fileRules {
		files, err := rules.FilesForLocationOlderThan(fileRule.FileLocation, time.Now().Add(-fileRule.Cutoff))
		filteredFiles := make([]rules.FileInfo, 0, len(files))

		fmt.Println(len(filteredFiles))

		if err != nil {
			log.Fatalf("could not run rule: %v", fileRule)
		}

		for _, file := range files {
			if fileRule.RuleMatcher.Matches(file.Name) {
				filteredFiles = append(filteredFiles, file)
			}
		}

		for _, file := range filteredFiles {
			err := fileRule.RuleAction.Run(file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Info("DONE")
}
