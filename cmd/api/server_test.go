package main

import (
	"bytes"
	"github.com/omid-h70/lucky-wheel-game/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	TEST_SERVER_CONFIG = "0.0.0.0:8000"
	TRANSFER_URL       = "/v1/health"
)

func Test_transfer_should_return_fail_when_server_is_not_ready_or_down(t *testing.T) {
	var jsonData = []byte(`{
		"UUID" : "f1bc8f04-0500-11ee-be56-0242ac120002"
	}`)

	request, err := http.NewRequest(http.MethodPost, TRANSFER_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	app := NewApp(db.NewMockDBHandler())

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(app.testYourLuck)
	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if recorder.Code != http.StatusOK {
		t.Error("Test Failed")
	}
}
