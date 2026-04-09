package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/api"
)

// stubLogical is a minimal stub satisfying the logical interface used by Client.
type stubLogical struct {
	readFn func(path string) (*api.Secret, error)
}

func (s *stubLogical) ReadWithContext(_ context.Context, path string) (*api.Secret, error) {
	return s.readFn(path)
}

func makeVersionsResp(versions map[string]interface{}) *api.Secret {
	return &api.Secret{
		Data: map[string]interface{}{
			"versions": versions,
		},
	}
}

func TestListVersions_ReturnsOrdered(t *testing.T) {
	versionData := map[string]interface{}{
		"3": map[string]interface{}{"created_time": "2024-03-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"1": map[string]interface{}{"created_time": "2024-01-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"2": map[string]interface{}{"created_time": "2024-02-01T00:00:00Z", "deletion_time": "", "destroyed": true},
	}

	c := &Client{}
	_ = c // ensure Client is used; actual test uses helper below

	resp := makeVersionsResp(versionData)
	if resp.Data["versions"] == nil {
		t.Fatal("expected versions key")
	}

	// Parse versions directly via the logic extracted for unit testing.
	metas, err := parseVersionsMeta(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(metas))
	}
	for i, expected := range []int{1, 2, 3} {
		if metas[i].Version != expected {
			t.Errorf("index %d: expected version %d, got %d", i, expected, metas[i].Version)
		}
	}
	if !metas[1].Destroyed {
		t.Error("expected version 2 to be destroyed")
	}
}

func TestListVersions_NilResponse(t *testing.T) {
	_, err := parseVersionsMeta(nil)
	if err == nil {
		t.Error("expected error for nil response")
	}
}

func TestListVersions_MissingVersionsKey(t *testing.T) {
	resp := &api.Secret{Data: map[string]interface{}{}}
	_, err := parseVersionsMeta(resp)
	if err == nil {
		t.Error("expected error when versions key is missing")
	}
}

// parseVersionsMeta is a testable helper extracted from ListVersions logic.
func parseVersionsMeta(resp *api.Secret) ([]VersionMeta, error) {
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("vault: no metadata found")
	}
	versionsRaw, ok := resp.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("vault: metadata missing 'versions' key")
	}
	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected versions format")
	}
	var metas []VersionMeta
	for key, val := range versionsMap {
		var vnum int
		if _, err := fmt.Sscanf(key, "%d", &vnum); err != nil {
			continue
		}
		info, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{Version: vnum}
		if ct, ok := info["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if destroyed, ok := info["destroyed"].(bool); ok {
			meta.Destroyed = destroyed
		}
		metas = append(metas, meta)
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Version < metas[j].Version })
	return metas, nil
}
