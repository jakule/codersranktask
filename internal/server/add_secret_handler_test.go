package server

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
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

			req, err := http.NewRequest("POST", "/v1/secret", body)
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

			handler := NewRouter(createMockCallParams(storageMock))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("server returned wrong status code: got %v want %v",
					status, tt.statusCode)
			}
		})
	}
}

func TestAddSecretMissingForm(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/secret", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	storageMock := mocks.NewMockStorage(ctrl)
	handler := NewRouter(createMockCallParams(storageMock))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("server returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func createMockCallParams(storage storage.Storage) *CallParams {
	return NewCallParams(zap.NewNop(), storage)
}
