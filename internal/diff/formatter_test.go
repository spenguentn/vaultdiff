package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{Writer: &buf}

	data := map[string]interface{}{"key": "value"}
	result := Compare("secret/app", data, data)

	if err := f.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{Writer: &buf}

	oldData := map[string]interface{}{"password": "old", "removed_key": "gone"}
	newData := map[string]interface{}{"password": "new", "added_key": "here"}
	result := Compare("secret/app", oldData, newData)

	if err := f.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ added_key") {
		t.Errorf("expected added_key in output: %s", out)
	}
	if !strings.Contains(out, "- removed_key") {
		t.Errorf("expected removed_key in output: %s", out)
	}
	if !strings.Contains(out, "~ password") {
		t.Errorf("expected modified password in output: %s", out)
	}
}

func TestTextFormatter_MaskSecrets(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{MaskSecrets: true, Writer: &buf}

	oldData := map[string]interface{}{}
	newData := map[string]interface{}{"token": "supersecret"}
	result := Compare("secret/app", oldData, newData)

	_ = f.Write(result)

	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Error("secret value should be masked")
	}
	if !strings.Contains(out, "s**********") {
		t.Errorf("expected masked value in output: %s", out)
	}
}

func TestSummary(t *testing.T) {
	oldData := map[string]interface{}{"a": "1", "b": "old"}
	newData := map[string]interface{}{"b": "new", "c": "3"}
	result := Compare("secret/app", oldData, newData)

	summary := Summary(result)
	if !strings.Contains(summary, "+1") {
		t.Errorf("expected +1 in summary: %s", summary)
	}
	if !strings.Contains(summary, "-1") {
		t.Errorf("expected -1 in summary: %s", summary)
	}
	if !strings.Contains(summary, "~1") {
		t.Errorf("expected ~1 in summary: %s", summary)
	}
}
