package server

import (
	"fmt"
	"net/http"

	swagger "github.com/jakule/codersranktask/internal"
)

func Index(c *swagger.CallParams, w http.ResponseWriter, _ *http.Request) {
	c.Infof("called index endpoint")
	_, _ = fmt.Fprintf(w, "Hello World!")
}
