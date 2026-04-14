package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeVaultResponse(t *testing.T, data map[string]interface{}) []byte {
	t.Helper()
	envelope := map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
		},
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshalling test response: %v", err)
	}
	return b
}

func TestGetSecret_Success(t *testing.T) {
	expected := map[string]interface{}{"DB_HOST": "localhost", "DB_PORT": "5432"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(makeVaultResponse(t, expected))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "test-token")
	data, err := client.GetSecret("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", data["DB_HOST"])
	}
	if data["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", data["DB_PORT"])
	}
}

func TestGetSecret_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["secret not found"]}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "test-token")
	_, err := client.GetSecret("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestParseKVv2Response_Valid(t *testing.T) {
	body := makeVaultResponse(t, map[string]interface{}{"KEY": "value"})
	data, err := parseKVv2Response(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", data["KEY"])
	}
}

func TestParseKVv2Response_InvalidJSON(t *testing.T) {
	_, err := parseKVv2Response([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
