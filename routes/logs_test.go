package routes

import (
	"net/http"
	"net/http/httptest"
	mock_services "static-site-hosting/services/mocks"
	"testing"
)

func TestLogsHandler(t *testing.T) {
	t.Run("Returns DeploymentLogs", func(t *testing.T) {
		// Create request to test endpoint
		req, err := http.NewRequest("GET", "/logs/test-site", nil)
		if err != nil {
			t.Fatal(err)
		}

		fakeDB := mock_services.NewFakeDB()
		// Create response recorder
		rr := httptest.NewRecorder()
		handler := DeploymentLogsHandler(fakeDB)

		// Execute handler
		handler.ServeHTTP(rr, req)

		// Verify status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Verify response body
		expected := "[{\"site_name\":\"\",\"timestamp\":\"2023-10-01T12:00:00Z\",\"ip_address\":\"127.0.0.1\",\"user_agent\":\"TestAgent\"}]\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %q want %q",
				rr.Body.String(), expected)
		}
	})

	t.Run("Returns method not allowed for non-GET requests", func(t *testing.T) {
		// Create request to test endpoint
		req, err := http.NewRequest("PUT", "/logs/non-existent-site", nil)
		if err != nil {
			t.Fatal(err)
		}

		fakeDB := mock_services.NewFakeDB()
		// Create response recorder
		rr := httptest.NewRecorder()
		handler := DeploymentLogsHandler(fakeDB)

		// Execute handler
		handler.ServeHTTP(rr, req)

		// Verify status code
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})
}
