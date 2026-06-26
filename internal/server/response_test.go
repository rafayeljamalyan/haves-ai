package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"haves/internal/db"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// assertEnvelope verifies a response body conforms exactly to the contract:
//
//	success: {"meta": {"error": null},                        "data": <anything non-null>}
//	failure: {"meta": {"error": {"message": str, "code": int}}, "data": null}
//
// It fails the test on any deviation: missing/extra top-level keys, a missing
// error key in meta, a non-null data on failure, or a malformed error object.
func assertEnvelope(t *testing.T, body []byte, wantErr bool) {
	t.Helper()

	var top map[string]json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		t.Fatalf("body is not a JSON object: %v (body: %s)", err, body)
	}

	// Exactly two top-level keys: meta and data.
	if _, ok := top["meta"]; !ok {
		t.Errorf("missing top-level %q key (body: %s)", "meta", body)
	}
	if _, ok := top["data"]; !ok {
		t.Errorf("missing top-level %q key (body: %s)", "data", body)
	}
	if len(top) != 2 {
		t.Errorf("expected exactly 2 top-level keys (meta, data), got %d: %s", len(top), body)
	}

	// meta must carry an error key (null or object).
	var meta map[string]json.RawMessage
	if err := json.Unmarshal(top["meta"], &meta); err != nil {
		t.Fatalf("meta is not a JSON object: %v (meta: %s)", err, top["meta"])
	}
	rawErr, ok := meta["error"]
	if !ok {
		t.Fatalf("meta missing %q key (meta: %s)", "error", top["meta"])
	}

	errIsNull := string(rawErr) == "null"
	dataIsNull := string(top["data"]) == "null"

	if wantErr {
		if errIsNull {
			t.Errorf("expected meta.error to be set on failure, got null")
		} else {
			// error object must be exactly {message: string, code: int}.
			var apiErr struct {
				Message *string  `json:"message"`
				Code    *float64 `json:"code"` // JSON numbers decode to float64
			}
			var extra map[string]json.RawMessage
			if err := json.Unmarshal(rawErr, &apiErr); err != nil {
				t.Fatalf("meta.error is not an object: %v (error: %s)", err, rawErr)
			}
			_ = json.Unmarshal(rawErr, &extra)
			if apiErr.Message == nil {
				t.Errorf("meta.error missing string %q (error: %s)", "message", rawErr)
			}
			if apiErr.Code == nil {
				t.Errorf("meta.error missing numeric %q (error: %s)", "code", rawErr)
			} else if *apiErr.Code != float64(int(*apiErr.Code)) {
				t.Errorf("meta.error.code must be an integer, got %v", *apiErr.Code)
			}
			if len(extra) != 2 {
				t.Errorf("meta.error must have exactly {message, code}, got %d keys: %s", len(extra), rawErr)
			}
		}
		if !dataIsNull {
			t.Errorf("expected data to be null on failure, got: %s", top["data"])
		}
	} else {
		if !errIsNull {
			t.Errorf("expected meta.error to be null on success, got: %s", rawErr)
		}
		if dataIsNull {
			t.Errorf("expected data to be non-null on success, got null")
		}
	}
}

func TestRespondShape(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respond(c, http.StatusOK, gin.H{"message": "pong"})

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	assertEnvelope(t, w.Body.Bytes(), false)

	// Exact byte-for-byte shape on success.
	const want = `{"meta":{"error":null},"data":{"message":"pong"}}`
	if got := w.Body.String(); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestRespondErrorShape(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, http.StatusNotFound, 1001, "thing not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	assertEnvelope(t, w.Body.Bytes(), true)

	// Exact byte-for-byte shape on failure.
	const want = `{"meta":{"error":{"message":"thing not found","code":1001}},"data":null}`
	if got := w.Body.String(); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

// newTestEngine builds the real router wired to a DB whose pool points at an
// unreachable address, so readiness checks fail without needing a live server.
func newTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), "postgres://u:p@127.0.0.1:1/none")
	if err != nil {
		t.Fatalf("building test pool: %v", err)
	}
	t.Cleanup(pool.Close)
	return New(&db.DB{Pool: pool})
}

// TestAllRoutesEnvelope drives every registered route through the real engine
// and asserts each response conforms to the envelope contract.
func TestAllRoutesEnvelope(t *testing.T) {
	engine := newTestEngine(t)

	cases := []struct {
		name     string
		method   string
		path     string
		wantCode int
		wantErr  bool
	}{
		{"health", http.MethodGet, "/health", http.StatusOK, false},
		{"ping", http.MethodGet, "/api/v1/ping", http.StatusOK, false},
		{"ready_db_down", http.MethodGet, "/ready", http.StatusServiceUnavailable, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			engine.ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Errorf("status = %d, want %d (body: %s)", w.Code, tc.wantCode, w.Body)
			}
			assertEnvelope(t, w.Body.Bytes(), tc.wantErr)
		})
	}
}
