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
	"github.com/jakule/codersranktask/internal/mocks"
	"github.com/jakule/codersranktask/internal/storage"
)

func TestGetSecretByHash(t *testing.T) {
	tests := []struct {
		name          string
		acceptHeader  string
		unmarshalFunc func(data []byte, v interface{}) error
	}{
		{"AcceptJson", "application/json", json.Unmarshal},
		{"AcceptXML", "application/xml", xml.Unmarshal},
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
				Times(1)
			storageMock.
				EXPECT().
				Delete("").
				Times(0)

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

func TestGetSecretAndDelete(t *testing.T) {
	const secret = "secretText"
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	tests := []struct {
		name       string
		secretData *storage.SecretData
	}{
		{
			"ExpireAfterViewsDelete",
			&storage.SecretData{
				Secret:           secret,
				ExpireAfterViews: 0,
				ExpireAfterTime:  nil,
				CreatedTime:      time.Now().Add(-5 * time.Minute),
			},
		},
		{
			"ExpireAfterTimeDelete",
			&storage.SecretData{
				Secret:           secret,
				ExpireAfterViews: 10,
				ExpireAfterTime:  &oneMinuteAgo,
				CreatedTime:      time.Now().Add(-5 * time.Minute),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const hashString = "c7cb197a-de61-4190-8735-17ac5a343826"
			url := fmt.Sprintf("/v1/secret/%s", hashString)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()
			storageMock := mocks.NewMockStorage(ctrl)
			storageMock.
				EXPECT().
				GetSecret(hashString).
				Return(tt.secretData, nil).
				Times(1)
			storageMock.
				EXPECT().
				Delete(hashString).
				Times(1)

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
			if err := json.Unmarshal(b, &secretResp); err != nil {
				t.Errorf("failed to unmarshal %v", err)
			}

			if secretResp.SecretText != secret {
				t.Errorf("expected %s, got %s", secret, secretResp.SecretText)
			}
		})
	}
}

func TestGetSecretNotFound(t *testing.T) {
	type storedData struct {
		secretData *storage.SecretData
		err        error
	}
	tests := []struct {
		name      string
		hash      string
		args      storedData
		getCalled int
	}{
		{
			"GetWrongUrl",
			"blah",
			storedData{nil, nil},
			0,
		},
		{
			"GetNotFound",
			"5ee201be-e30b-4184-ae38-156748926f1f",
			storedData{nil, storage.ErrHashNotfound},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/v1/secret/%s", tt.hash)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()
			storageMock := mocks.NewMockStorage(ctrl)
			storageMock.
				EXPECT().
				GetSecret(tt.hash).
				Return(tt.args.secretData, tt.args.err).
				Times(tt.getCalled)
			storageMock.
				EXPECT().
				Delete("").
				Times(0)

			handler := NewRouter(createMockCallParams(storageMock))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("server returned wrong status code: got %v want %v",
					status, http.StatusNotFound)
			}
		})
	}
}
