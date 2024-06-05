package conf

import (
	"fmt"
	"os"
	"testing"
)

func TestGetTemplates(t *testing.T) {
	targetDir := t.TempDir()
	if err := os.MkdirAll(targetDir+"/tpl", 0755); err != nil {
		t.Error(err)
		t.FailNow()
	}
	want := []string{"bookmark.tpl", "default.tpl.md", "weekly-planning.tpl"}
	for i := range want {
		path := fmt.Sprintf("%s/tpl/%s", targetDir, want[i])
		if _, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644); err != nil {
			t.Errorf("cannot set up test: %v", err)
			t.FailNow()
		}
	}
	got, err := GetTemplates(targetDir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(got) != len(want) {
		t.Errorf("Expected %d templates, got %d", len(want), len(got))
		t.FailNow()
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("Expected template '%s' at %d, got '%s'", want[i], i, got[i])
		}
	}
}
