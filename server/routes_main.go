package server

func routeMain(c *requestContext) {
	c.w.Header().Set("Content-Type", "text/plain")
	c.w.Write([]byte("Beep boop blop I am a jobmgr API server"))
}
