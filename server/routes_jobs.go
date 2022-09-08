package server

import (
	"strconv"

	"github.com/thatoddmailbox/jobmgr/data"
)

type jobResponse struct {
	Status    string          `json:"status"`
	Job       data.Job        `json:"job"`
	Result    *string         `json:"result"`
	Artifacts []data.Artifact `json:"artifacts"`
}

func routeJobsGet(c *requestContext) {
	idString := c.r.FormValue("id")
	if idString == "" {
		c.WriteJSON(errorResponse{"error", "missing_params"})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.WriteJSON(errorResponse{"error", "invalid_params"})
		return
	}

	job, err := data.GetJobByID(id)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	var result *string

	if job.Status == data.JobStatusCompleted || job.Status == data.JobStatusFailed {
		resultString, err := data.GetResultForJob(&job)
		if err != nil {
			c.InternalServerError(err)
			return
		}

		result = &resultString
	}

	var artifacts []data.Artifact

	if job.Status == data.JobStatusCompleted {
		artifacts, err = data.GetArtifactsForJob(&job)
		if err != nil {
			c.InternalServerError(err)
			return
		}
	}

	c.WriteJSON(jobResponse{"ok", job, result, artifacts})
}
