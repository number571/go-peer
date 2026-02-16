package gopeer

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"
)

const (
	cVersion = "v1.7.13"
)

var (
	//go:embed CHANGELOG.md
	tgCHANGELOG string
)

func TestGoPeerVersion(t *testing.T) {
	t.Parallel()

	re := regexp.MustCompile(`##\s+(v\d+\.\d+\.\d+~?)\s+`)
	match := re.FindAllStringSubmatch(tgCHANGELOG, -1)
	if len(match) < 2 {
		t.Fatal("versions not found")
	}

	version := cVersion
	if strings.HasSuffix(version, "~") {
		if match[0][1] != version {
			t.Fatal("the versions do not match")
		}
	} else {
		// current version is always previous version in the changelog
		if match[1][1] != version {
			t.Fatal("the versions do not match")
		}
	}

	if match[0][1] == match[1][1] {
		t.Fatal("the same versions inline")
	}

	for i := 0; i < len(match); i++ {
		for j := i + 1; j < len(match)-1; j++ {
			if match[i][1] == match[j][1] {
				t.Fatalf("found the same versions (i=%d, j=%d)", i, j)
			}
		}
	}

	if strings.Count(tgCHANGELOG, "*??? ??, ????*") != 1 {
		t.Fatal("is there no new version or more than one new version?")
	}
}
