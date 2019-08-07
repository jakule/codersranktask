package server

import (
	"fmt"
	"net/http"
)

func Index(c *CallParams, w http.ResponseWriter, _ *http.Request) {
	c.Infof("called index endpoint")
	_, _ = fmt.Fprintf(w, "Hello World!")
}
