package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionMeta holds metadata about a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KVv2 secret at path.
func (c *Client) ListVersions(ctx context.Context, mount, path string) ([]VersionMeta, error) {
	secretPath := fmt.Sprintf("%s/metadata/%s", mount, path)

	resp, err := c.logical.ReadWithContext(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("vault: list versions for %q: %w", path, err)
	}
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("vault: no metadata found for %q", path)
	}

	versionsRaw, ok := resp.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("vault: metadata missing 'versions' key for %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected versions format for %q", path)
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
		if dt, ok := info["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if destroyed, ok := info["destroyed"].(bool); ok {
			meta.Destroyed = destroyed
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return metas, nil
}
