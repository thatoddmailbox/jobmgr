package server

import (
	"net/http"
	"strconv"

	"github.com/thatoddmailbox/jobmgr/data"
)

func routeArtifactsDownload(c *requestContext) {
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

	artifact, err := data.GetArtifactByID(id)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	job, err := data.GetJobByID(artifact.JobID)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	// TODO: check job user id
	_ = job

	url, err := data.BuildURLForArtifact(&artifact)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	http.Redirect(c.w, c.r, url, http.StatusTemporaryRedirect)
}
