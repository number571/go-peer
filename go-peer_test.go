package gopeer

import (
	"os"
	"regexp"
	"testing"
)

func TestGoPeerVersion(t *testing.T) {
	changelog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Error(err)
		return
	}

	re := regexp.MustCompile(`##\s+(v\d+\.\d+\.\d+)\s+`)
	match := re.FindAllStringSubmatch(string(changelog), -1)
	if len(match) < 2 {
		t.Error("versions not found")
		return
	}

	// current version is always previous version in the changelog
	if match[1][1] != CVersion {
		t.Error("the versions do not match")
		return
	}
}
