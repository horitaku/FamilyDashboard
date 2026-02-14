package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/models"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router)
	return router
}

func performRequest(router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	router.ServeHTTP(rec, req)
	return rec
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, out any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
		t.Fatalf("decode json: %v", err)
	}
}

func TestGetStatus(t *testing.T) {
	router := setupTestRouter()
	rec := performRequest(router, http.MethodGet, "/api/status")

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload models.StatusResponse
	decodeJSON(t, rec, &payload)

	if !payload.OK {
		t.Fatalf("payload ok = false")
	}

	if _, err := time.Parse(time.RFC3339, payload.Now); err != nil {
		t.Fatalf("now parse error: %v", err)
	}

	if payload.LastUpdated.Weather == "" || payload.LastUpdated.Calendar == "" || payload.LastUpdated.Tasks == "" {
		t.Fatalf("lastUpdated has empty value")
	}
}

func TestGetCalendar(t *testing.T) {
	router := setupTestRouter()
	rec := performRequest(router, http.MethodGet, "/api/calendar")

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload models.CalendarResponse
	decodeJSON(t, rec, &payload)

	if len(payload.Days) == 0 {
		t.Fatalf("days is empty")
	}

	if _, err := time.Parse("2006-01-02", payload.Days[0].Date); err != nil {
		t.Fatalf("date parse error: %v", err)
	}
}

func TestGetTasks(t *testing.T) {
	router := setupTestRouter()
	rec := performRequest(router, http.MethodGet, "/api/tasks")

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload models.TasksResponse
	decodeJSON(t, rec, &payload)

	if len(payload.Items) == 0 {
		t.Fatalf("items is empty")
	}

	for _, item := range payload.Items {
		if item.CreatedAt.IsZero() {
			t.Fatalf("createdAt is zero")
		}
		if item.DueDate != nil {
			if _, err := time.Parse("2006-01-02", *item.DueDate); err != nil {
				t.Fatalf("dueDate parse error: %v", err)
			}
		}
	}
}

func TestGetWeather(t *testing.T) {
	router := setupTestRouter()
	rec := performRequest(router, http.MethodGet, "/api/weather")

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload models.WeatherResponse
	decodeJSON(t, rec, &payload)

	if payload.Location == "" {
		t.Fatalf("location is empty")
	}
	if payload.Current.Condition == "" {
		t.Fatalf("current.condition is empty")
	}
}

func TestHealth(t *testing.T) {
	router := setupTestRouter()
	rec := performRequest(router, http.MethodGet, "/api/health")

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload map[string]bool
	decodeJSON(t, rec, &payload)

	if ok := payload["ok"]; !ok {
		t.Fatalf("health ok = false")
	}
}
