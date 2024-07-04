package version_test

import (
	"github.com/yanshicheng/ikube-gin-starter/pkg/version"
	"testing"
)

func TestVersion(t *testing.T) {
	version.IkubeopsGoVersion = "go1.22.3"
	version.IkubeopsCommit = "sanodijewnfiw"
	version.IkubeopsBranch = "master"
	version.IkubeopsBuildTime = "2024-07-04 11:30:22"
	version.IkubeopsTag = "v1.1.1"
	t.Logf("\n%s", version.FullTagVersion())
	t.Logf("\n%s", version.ShortTagVersion())
}
