package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thatoddmailbox/jobmgr/config"
	"github.com/thatoddmailbox/jobmgr/data"
	"github.com/thatoddmailbox/jobmgr/server"
)

func main() {
	log.Println("jobmgr")

	err := config.Load()
	if err != nil {
		panic(err)
	}

	err = data.Init()
	if err != nil {
		panic(err)
	}

	// TODO: theoretically in the future we'd want separate workers and servers
	// but this is okay for now i guess
	go server.StartServer()

	for {
		data.JobQueueNotification.L.Lock()

		job, err := data.GetNextJobInQueue()
		if err != nil {
			panic(err)
		}

		if job == nil {
			// there's no job available
			// wait for us to get notified about something being queued
			data.JobQueueNotification.Wait()
			data.JobQueueNotification.L.Unlock()

			// note that we're doing this in a loop
			// it's possible for Wait() to return without the condition being true
			// in that case, we'll just check the db, see that there's nothing, and go back to sleep
			continue
		}

		data.JobQueueNotification.L.Unlock()

		err = data.MarkJobStarted(job)
		if err != nil {
			panic(err)
		}

		result, tempDir, err := runJob(job, job.Parameters)
		if err != nil {
			// job failed!
			tx, err2 := data.DB.Begin()
			if err2 != nil {
				panic(err2)
			}

			err2 = data.InsertResult(tx, job, err.Error()+"\n\nJob directory: "+tempDir)
			if err2 != nil {
				tx.Rollback()
				panic(err2)
			}

			err2 = data.MarkJobFailed(tx, job)
			if err2 != nil {
				tx.Rollback()
				panic(err2)
			}

			err2 = tx.Commit()
			if err2 != nil {
				tx.Rollback()
				panic(err2)
			}

			continue
		}

		// TODO: save artifacts
		artifactsDir := filepath.Join(tempDir, artifactsDirName)
		err = filepath.WalkDir(artifactsDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// bubble up any errors
				return err
			}

			if d.IsDir() {
				return nil
			}

			relativePath := strings.Replace(path, artifactsDir, "", 1)[1:]

			f, err := os.Open(path)
			if err != nil {
				// bubble up any errors
				// TODO: should we handle this better?
				return err
			}
			defer f.Close()

			err = data.UploadArtifact(relativePath, f, job)
			if err != nil {
				// bubble up any errors
				// TODO: should we handle this better?
				return err
			}

			return nil
		})
		if err != nil {
			panic(err)
		}

		tx, err := data.DB.Begin()
		if err != nil {
			panic(err)
		}

		err = data.InsertResult(tx, job, result)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = data.MarkJobCompleted(tx, job)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		// clean up the temp directory
		err = os.RemoveAll(tempDir)
		if err != nil {
			panic(err)
		}
	}
}
