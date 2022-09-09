package data

import (
	"context"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid"

	"github.com/thatoddmailbox/jobmgr/config"
)

type Artifact struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	MIME    string `json:"mime"`
	Size    int64  `json:"size"`
	UUID    string `json:"uuid"`
	Created int64  `json:"created"`
	JobID   int    `json:"jobID"`
}

var mimeMap = map[string]string{
	"pdf": "application/pdf",
	"zip": "application/zip",

	"mp3": "audio/mpeg",
	"ogg": "audio/ogg",
	"oga": "audio/ogg",
	"wav": "audio/wav",

	"gif":  "image/gif",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"webp": "image/webm",

	"csv": "text/csv",
	"txt": "text/plain",

	"mp4":  "video/mp4",
	"ogv":  "video/ogg",
	"webm": "video/webm",
}

func getKeyForArtifact(a *Artifact) string {
	return strconv.Itoa(a.JobID) + "/" + a.UUID
}

func GetArtifactByID(id int) (Artifact, error) {
	rows, err := DB.Query("SELECT id, name, mime, size, uuid, created, jobID FROM artifacts WHERE id = ?", id)
	if err != nil {
		return Artifact{}, err
	}

	defer rows.Close()
	if !rows.Next() {
		return Artifact{}, ErrNotFound
	}

	artifact := Artifact{}
	err = rows.Scan(&artifact.ID, &artifact.Name, &artifact.MIME, &artifact.Size, &artifact.UUID, &artifact.Created, &artifact.JobID)
	if err != nil {
		return Artifact{}, err
	}

	return artifact, nil
}

func GetArtifactsForJob(job *Job) ([]Artifact, error) {
	rows, err := DB.Query("SELECT id, name, mime, size, uuid, created, jobID FROM artifacts WHERE jobID = ?", job.ID)
	if err != nil {
		return []Artifact{}, err
	}

	defer rows.Close()

	result := []Artifact{}
	for rows.Next() {
		item := Artifact{}
		err = rows.Scan(&item.ID, &item.Name, &item.MIME, &item.Size, &item.UUID, &item.Created, &item.JobID)
		if err != nil {
			return []Artifact{}, err
		}

		result = append(result, item)
	}

	return result, nil
}

func UploadArtifact(name string, f fs.File, job *Job) error {
	a := Artifact{
		Name:    name,
		MIME:    "",
		Size:    0,
		UUID:    "",
		Created: time.Now().Unix(),
		JobID:   job.ID,
	}

	ext := strings.ToLower(filepath.Ext(name))[1:]
	mime, ok := mimeMap[ext]
	if !ok {
		mime = "application/octet-stream"
	}
	a.MIME = mime

	info, err := f.Stat()
	if err != nil {
		return err
	}
	a.Size = info.Size()

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	a.UUID = uuid.String()

	key := getKeyForArtifact(&a)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &config.Current.AWS.ArtifactsBucket,
		Key:    &key,
		Body:   f,
	})
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"INSERT INTO artifacts(name, mime, size, uuid, created, jobID) VALUES(?, ?, ?, ?, ?, ?)",
		a.Name,
		a.MIME,
		a.Size,
		a.UUID,
		a.Created,
		a.JobID,
	)
	if err != nil {
		return err
	}

	return nil
}
