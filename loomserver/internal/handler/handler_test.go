package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rotscher/httploom/loomserver/internal/api"
	"github.com/rotscher/httploom/loomserver/internal/weave"
)

func TestAPI(t *testing.T) {
	srv := httptest.NewServer(api.Handler(New(weave.NewStore(), "test")))
	t.Cleanup(srv.Close)

	t.Run("health", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/health")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var health api.Health
		if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK || health.Status != "ok" || health.Version != "test" {
			t.Fatalf("unexpected health response: %d %+v", resp.StatusCode, health)
		}
	})

	var created api.Weave

	t.Run("create weave", func(t *testing.T) {
		resp, err := http.Post(
			srv.URL+"/weaves",
			"application/json",
			strings.NewReader(`{"text":"The quick brown fox jumps over the lazy dog","width":20}`),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
		if created.Id == "" {
			t.Fatal("expected an id")
		}
		if want := "The quick brown fox\njumps over the lazy\ndog"; created.Text != want {
			t.Fatalf("text = %q, want %q", created.Text, want)
		}
	})

	t.Run("get weave", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/weaves/" + created.Id)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var got api.Weave
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK || got.Id != created.Id {
			t.Fatalf("unexpected response: %d %+v", resp.StatusCode, got)
		}
	})

	t.Run("list weaves", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/weaves")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var list api.WeaveList
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK || len(list.Items) != 1 {
			t.Fatalf("unexpected response: %d %+v", resp.StatusCode, list)
		}
	})

	t.Run("rejects empty text", func(t *testing.T) {
		resp, err := http.Post(srv.URL+"/weaves", "application/json", strings.NewReader(`{"text":"  "}`))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for unknown weave", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/weaves/missing")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})
}
