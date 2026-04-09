package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/snapshot"
)

func makeSnap(path string) *snapshot.Snapshot {
	return snapshot.New(path, 1, map[string]string{"KEY": "val"}, snapshot.Meta{})
}

func TestStore_SaveAndGet(t *testing.T) {
	store := snapshot.NewStore()
	snap := makeSnap("secret/app")
	store.Save("app-v1", snap)

	got, err := store.Get("app-v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", got.Path)
	}
}

func TestStore_Get_NotFound(t *testing.T) {
	store := snapshot.NewStore()
	_, err := store.Get("missing")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestStore_Delete(t *testing.T) {
	store := snapshot.NewStore()
	store.Save("x", makeSnap("secret/x"))
	if err := store.Delete("x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Len() != 0 {
		t.Error("expected empty store after delete")
	}
}

func TestStore_Delete_NotFound(t *testing.T) {
	store := snapshot.NewStore()
	if err := store.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent snapshot")
	}
}

func TestStore_List(t *testing.T) {
	store := snapshot.NewStore()
	store.Save("a", makeSnap("secret/a"))
	store.Save("b", makeSnap("secret/b"))
	names := store.List()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestStore_Len(t *testing.T) {
	store := snapshot.NewStore()
	if store.Len() != 0 {
		t.Error("expected empty store")
	}
	store.Save("one", makeSnap("secret/one"))
	if store.Len() != 1 {
		t.Errorf("expected 1, got %d", store.Len())
	}
}
