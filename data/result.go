package data

import "database/sql"

func GetResultForJob(job *Job) (string, error) {
	rows, err := DB.Query("SELECT data FROM results WHERE jobID = ? LIMIT 1", job.ID)
	if err != nil {
		return "", err
	}

	defer rows.Close()
	if !rows.Next() {
		return "", ErrNotFound
	}

	data := ""
	err = rows.Scan(&data)
	if err != nil {
		return "", err
	}

	return data, nil
}

func InsertResult(tx *sql.Tx, job *Job, data string) error {
	_, err := tx.Exec("INSERT INTO results(data, jobID) VALUES(?, ?)", data, job.ID)
	if err != nil {
		return err
	}

	return nil
}
