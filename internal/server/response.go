package server

import "github.com/gin-gonic/gin"

// apiError is the structured error payload: a human-readable message and a
// machine-readable numeric code for clients to branch on.
type apiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// meta carries response-level metadata. error is always present: nil on
// success, an apiError on failure. Other fields (e.g. pagination) can be added
// here later.
type meta struct {
	Error *apiError `json:"error"`
}

// envelope is the standard shape for every API response: a meta object and a
// data payload. On success data holds the payload and meta.error is null; on
// failure data is null and meta.error carries the message.
type envelope struct {
	Meta meta `json:"meta"`
	Data any  `json:"data"`
}

// respond writes data wrapped in the standard envelope with a null error.
func respond(c *gin.Context, status int, data any) {
	c.JSON(status, envelope{Meta: meta{Error: nil}, Data: data})
}

// respondError writes an error in the same envelope: meta.error carries the
// message and code, and data is null. status sets the HTTP status; code is the
// machine-readable identifier in the body.
func respondError(c *gin.Context, status, code int, msg string) {
	c.JSON(status, envelope{Meta: meta{Error: &apiError{Message: msg, Code: code}}, Data: nil})
}
