package main

import (
	"os"
	"path/filepath"

	"github.com/thatoddmailbox/jobmgr/data"
)

const artifactsDirName = "artifacts"

func runJob(job *data.Job) (string, string, error) {
	wd, err := os.MkdirTemp("", "jobmgr-")
	if err != nil {
		return "", "", err
	}

	artifactsDir := filepath.Join(wd, artifactsDirName)

	err = os.Mkdir(artifactsDir, 0777)
	if err != nil {
		return "", "", err
	}

	// TODO: run the job here :)s

	return ":)", wd, nil
}
