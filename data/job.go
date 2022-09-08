package data

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JobStatus int

const (
	JobStatusQueued    JobStatus = 0
	JobStatusStarted   JobStatus = 1
	JobStatusCompleted JobStatus = 2
	JobStatusFailed    JobStatus = 3
)

type Job struct {
	ID         string
	Status     JobStatus
	Priority   int
	Name       string
	Parameters map[string]string
	Created    int
	Started    *int
	Completed  *int
	UserID     int
}

func GetNextJobInQueue() (*Job, error) {
	rows, err := DB.Query("SELECT id, status, priority, name, parameters, created, started, completed, userID FROM jobs WHERE status = 0 ORDER BY priority ASC, created ASC LIMIT 1")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	job := Job{}
	parameterString := ""
	err = rows.Scan(&job.ID, &job.Status, &job.Priority, &job.Name, &parameterString, &job.Created, &job.Started, &job.Completed, &job.UserID)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(parameterString), &job.Parameters)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func MarkJobCompleted(tx *sql.Tx, job *Job) error {
	_, err := tx.Exec("UPDATE jobs SET status = ?, completed = ? WHERE id = ?", JobStatusCompleted, time.Now().Unix(), job.ID)
	if err != nil {
		return err
	}

	return nil
}

func MarkJobFailed(tx *sql.Tx, job *Job) error {
	_, err := tx.Exec("UPDATE jobs SET status = ?, completed = ? WHERE id = ?", JobStatusFailed, time.Now().Unix(), job.ID)
	if err != nil {
		return err
	}

	return nil
}

func MarkJobStarted(job *Job) error {
	_, err := DB.Exec("UPDATE jobs SET status = ?, started = ? WHERE id = ?", JobStatusStarted, time.Now().Unix(), job.ID)
	if err != nil {
		return err
	}

	return nil
}
