package trello

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBoards(t *testing.T) {
	boards := []Board{{ID: "board1", Name: "Sprint Board"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/members/me/boards" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("filter") != "open" {
			t.Error("expected filter=open")
		}
		if r.URL.Query().Get("key") == "" || r.URL.Query().Get("token") == "" {
			t.Error("missing auth params")
		}
		json.NewEncoder(w).Encode(boards)
	}))
	defer server.Close()

	client := NewClient("testkey", "testtoken", WithBaseURL(server.URL))
	result, err := client.GetBoards()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Name != "Sprint Board" {
		t.Errorf("unexpected boards: %+v", result)
	}
}
