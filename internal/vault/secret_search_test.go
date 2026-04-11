package vault

import (
	"testing"
)

var searchData = map[string]map[string]string{
	"secret/app/config": {
		"db_password": "hunter2",
		"db_user":     "admin",
		"api_key":     "abc123",
	},
	"secret/app/feature": {
		"feature_flag": "true",
		"db_host":      "localhost",
	},
	"infra/network/vpn": {
		"psk": "secret",
	},
}

func TestSearchSecrets_NoQuery(t *testing.T) {
	results := SearchSecrets(searchData, SearchOptions{})
	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}
}

func TestSearchSecrets_QueryMatchesSubstring(t *testing.T) {
	results := SearchSecrets(searchData, SearchOptions{Query: "db_"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results for 'db_', got %d", len(results))
	}
}

func TestSearchSecrets_QueryCaseInsensitive(t *testing.T) {
	results := SearchSecrets(searchData, SearchOptions{Query: "API_KEY"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'API_KEY', got %d", len(results))
	}
	if results[0].Key != "api_key" {
		t.Errorf("expected key 'api_key', got %q", results[0].Key)
	}
}

func TestSearchSecrets_MountFilter(t *testing.T) {
	results := SearchSecrets(searchData, SearchOptions{Mount: "infra"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result for mount 'infra', got %d", len(results))
	}
	if results[0].Mount != "infra" {
		t.Errorf("expected mount 'infra', got %q", results[0].Mount)
	}
}

func TestSearchSecrets_MaxResults(t *testing.T) {
	results := SearchSecrets(searchData, SearchOptions{MaxResults: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchResult_FullPath(t *testing.T) {
	r := SearchResult{Mount: "secret", Path: "app/config", Key: "db_password"}
	want := "secret/app/config#db_password"
	if r.FullPath() != want {
		t.Errorf("FullPath() = %q, want %q", r.FullPath(), want)
	}
}

func TestSplitLocation_Valid(t *testing.T) {
	m, p, ok := splitLocation("secret/app/config")
	if !ok || m != "secret" || p != "app/config" {
		t.Errorf("unexpected split: mount=%q path=%q ok=%v", m, p, ok)
	}
}

func TestSplitLocation_NoSlash(t *testing.T) {
	m, p, ok := splitLocation("secret")
	if ok || m != "secret" || p != "" {
		t.Errorf("unexpected split: mount=%q path=%q ok=%v", m, p, ok)
	}
}
