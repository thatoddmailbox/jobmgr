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
	ID          int               `json:"id"`
	Status      JobStatus         `json:"-"`
	StatusText_ string            `json:"status"`
	Priority    int               `json:"priority"`
	Name        string            `json:"name"`
	Parameters  map[string]string `json:"parameters"`
	Created     int               `json:"created"`
	Started     *int              `json:"started"`
	Completed   *int              `json:"completed"`
	UserID      int               `json:"userID"`
}

func hydrateJob(j Job) Job {
	j.StatusText_ = map[JobStatus]string{
		JobStatusQueued:    "queued",
		JobStatusStarted:   "started",
		JobStatusCompleted: "completed",
		JobStatusFailed:    "failed",
	}[j.Status]
	return j
}

func EnqueueJob(name string, parameters map[string]string, priority int, userID int) (int64, error) {
	parameterBytes, err := json.Marshal(parameters)
	if err != nil {
		return 0, err
	}

	result, err := DB.Exec(
		"INSERT INTO jobs(status, priority, name, parameters, created, userID) VALUES(?, ?, ?, ?, ?, ?)",
		JobStatusQueued,
		priority,
		name,
		parameterBytes,
		time.Now().Unix(),
		userID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetJobByID(id int) (Job, error) {
	rows, err := DB.Query("SELECT id, status, priority, name, parameters, created, started, completed, userID FROM jobs WHERE id = ?", id)
	if err != nil {
		return Job{}, err
	}

	defer rows.Close()
	if !rows.Next() {
		return Job{}, ErrNotFound
	}

	job := Job{}
	parameterString := ""
	err = rows.Scan(&job.ID, &job.Status, &job.Priority, &job.Name, &parameterString, &job.Created, &job.Started, &job.Completed, &job.UserID)
	if err != nil {
		return Job{}, err
	}

	err = json.Unmarshal([]byte(parameterString), &job.Parameters)
	if err != nil {
		return Job{}, err
	}

	return hydrateJob(job), nil
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

	job = hydrateJob(job)

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
