/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AddSecretRequest struct {
	Secret           string
	ExpireAfterViews int
	ExpireAfter      int
}

func AddSecret(c *CallParams, w http.ResponseWriter, r *http.Request) {
	secret, err := validateAddSecret(r)
	if err != nil {
		c.Errorf("validation failed: %v", err)
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	id, err := c.Storage().CreateSecret(secret.Secret)
	if err != nil {
		c.Errorf("insert failed : %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	secretModel := &Secret{
		Hash:           id,
		SecretText:     secret.Secret,
		CreatedAt:      time.Now().UTC(),
		ExpiresAt:      time.Time{},
		RemainingViews: int32(secret.ExpireAfterViews),
	}

	var data []byte
	var contentType string

	acceptHeader := r.Header.Get("Accept")
	switch acceptHeader {
	case "application/json":
		contentType = "application/json; charset=UTF-8"
		data, err = json.Marshal(secretModel)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	case "application/xml":
		contentType = "application/xml; charset=UTF-8"
		data, err = xml.Marshal(secretModel)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	default:
		c.Errorf("unsupported response type : %s", acceptHeader)
		http.Error(w, "unsupported response type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
	return
}

func validateAddSecret(r *http.Request) (*AddSecretRequest, error) {
	err := r.ParseMultipartForm(1024 * 1024)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form : %v", err)
	}

	secret := r.FormValue("secret")
	if secret == "" {
		return nil, fmt.Errorf("missing secret field")
	}

	expireAfterViews, err := parseFormInt(r, "expireAfterViews")
	if err != nil {
		return nil, err
	}

	expireAfter, err := parseFormInt(r, "expireAfter")
	if err != nil {
		return nil, fmt.Errorf("expireAfterViews is not an integer")
	}

	return &AddSecretRequest{
		Secret:           secret,
		ExpireAfterViews: expireAfterViews,
		ExpireAfter:      expireAfter,
	}, nil
}

func parseFormInt(r *http.Request, fieldName string) (int, error) {
	s := r.FormValue(fieldName)
	if s == "" {
		return 0, fmt.Errorf("expireAfterViews is missing")
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("expireAfterViews is not an integer")
	}

	return val, nil
}

func GetSecretByHash(c *CallParams, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
