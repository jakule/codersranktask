package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jakule/codersranktask/internal/storage"
	"github.com/pkg/errors"
)

type addSecretRequest struct {
	secret           string
	expireAfterViews int
	expireAfter      int
}

func AddSecret(c *CallParams, w http.ResponseWriter, r *http.Request) {
	secret, err := validateAddSecret(r)
	if err != nil {
		c.Errorf("validation failed: %v", err)
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	var expireAfterTime *time.Time
	if secret.expireAfter > 0 {
		now := time.Now().UTC()
		t := now.Add(time.Duration(secret.expireAfter) * time.Minute)
		expireAfterTime = &t
	}

	secretData := &storage.SecretData{
		Secret:           secret.secret,
		ExpireAfterViews: secret.expireAfterViews,
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
		SecretText:     secret.secret,
		CreatedAt:      time.Now().UTC(),
		ExpiresAt:      expireAfterTime,
		RemainingViews: int32(secret.expireAfterViews),
	}

	writeResponse(c, w, r, secretModel)
}

func validateAddSecret(r *http.Request) (*addSecretRequest, error) {
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

	return &addSecretRequest{
		secret:           secret,
		expireAfterViews: expireAfterViews,
		expireAfter:      expireAfter,
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
