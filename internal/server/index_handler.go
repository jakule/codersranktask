package server

import (
	"fmt"
	"net/http"

	"github.com/jakule/codersranktask/internal"
)

func Index(c *internal.CallParams, w http.ResponseWriter, _ *http.Request) {
	c.Infof("called index endpoint")
	_, _ = fmt.Fprintf(w, "Hello World!")
}
