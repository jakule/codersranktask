package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jakule/codersranktask/internal/mocks"
)

func TestIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/", nil)
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

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("server returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
