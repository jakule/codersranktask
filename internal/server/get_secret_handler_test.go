package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jakule/codersranktask/internal/mocks"
	"github.com/jakule/codersranktask/internal/storage"
)

func TestGetSecretByHash(t *testing.T) {
	tests := []struct {
		name          string
		acceptHeader  string
		unmarshalFunc func(data []byte, v interface{}) error
	}{
		{"json", "application/json", json.Unmarshal},
		{"json", "application/xml", xml.Unmarshal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const hashString = "c7cb197a-de61-4190-8735-17ac5a343826"
			url := fmt.Sprintf("/v1/secret/%s", hashString)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}
			secret := "secretText"
			req = mux.SetURLVars(req, map[string]string{"hash": hashString})
			req.Header.Set("Accept", tt.acceptHeader)

			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()
			storageMock := mocks.NewMockStorage(ctrl)
			storageMock.
				EXPECT().
				GetSecret(hashString).
				Return(&storage.SecretData{
					Secret:           secret,
					ExpireAfterViews: 5,
					ExpireAfterTime:  nil,
					CreatedTime:      time.Now().Add(-5 * time.Minute),
				}, nil).
				AnyTimes()

			handler := NewRouter(createMockCallParams(storageMock))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("server returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			b, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Errorf("failed to parse body %v", err)
			}
			var secretResp Secret
			if err := tt.unmarshalFunc(b, &secretResp); err != nil {
				t.Errorf("failed to unmarshal %v", err)
			}

			if secretResp.SecretText != secret {
				t.Errorf("expected %s, got %s", secret, secretResp.SecretText)
			}

		})
	}
}
