package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	"github.com/rihow/FamilyDashboard/internal/models"
	"github.com/rihow/FamilyDashboard/internal/services/nextcloud"
	"github.com/rihow/FamilyDashboard/internal/services/weather"
	"github.com/rihow/FamilyDashboard/internal/status"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Config{
		RefreshIntervals: config.RefreshIntervals{
			WeatherSec:  300,
			CalendarSec: 300,
			TasksSec:    300,
		},
		Location: config.Location{
			CityName: "姫路市",
			Country:  "JP",
		},
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			CalendarNames: []string{"family"},
			TaskListNames: []string{"tasks"},
		},
	}

	fc := cache.New(t.TempDir())
	seedCache(t, fc, cfg)

	weatherClient := weather.NewClient(fc, "http://localhost:8080")
	nextcloudClient, _ := nextcloud.NewClient(fc, cfg)
	errorStore := status.NewErrorStore()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Set("cache", fc)
		ctx.Set("weather", weatherClient)
		ctx.Set("nextcloud", nextcloudClient)
		ctx.Set("errorStore", errorStore)
		ctx.Next()
	})

	SetupRoutes(router)
	return router
}

func seedCache(t *testing.T, fc *cache.FileCache, cfg *config.Config) {
	t.Helper()

	weatherKey := fmt.Sprintf("weather:%s:%s", cfg.Location.Country, cfg.Location.CityName)
	weatherPayload := &models.WeatherResponse{
		Location: cfg.Location.CityName,
		Current: models.CurrentWeather{
			Temperature: 10,
			Condition:   "晴",
			Icon:        "01d",
			Humidity:    40,
			WindSpeed:   2,
		},
		Today: models.TodayWeather{
			MaxTemp: 15,
			MinTemp: 5,
			Summary: "晴",
		},
		PrecipSlots: []models.PrecipSlot{{Time: "09:00", Precip: 10}},
		Alerts:      []models.WeatherAlert{},
	}
	if _, err := fc.Write(weatherKey, weatherPayload, map[string]string{"source": "test"}); err != nil {
		t.Fatalf("seed weather cache: %v", err)
	}

	calendarPayload := &models.CalendarResponse{
		Days: []models.CalendarDay{
			{
				Date:   time.Now().Format("2006-01-02"),
				AllDay: []models.Event{},
				Timed: []models.Event{
					{
						ID:       "seed-event",
						Title:    "テスト",
						Start:    time.Now().Format(time.RFC3339),
						End:      time.Now().Add(1 * time.Hour).Format(time.RFC3339),
						Color:    "#A4BDFC",
						Calendar: "shared",
					},
				},
			},
		},
	}
	if _, err := fc.Write("nextcloud_calendar_events_all", calendarPayload, map[string]string{"source": "test"}); err != nil {
		t.Fatalf("seed calendar cache: %v", err)
	}

	tasksPayload := &models.TasksResponse{
		Items: []models.TaskItem{
			{
				ID:        "seed-task",
				Title:     "テスト",
				Notes:     "",
				Status:    "needsAction",
				DueDate:   stringPtr(time.Now().Format("2006-01-02")),
				Priority:  1,
				CreatedAt: time.Now().Add(-1 * time.Hour),
			},
		},
	}
	if _, err := fc.Write("nextcloud_tasks_items_all", tasksPayload, map[string]string{"source": "test"}); err != nil {
		t.Fatalf("seed tasks cache: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
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
	router := setupTestRouter(t)
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

	if _, err := time.Parse(time.RFC3339, payload.LastUpdated.Weather); err != nil {
		t.Fatalf("weather lastUpdated parse error: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, payload.LastUpdated.Calendar); err != nil {
		t.Fatalf("calendar lastUpdated parse error: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, payload.LastUpdated.Tasks); err != nil {
		t.Fatalf("tasks lastUpdated parse error: %v", err)
	}
}

func TestGetCalendar(t *testing.T) {
	router := setupTestRouter(t)
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
	router := setupTestRouter(t)
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
	router := setupTestRouter(t)
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
	router := setupTestRouter(t)
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
