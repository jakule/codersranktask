package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jakule/codersranktask/internal/mocks"
	"github.com/jakule/codersranktask/internal/storage"
	"go.uber.org/zap"
)

func TestAddSecret(t *testing.T) {
	tests := []struct {
		name       string
		fields     map[string]string
		statusCode int
	}{
		{
			"OK",
			map[string]string{
				"secret":           "secretString",
				"expireAfterViews": "15",
				"expireAfter":      "0",
			},
			http.StatusOK,
		},
		{
			"ExpireAfterViewsZero",
			map[string]string{
				"secret":           "secretString",
				"expireAfterViews": "0",
				"expireAfter":      "60",
			},
			http.StatusBadRequest,
		},
		{
			"MissingSecret",
			map[string]string{
				"expireAfterViews": "15",
				"expireAfter":      "60",
			},
			http.StatusBadRequest,
		},
		{
			"MissingExpireAfter",
			map[string]string{
				"secret":           "secretString",
				"expireAfterViews": "15",
			},
			http.StatusBadRequest,
		},
		{
			"MissingExpireAfterViews",
			map[string]string{
				"secret":      "secretString",
				"expireAfter": "60",
			},
			http.StatusBadRequest,
		},
		{
			"WrongExpireAfter",
			map[string]string{
				"secret":           "secretString",
				"expireAfterViews": "21fse",
			},
			http.StatusBadRequest,
		},
		{
			"WrongExpireAfterViews",
			map[string]string{
				"secret":      "secretString",
				"expireAfter": "vd",
			},
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			for k, v := range tt.fields {
				_ = writer.WriteField(k, v)
			}
			_ = writer.Close()

			req, err := http.NewRequest("POST", "/secret", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()
			storageMock := mocks.NewMockStorage(ctrl)
			storageMock.EXPECT().CreateSecret(&storage.SecretData{
				Secret: tt.fields["secret"],
				//Secret:           "secret",
				ExpireAfterViews: 15,
				ExpireAfterTime:  nil,
			}).AnyTimes()

			handler := handlerWrapper(createMockCallParams(storageMock), AddSecret)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("server returned wrong status code: got %v want %v",
					status, tt.statusCode)
			}
		})
	}
}

func TestAddSecretMissingForm(t *testing.T) {
	req, err := http.NewRequest("POST", "/secret", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	storageMock := mocks.NewMockStorage(ctrl)
	handler := handlerWrapper(createMockCallParams(storageMock), AddSecret)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("server returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetSecretByHash(t *testing.T) {
	req, err := http.NewRequest("GET", "/secret", nil)
	if err != nil {
		t.Fatal(err)
	}
	hashString := "c7cb197a-de61-4190-8735-17ac5a343826"
	secret := "secretText"
	req = mux.SetURLVars(req, map[string]string{"hash": hashString})
	req.Header.Set("Accept", "application/json")

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

	handler := handlerWrapper(createMockCallParams(storageMock), GetSecretByHash)
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
}

func handlerWrapper(params *CallParams, inner paramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner(params, w, r)
	})
}

func createMockCallParams(storage storage.Storage) *CallParams {
	return NewCallParams(zap.NewNop(), storage)
}
