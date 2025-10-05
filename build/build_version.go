// nolint: err113
package build

import (
	_ "embed"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	//go:embed version.yml
	gVersionVal []byte
	gVersion    string
)

func init() {
	var versionYAML struct {
		FVersion string `yaml:"version"`
	}
	if err := encoding.DeserializeYAML(gVersionVal, &versionYAML); err != nil {
		panic(err)
	}
	gVersion = versionYAML.FVersion
}

func GetVersion() string {
	return gVersion
}
