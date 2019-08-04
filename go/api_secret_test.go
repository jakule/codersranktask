package swagger

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
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
				"expireAfter":      "60",
			},
			http.StatusOK,
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
			handler := handlerWrapperLogger(AddSecret)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
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
	handler := handlerWrapperLogger(AddSecret)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
