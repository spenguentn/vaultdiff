package filter_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "db/password", Change: diff.Modified, Left: "old", Right: "new"},
	{Key: "db/user", Change: diff.Unchanged, Left: "admin", Right: "admin"},
	{Key: "app/secret", Change: diff.Added, Left: "", Right: "value"},
	{Key: "app/debug", Change: diff.Removed, Left: "true", Right: ""},
}

func TestApply_NoOptions(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{})
	if len(got) != len(sampleResults) {
		t.Fatalf("expected %d results, got %d", len(sampleResults), len(got))
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{OnlyChanged: true})
	for _, r := range got {
		if r.Change == diff.Unchanged {
			t.Errorf("unexpected unchanged key %q in results", r.Key)
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 changed results, got %d", len(got))
	}
}

func TestApply_Prefix(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{Prefix: "app/"})
	if len(got) != 2 {
		t.Fatalf("expected 2 results with prefix app/, got %d", len(got))
	}
	for _, r := range got {
		if r.Key != "app/secret" && r.Key != "app/debug" {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestApply_KeyAllowlist(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{Keys: []string{"db/password", "app/debug"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_PrefixAndOnlyChanged(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{Prefix: "db/", OnlyChanged: true})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Key != "db/password" {
		t.Errorf("expected db/password, got %q", got[0].Key)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := filter.Apply(nil, filter.Options{OnlyChanged: true, Prefix: "x/"})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}
