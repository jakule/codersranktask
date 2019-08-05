package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	swagger "github.com/jakule/codersranktask/internal"
	"github.com/jakule/codersranktask/internal/storage"
)

func GetSecretByHash(c *swagger.CallParams, w http.ResponseWriter, r *http.Request) {
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

	if secret.ExpireAfterViews < 0 ||
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
