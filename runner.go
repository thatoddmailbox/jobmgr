package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/thatoddmailbox/jobmgr/data"
)

const artifactsDirName = "artifacts"

func runJob(job *data.Job, parameters map[string]string) (string, string, error) {
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

	env := []string{}
	for _, p := range jobspec.Parameter {
		rawValue, exists := parameters[p.Name]
		if !exists {
			// TODO: default values?
			return "", "", fmt.Errorf("parameter '%s' is required but not set", p.Name)
		}

		value, ok := p.ParseValue(rawValue)
		if !ok {
			return "", "", fmt.Errorf("parameter '%s' has invalid value '%s'", p.Name, value)
		}

		filteredName := strings.ToUpper(
			strings.ReplaceAll(p.Name, " ", "_"),
		)

		envName := "JOBMGR_PARAMETER_" + filteredName

		env = append(env, envName+"="+value)
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

	env = append(env, "JOBMGR_ARTIFACTS_DIR="+artifactsDir)

	context, cancel := context.WithTimeout(context.Background(), jobspec.Timeout.Duration)
	defer cancel()

	cmd := exec.CommandContext(
		context,
		jobspec.Command,
		jobspec.Arguments...,
	)
	cmd.Dir = jobspec.WorkingDirectory
	cmd.Env = env

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
