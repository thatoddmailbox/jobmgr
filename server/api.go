package server

type statusResponse struct {
	Status string `json:"status"`
}

type createdResponse struct {
	Status  string `json:"status"`
	Created int64  `json:"created"`
}

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type errorHintResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Hint   string `json:"hint"`
}
