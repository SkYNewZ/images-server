package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_server_handleHealthCheck(t *testing.T) {
	s := new(server)

	req, err := http.NewRequest("GET", "/_health", nil)
	if err != nil {
		t.Fatal(err)
	}

	buildNumber = "foo"

	rr := httptest.NewRecorder()
	err = s.handleHealthCheck(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ok":true,"version":"foo"}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
