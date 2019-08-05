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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jakule/codersranktask/go/storage"
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

	var expireAfterTime *time.Time
	if secret.ExpireAfter > 0 {
		now := time.Now().UTC()
		t := now.Add(time.Duration(secret.ExpireAfter) * time.Minute)
		expireAfterTime = &t
	}

	secretData := &storage.SecretData{
		Secret:           secret.Secret,
		ExpireAfterViews: secret.ExpireAfterViews,
		ExpireAfterTime:  expireAfterTime,
	}
	id, err := c.Storage().CreateSecret(secretData)
	if err != nil {
		c.Errorf("insert failed : %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	c.Infof("created secret %s", id)

	secretModel := &Secret{
		Hash:           id,
		SecretText:     secret.Secret,
		CreatedAt:      time.Now().UTC(),
		ExpiresAt:      expireAfterTime,
		RemainingViews: int32(secret.ExpireAfterViews),
	}

	writeResponse(c, w, r, secretModel)
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

	if expireAfterViews <= 0 {
		return nil, errors.New("expireAfterViews has to be grater than 0")
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

func writeResponse(c *CallParams, w http.ResponseWriter, r *http.Request, secretModel *Secret) {
	var data []byte
	var contentType string
	var err error

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
}

func GetSecretByHash(c *CallParams, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]

	secret, err := c.Storage().GetSecret(hash)
	switch err {
	case storage.ErrHashNotfound:
		c.Infof("secret not found")
		http.Error(w, "secret not found", http.StatusBadRequest)
		return
	case nil:
		break
	default:
		c.Errorf("failed to get secret %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if secret.ExpireAfterViews <= 0 ||
		(secret.ExpireAfterTime != nil && secret.ExpireAfterTime.Before(time.Now())) {
		go func() {
			c.Infof("deleting secret %s", hash)
			err := c.Storage().Delete(hash)
			if err != nil {
				// log but do not fail
				c.Errorf("failed to delete secret %v", err)
			}
		}()
		http.Error(w, "secret not found", http.StatusBadRequest)
		return
	}

	secretModel := &Secret{
		Hash:           hash,
		SecretText:     secret.Secret,
		CreatedAt:      secret.CreatedTime,
		ExpiresAt:      secret.ExpireAfterTime,
		RemainingViews: int32(secret.ExpireAfterViews),
	}

	writeResponse(c, w, r, secretModel)
}
