package main

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thatoddmailbox/jobmgr/data"
)

const artifactsDirName = "artifacts"

func runJob(job *data.Job) (string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	if strings.ContainsAny(job.Name, "./\\ ") {
		return "", "", errors.New("invalid job name")
	}

	jobspecPath := filepath.Join(wd, "jobspecs", job.Name+".toml")
	jobspec, err := data.ParseJobSpec(jobspecPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", "", errors.New("job does not exist")
		}

		return "", "", err
	}

	tempDir, err := os.MkdirTemp("", "jobmgr-")
	if err != nil {
		return "", tempDir, err
	}

	artifactsDir := filepath.Join(tempDir, artifactsDirName)

	err = os.Mkdir(artifactsDir, 0777)
	if err != nil {
		return "", tempDir, err
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		context,
		jobspec.Command,
		jobspec.Arguments...,
	)
	cmd.Dir = jobspec.WorkingDirectory

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", tempDir, err
	}
	defer stdoutPipe.Close()

	err = cmd.Start()
	if err != nil {
		return "", tempDir, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return "", tempDir, err
	}

	err = cmd.Wait()
	if err != nil {
		return "", tempDir, err
	}

	return string(stdout), tempDir, nil
}
