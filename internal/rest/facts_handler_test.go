package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chrismason/pet-me/internal/config"
	"github.com/chrismason/pet-me/internal/models"
)

func TestFactsHandler(t *testing.T) {
	tt := []struct {
		name             string
		count            string
		server           *httptest.Server
		expectedResponse *models.FactResponse
		status           int
		expectedError    string
	}{
		{
			name:          "invalid count parameter",
			count:         "a",
			status:        http.StatusBadRequest,
			expectedError: "invalid parameter of 'factCount'",
		},
		{
			name:          "single fact",
			count:         "1",
			status:        http.StatusOK,
			expectedError: "",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"fact": "A fact", "length": 6}`))
			})),
			expectedResponse: &models.FactResponse{
				Facts: []string{"A fact"},
			},
		},
		{
			name:          "single fact with error",
			count:         "1",
			status:        http.StatusInternalServerError,
			expectedError: "failed calling fact API",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("500"))
			})),
			expectedResponse: nil,
		},
		{
			name:          "multiple facts",
			count:         "2",
			status:        http.StatusOK,
			expectedError: "",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data": [{"fact": "A fact", "length": 6},{"fact": "Another fact", "length": 11}]}`))
			})),
			expectedResponse: &models.FactResponse{
				Facts: []string{"A fact", "Another fact"},
			},
		},
		{
			name:          "multiple facts with error",
			count:         "2",
			status:        http.StatusInternalServerError,
			expectedError: "failed calling fact API",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("500"))
			})),
			expectedResponse: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.ServerConfig{}
			if tc.server != nil {
				cfg.FactsAPI = tc.server.URL
				defer tc.server.Close()
			}
			req, err := http.NewRequest("GET", "localhost:8080/facts/cats?factCount="+tc.count, nil)
			if err != nil {
				t.Errorf("could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			status, err := getCatFactsInner(&cfg, rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if status != tc.status {
				t.Errorf("expected status %v, got %v", tc.status, res.StatusCode)
			}

			if tc.expectedResponse != nil {
				facts := models.FactResponse{}
				err = json.NewDecoder(res.Body).Decode(&facts)
				if err != nil {
					t.Fatalf("error parsing response")
				}
				if !reflect.DeepEqual(&facts, tc.expectedResponse) {
					t.Fatalf("expected %v, found %v", tc.expectedResponse, facts)
				}
			}
			if tc.expectedError != "" {
				if err.Error() != tc.expectedError {
					t.Errorf("expected error response of '%v', got '%v'", tc.expectedError, err)
				}
			}
		})
	}
}
