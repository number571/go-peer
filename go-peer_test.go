package gopeer

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestGoPeerVersion(t *testing.T) {
	t.Parallel()

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

	if match[0][1] == match[1][1] {
		t.Error("the same versions inline")
		return
	}

	for i := 0; i < len(match); i++ {
		for j := i + 1; j < len(match)-1; j++ {
			if match[i][1] == match[j][1] {
				t.Errorf("found the same versions (i=%d, j=%d)", i, j)
				return
			}
		}
	}

	if bytes.Count(changelog, []byte("*??? ??, ????*")) != 1 {
		t.Error("is there no new version or more than one new version?")
		return
	}
}
