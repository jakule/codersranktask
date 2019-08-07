package server

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/jakule/codersranktask/internal/storage"
)

// uuidRegexFormat is a regex to validate UUIDv4 string
const uuidRegexFormat = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

var uuidRegex = regexp.MustCompile(uuidRegexFormat)

func GetSecretByHash(c *CallParams, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]

	// Check if hash comes as UUID string
	if !uuidRegex.MatchString(hash) {
		c.Infof("secret not found")
		http.Error(w, "secret not found", http.StatusBadRequest)
		return
	}

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

	if secretExpired(secret) {
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

func secretExpired(secret *storage.SecretData) bool {
	return secret.ExpireAfterViews < 0 ||
		(secret.ExpireAfterTime != nil && secret.ExpireAfterTime.Before(time.Now()))
}
