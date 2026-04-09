package vault

import (
	"context"
	"fmt"
	"time"
)

// SecretMetadata holds metadata about a KV v2 secret path.
type SecretMetadata struct {
	Path           string
	CurrentVersion int
	OldestVersion  int
	CreatedTime    time.Time
	UpdatedTime    time.Time
	Versions       map[string]VersionMeta
}

// VersionMeta holds per-version metadata.
type VersionMeta struct {
	CreatedTime  time.Time
	DeletionTime time.Time
	Destroyed    bool
}

// ReadMetadata fetches KV v2 metadata for the given mount and secret path.
func (c *Client) ReadMetadata(ctx context.Context, mount, path string) (*SecretMetadata, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)
	resp, err := c.logical.ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", path, err)
	}
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("no metadata found for %q", path)
	}
	return parseMetadata(path, resp.Data)
}

func parseMetadata(path string, data map[string]interface{}) (*SecretMetadata, error) {
	meta := &SecretMetadata{Path: path, Versions: make(map[string]VersionMeta)}

	if v, ok := data["current_version"]; ok {
		if n, ok := v.(json.Number); ok {
			cv, _ := n.Int64()
			meta.CurrentVersion = int(cv)
		}
	}
	if v, ok := data["oldest_version"]; ok {
		if n, ok := v.(json.Number); ok {
			ov, _ := n.Int64()
			meta.OldestVersion = int(ov)
		}
	}
	if v, ok := data["created_time"]; ok {
		if s, ok := v.(string); ok {
			meta.CreatedTime, _ = time.Parse(time.RFC3339Nano, s)
		}
	}
	if v, ok := data["updated_time"]; ok {
		if s, ok := v.(string); ok {
			meta.UpdatedTime, _ = time.Parse(time.RFC3339Nano, s)
		}
	}

	if versions, ok := data["versions"].(map[string]interface{}); ok {
		for k, raw := range versions {
			vm, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			var vmeta VersionMeta
			if ct, ok := vm["created_time"].(string); ok {
				vmeta.CreatedTime, _ = time.Parse(time.RFC3339Nano, ct)
			}
			if dt, ok := vm["deletion_time"].(string); ok && dt != "" {
				vmeta.DeletionTime, _ = time.Parse(time.RFC3339Nano, dt)
			}
			if d, ok := vm["destroyed"].(bool); ok {
				vmeta.Destroyed = d
			}
			meta.Versions[k] = vmeta
		}
	}
	return meta, nil
}
