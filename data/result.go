package data

import "database/sql"

func InsertResult(tx *sql.Tx, job *Job, data string) error {
	_, err := tx.Exec("INSERT INTO results(data, jobID) VALUES(?, ?)", data, job.ID)
	if err != nil {
		return err
	}

	return nil
}
